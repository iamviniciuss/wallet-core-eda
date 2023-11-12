package infra

import infra "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/errors"

type HttpService interface {
	Get(path string, callback func(map[string]string, []byte, QueryParams) (interface{}, *infra.IntegrationError))
	Post(path string, callback func(map[string]string, []byte, QueryParams) (interface{}, *infra.IntegrationError))
	Put(path string, callback func(map[string]string, []byte, QueryParams) (interface{}, *infra.IntegrationError))
	ListenAndServe(port string) error
}

type QueryParams interface {
	GetParam(key string) []byte
	AddParam(key string, value string)
}
