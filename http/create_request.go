package http

import (
	"net/http"

	"github.com/anboo/throttler/service/storage"
	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	storage storage.Storage
}

func NewCreateRequestHandler(storage storage.Storage) *CreateRequest {
	return &CreateRequest{
		storage: storage,
	}
}

func (h *CreateRequest) Handler(c *gin.Context) {
	r, err := h.storage.Create(c.Request.Context(), storage.Request{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": r.ID})
}
