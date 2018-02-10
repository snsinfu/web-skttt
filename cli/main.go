package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/gorilla/websocket"
)

const usage = `
WebSocket client.

Usage:
  skttt [-t <second>] <url>
  skttt -h

Options:
  -t, --timeout <second>  Specify timeout in seconds [default: 5]
  -h, --help              Show usage and exit
`

type command struct {
	URL     string `docopt:"<url>"`
	Timeout time.Duration
	Help    bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		return err
	}

	cmd := command{}
	if err := opts.Bind(&cmd); err != nil {
		return err
	}
	cmd.Timeout *= time.Second

	return cmd.run()
}

func (cmd *command) run() error {
	ws, res, err := websocket.DefaultDialer.Dial(cmd.URL, http.Header{})
	if err != nil {
		return err
	}
	defer ws.Close()

	if res.StatusCode >= 400 {
		return errors.New(res.Status)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGPIPE)

	readCh := make(chan error, 1)
	go func() {
		readCh <- cmd.read(ws)
	}()

	select {
	case <-sigCh:
		ws.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseGoingAway, ""),
			time.Now().Add(cmd.Timeout),
		)

	case err := <-readCh:
		if err != nil {
			return err
		}
	}

	return nil
}

func (cmd *command) read(ws *websocket.Conn) error {
	for {
		ty, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}
			return err
		}

		switch ty {
		case websocket.TextMessage:
			fmt.Println(string(msg))

		case websocket.BinaryMessage:
			fmt.Println(base64.StdEncoding.EncodeToString(msg))

		case websocket.CloseMessage:
			return nil

		default:
		}
	}
}
