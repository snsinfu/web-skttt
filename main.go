package main

import (
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/snsinfu/web-skttt/action"
	"github.com/snsinfu/web-skttt/domain"
)

var (
	defaultDomainConfig = domain.Config{
		NameLen:    8,
		KeyLen:     16,
		Expiration: 24 * time.Hour,
	}
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("1K"))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		ContentSecurityPolicy: "frame-ancestors 'none'",
	}))

	dom, err := makeDomain()
	if err != nil {
		e.Logger.Fatal(err)
	}

	act, err := makeAction(dom)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/:topic", act.GetTopic)
	e.POST("/:topic", act.PostTopic)
	e.POST("/:topic/:key", act.PostTopic)
	e.POST("/request", act.PostRequest)
	e.Static("/", "public")

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}

func makeDomain() (*domain.Domain, error) {
	config := defaultDomainConfig

	if err := intEnv("SKTTT_TOPIC_NAME_LENGTH", &config.NameLen); err != nil {
		return nil, err
	}

	if err := intEnv("SKTTT_TOPIC_KEY_LENGTH", &config.KeyLen); err != nil {
		return nil, err
	}

	if err := durationEnv("SKTTT_TOPIC_EXPIRATION", &config.Expiration); err != nil {
		return nil, err
	}

	return domain.New(config)
}

func makeAction(dom *domain.Domain) (*action.Action, error) {
	return action.New(dom)
}

func intEnv(key string, dest *int) error {
	if env, ok := os.LookupEnv(key); ok {
		val, err := strconv.ParseInt(env, 10, 0)
		if err != nil {
			return err
		}
		*dest = int(val)
	}
	return nil
}

func durationEnv(key string, dest *time.Duration) error {
	if env, ok := os.LookupEnv(key); ok {
		val, err := strconv.ParseInt(env, 10, 0)
		if err != nil {
			return err
		}
		*dest = time.Duration(val) * time.Second
	}
	return nil
}
