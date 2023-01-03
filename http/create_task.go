package http

import (
	"github.com/anboo/throttler/service"
	"github.com/gin-gonic/gin"
)

type CreateTask struct {
	client *service.HttpClient
}

func (h *CreateTask) Handler(c *gin.Context) {

}
