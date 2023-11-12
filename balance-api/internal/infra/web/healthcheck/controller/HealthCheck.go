package infra

import (
	"encoding/json"
	"fmt"
	"os"

	errors "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/errors"
	http "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/http"
)

type AWSSecret struct {
	AWS_SECRET string `json:AWS_SECRET`
}

// HealthCheckCtrl godoc
// @Summary Health Check
// @Description check if app is running
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthCheck
// @Failure 500 {object} HTTPError
// @Router /health-check [get]
// @Security		ApiKeyAuth
// @Param			Authorization	header		string	true	"Bearer token"
// @Tags         External
func HealthCheckCtrl(params map[string]string, body []byte, queryArgs http.QueryParams) (interface{}, *errors.IntegrationError) {
	fmt.Println("####################### HEALTH CHECK INIT #######################")
	awsSecret := os.Getenv("AWS_SECRET")
	aWSSecretStruct := AWSSecret{}
	json.Unmarshal([]byte(awsSecret), &aWSSecretStruct)
	fmt.Println(aWSSecretStruct.AWS_SECRET)
	fmt.Println("####################### HEALTH CHECK FINISH #######################")

	return "OK", nil
}

type HealthCheck struct {
	Status string `json:"status"`
}

type HTTPError struct {
	Status  string
	Message string
}
