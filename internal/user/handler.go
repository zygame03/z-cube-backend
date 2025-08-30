package user

import "github.com/gin-gonic/gin"

type handler struct {
	svc *service
}

func NewHandler(svc *service) *handler {
	return &handler{
		svc: svc,
	}
}

func (h *handler) RegisterRoutes(e *gin.Engine) {

}
