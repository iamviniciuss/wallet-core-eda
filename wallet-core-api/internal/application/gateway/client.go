package gateway

import (
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/entity"
)

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
