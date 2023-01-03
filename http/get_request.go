package http

import (
	"net/http"

	"github.com/anboo/throttler/service/storage"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type GetRequest struct {
	storage storage.Storage
}

func NewGetRequestHandler(storage storage.Storage) *CreateRequest {
	return &CreateRequest{
		storage: storage,
	}
}

func (h *GetRequest) Handler(c *gin.Context) {
	requestId := c.Param("id")

	r, err := h.storage.ByID(c.Request.Context(), requestId)
	switch {
	case errors.Is(err, storage.ErrorNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	case err != nil:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": r.ID, "status": r.Status, "created_at": r.CreatedAt})
}
