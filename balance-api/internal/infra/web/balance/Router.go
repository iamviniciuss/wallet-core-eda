package infra

import (
	infra "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/http"
	controller "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/web/balance/controller"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/pkg/uow"
)

func BalanceRouter(http infra.HttpService, Uow uow.UowInterface) {
	http.Get("/balances/:account_id", controller.NewGetBalanceByIdCtrl(Uow).Execute)
}
