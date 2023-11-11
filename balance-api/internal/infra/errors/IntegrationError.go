package infra

import "fmt"

type IntegrationError struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (ie *IntegrationError) Error() string {
	return fmt.Sprintf("Status %d: Integration Error: %s", ie.StatusCode, ie.Message)
}
