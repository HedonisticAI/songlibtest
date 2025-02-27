package songdata

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrBadParams = errors.New("bad request params")

type AddSongRequest struct {
	Group string `json:"group,omitempty"`
	Song  string `json:"song,omitempty"`
}

type EnrichmentData struct {
	Date string `json:"releaseDate"`
	Link string `json:"link"`
	Text string `json:"text"`
}

type SongDataUsecase interface {
	AddSong(c *gin.Context)
}
