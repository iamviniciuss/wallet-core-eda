package infra

import (
	infra "github.com.br/devfullcycle/fc-ms-wallet/balance-api/internal/infra/http"
	controller "github.com.br/devfullcycle/fc-ms-wallet/balance-api/internal/infra/web/balance/controller"
	"github.com.br/devfullcycle/fc-ms-wallet/balance-api/pkg/uow"
)

func BalanceRouter(http infra.HttpService, Uow uow.UowInterface) {
	http.Get("/balances/:account_id", controller.NewGetBalanceByIdCtrl(Uow).Execute)
}
