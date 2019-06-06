package delivery

import "net/http"

type ServiceDeliveryHTTP interface {
	GetEngine() http.Handler
}
