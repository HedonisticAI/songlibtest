package paginator

import (
	"fmt"
	"net/http"
	"songlibtest/pkg/logger"
	"songlibtest/pkg/postgres"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Paginator struct {
	Logger         *logger.Logger
	PageSize       int
	VerseDelimiter string
	DB             *postgres.Postgres
}

func NewPaginator(Logger *logger.Logger, PageSize int, VerseDelimiter string, DB *postgres.Postgres) PaginatorUseCase {
	return &Paginator{Logger: Logger, PageSize: PageSize, DB: DB, VerseDelimiter: VerseDelimiter}
}

func formQuery(c *gin.Context, Type string, Begin int, End int) string {
	QuertString := "select" + Type + "from Songs "
	count := 0
	if c.Query("name") != "" {
		if count == 0 {
			QuertString += "where "
		} else {
			QuertString += "AND "
		}
		QuertString += "Song = " + c.Query("name")
		count++
	}
	if c.Query("group") != "" {
		if count == 0 {
			QuertString += "where "
		} else {
			QuertString += "AND "
		}
		QuertString += "GroupName = " + c.Query("group")
		count++
	}
	if c.Query("releaseDate") != "" {
		if count == 0 {
			QuertString += "where "
		} else {
			QuertString += "AND "
		}
		QuertString += "ReleaseDate = " + c.Query("releaseDate")
		count++
	}
	if c.Query("link") != "" {
		if count == 0 {
			QuertString += "where "
		} else {
			QuertString += "AND "
		}
		QuertString += "Link = " + c.Query("link")
		count++
	}
	if c.Query("text") != "" {
		if count == 0 {
			QuertString += "where "
		} else {
			QuertString += "AND "
		}
		QuertString += "SongText = " + c.Query("text")
		count++
	}
	if Type == TypeAll {
		if count == 0 {
			QuertString += "where "
		} else {
			QuertString += "AND "
		}
		QuertString += fmt.Sprintf(" ID > %d AND ID < %d", Begin, End)
		QuertString += "ORDER BY ID"
	}
	return QuertString
}

func (Paginator *Paginator) GetSongsWithPagination(c *gin.Context) {
	var Songs SongsResponse
	Paginator.Logger.Info("Got GetSong request")
	Page, err := strconv.Atoi(c.Query("page"))
	Begin := (Page - 1) * Paginator.PageSize
	End := Page*Paginator.PageSize - 1
	if err != nil || Page < 0 {
		c.String(http.StatusBadRequest, ErrBadPageNumber.Error())
		return
	}
	QuertString := formQuery(c, TypeAll, Begin, End)
	rows, err := Paginator.DB.DB.Query(QuertString)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	for rows.Next() {
		Song := SongData{}
		rows.Scan(&Song.ID, &Song.Name, &Song.Group, &Song.Date, &Song.Text, &Song.Link)
		Songs.Songs = append(Songs.Songs, Song)
	}
	c.JSON(http.StatusOK, Songs)
}

func (Paginator *Paginator) GetSongTextWithPagination(c *gin.Context) {
	var Text string
	Paginator.Logger.Info("Got GetSongText request")
	Page, err := strconv.Atoi(c.Query("page"))
	Begin := (Page - 1) * Paginator.PageSize
	End := Page*Paginator.PageSize -1
	if err != nil || Page < 0 {
		c.String(http.StatusBadRequest, ErrBadPageNumber.Error())
		return
	}
	QueryString := formQuery(c, TypeText, Begin, End)
	row := Paginator.DB.DB.QueryRow(QueryString)
	if row.Err() != nil {
		c.String(http.StatusBadRequest, row.Err().Error())
		return
	}
	row.Scan(&Text)
	str := strings.Split(Text, Paginator.VerseDelimiter)
	Verses := SongText{Verses: str[Begin:End]}
	c.JSON(http.StatusOK, Verses)
}
