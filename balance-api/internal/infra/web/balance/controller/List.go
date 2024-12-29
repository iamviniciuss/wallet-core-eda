package infra

import (
	"context"

	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/application/usecase"
	errors "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/errors"
	http "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/http"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/pkg/uow"
)

type GetBalanceByIdCtrlOutput struct {
	Balance float64 `json:"balance"`
}

type GetBalanceByIdCtrl struct {
	Uow uow.UowInterface
}

func NewGetBalanceByIdCtrl(Uow uow.UowInterface) *GetBalanceByIdCtrl {
	return &GetBalanceByIdCtrl{Uow}
}

func (ctrl *GetBalanceByIdCtrl) Execute(params map[string]string, body []byte, queryArgs http.QueryParams) (interface{}, *errors.IntegrationError) {
	ctx := context.Background()

	output, err := usecase.NewGetBalanceByIdUseCase(ctrl.Uow).Execute(ctx, string(params["account_id"]))

	if err != nil {
		return nil, &errors.IntegrationError{
			StatusCode: 400,
			Message:    err.Error(),
		}
	}

	return &GetBalanceByIdCtrlOutput{
		Balance: output,
	}, nil
}
