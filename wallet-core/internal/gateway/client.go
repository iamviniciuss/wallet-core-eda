package gateway

import (
	"github.com/iamviniciuss/wallet-core-eda/wallet-core/internal/entity"
)

type ClientGateway interface {
	Get(id string) (*entity.Client, error)
	Save(client *entity.Client) error
}
