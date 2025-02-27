package songdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"songlibtest/pkg/http_server"
	"songlibtest/pkg/logger"
	"songlibtest/pkg/postgres"

	"github.com/gin-gonic/gin"
)

type SongData struct {
	DB                 *postgres.Postgres
	Logger             *logger.Logger
	EnrichmentAddr     string
	EnrichmentEndpoint string
}

func New(DB *postgres.Postgres, Logger *logger.Logger, EnrichmentAddr string) SongDataUsecase {
	return &SongData{DB: DB, Logger: Logger, EnrichmentAddr: EnrichmentAddr}
}

func (S *SongData) AddSong(c *gin.Context) {
	var Request AddSongRequest
	var Enrich EnrichmentData
	var ID int
	S.Logger.Info("got AddSong request")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		S.Logger.Debug(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(body, &Request)
	if err != nil {
		S.Logger.Debug(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if Request.Group == "" || Request.Song == "" {
		S.Logger.Debug(ErrBadParams)
		c.String(http.StatusBadRequest, ErrBadParams.Error())
		return
	}
	EnirchRequestStr := fmt.Sprintf("%s%s?group=%s&song=%s", S.EnrichmentAddr, S.EnrichmentEndpoint, Request.Group, Request.Song)
	code, body, err := http_server.SimpleRequest(EnirchRequestStr)
	if code != http.StatusOK || err != nil {
		S.Logger.Debug("Enrichment Request Failed", code, err.Error())
		c.String(http.StatusInternalServerError, "resp_code: %d err: %s", code, err.Error())
		return
	}
	err = json.Unmarshal(body, &Enrich)
	if err != nil {
		S.Logger.Debug(err.Error())
		c.String(http.StatusInternalServerError, err.Error())
	}
	S.DB.DB.QueryRow("insert into Songs (Song, GroupName, ReleaseDate, SongText, Link) values ($1, $2, $3, $4, $5)", Request.Song, Request.Group, Enrich.Date, Enrich.Text, Enrich.Link).Scan(&ID)
	c.String(http.StatusOK, "New Song Added with ID: %d", ID)
}
