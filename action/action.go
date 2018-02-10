package action

import (
	"github.com/snsinfu/web-skttt/domain"
)

type Action struct {
	dom *domain.Domain
}

func New(dom *domain.Domain) (*Action, error) {
	return &Action{dom}, nil
}
