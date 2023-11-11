package gateway

import (
	"github.com.br/devfullcycle/fc-ms-wallet/balance-api/internal/entity"
)

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
