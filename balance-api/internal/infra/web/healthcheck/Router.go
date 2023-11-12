package infra

import (
	infra "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/http"
	healthCheckCtrl "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/web/healthcheck/controller"
)

func HealthCheckRouter(http infra.HttpService) {
	http.Get("/", healthCheckCtrl.HealthCheckCtrl)
}
