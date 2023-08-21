package restful

type ErrorResponse struct {
	Code          string `json:"code"`
	CorrelationID string `json:"correlation_id"`
	Message       string `json:"message"`
}
