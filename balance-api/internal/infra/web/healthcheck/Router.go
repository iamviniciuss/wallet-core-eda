package infra

import (
	infra "github.com.br/devfullcycle/fc-ms-wallet/balance-api/internal/infra/http"
	healthCheckCtrl "github.com.br/devfullcycle/fc-ms-wallet/balance-api/internal/infra/web/healthcheck/controller"
)

func HealthCheckRouter(http infra.HttpService) {
	http.Get("/", healthCheckCtrl.HealthCheckCtrl)
}
