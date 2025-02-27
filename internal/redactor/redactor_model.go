package redactor

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrBadID = errors.New("no entry with this ID found")
var ErrNoParams = errors.New("no params found")

type ChangeRequest struct {
	Name  string `json:"name,omitempty"`
	Group string `json:"group,omitempty"`
	Date  string `json:"date,omitempty"`
	Text  string `json:"text,omitempty"`
	Link  string `json:"link,omitempty"`
}

type RedactorUsecase interface {
	Delete(c *gin.Context)
	Change(c *gin.Context)
}
