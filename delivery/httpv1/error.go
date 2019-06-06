package httpv1

type HandlerError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
