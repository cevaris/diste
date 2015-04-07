package diste

const TCP_CONN = "tcp"

type ServiceRequest struct {
	Args map[string]string
}

type ServiceResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}
