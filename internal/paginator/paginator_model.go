package paginator

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrBadPageNumber = errors.New("bad Page number")
var TypeAll string = "*"
var TypeText = "SongText"

type SongData struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
	Date  string `json:"date"`
	Text  string `json:"text"`
	Link  string `json:"link"`
}

type SongsResponse struct {
	Songs []SongData `json:"songs"`
}

type SongText struct {
	Verses []string `json:"verses"`
}

type PaginatorUseCase interface {
	GetSongsWithPagination(c *gin.Context)
	GetSongTextWithPagination(c *gin.Context)
}
