package redactor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"songlibtest/pkg/logger"
	"songlibtest/pkg/postgres"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Redactor struct {
	DB     *postgres.Postgres
	Logger *logger.Logger
}

func New(DB *postgres.Postgres, Logger *logger.Logger) RedactorUsecase {
	return &Redactor{DB: DB, Logger: Logger}
}

func (Redactor *Redactor) Delete(c *gin.Context) {
	Redactor.Logger.Info("Got new Delete request")
	StringID := c.Query("ID")
	ID, err := strconv.Atoi(StringID)
	if err != nil {
		Redactor.Logger.Debug(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	Redactor.Logger.Debug("Got new Delete request with ID: %d", ID)
	result, err := Redactor.DB.DB.Exec("delete * from Songs where id = $1", ID)
	rows, _ := result.RowsAffected()
	if err != nil {
		Redactor.Logger.Debug(err)
		c.String(http.StatusNoContent, err.Error())
		return
	}
	if rows == 0 {
		Redactor.Logger.Debug(ErrBadID)
		c.String(http.StatusNoContent, ErrBadID.Error())
		return
	}
	Redactor.Logger.Debug("Deleted row with ID: %d", ID)
	c.String(http.StatusOK, "Song Deleted")
}

func (Redactor *Redactor) Change(c *gin.Context) {
	var Args ChangeRequest
	Redactor.Logger.Info("Got new Change request")
	StringID := c.Query("ID")
	ID, err := strconv.Atoi(StringID)
	if err != nil {
		Redactor.Logger.Debug(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	Redactor.Logger.Debug("Got new Change request with ID: %d", ID)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		Redactor.Logger.Debug(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(body, &Args)
	if err != nil {
		Redactor.Logger.Debug(err)
		c.String(http.StatusOK, err.Error())
		return
	}
	str, err := prepareUpdateString(ID, Args)
	if err != nil {
		Redactor.Logger.Debug(err)
		c.String(http.StatusOK, ErrNoParams.Error())
		return
	}
	result, err := Redactor.DB.DB.Exec(str)
	rows, _ := result.RowsAffected()
	if err != nil || rows == 0 {
		Redactor.Logger.Debug(err)
		c.String(http.StatusNoContent, err.Error())
		return
	}
}

func prepareUpdateString(ID int, Args ChangeRequest) (string, error) {
	res := "update Songs"
	count := 0
	if Args.Name != "" {
		count++
		if count == 1 {
			res += " set"
		}
		res += " Song =" + Args.Name
	}
	if Args.Group != "" {
		count++
		if count == 1 {
			res += " set"
		}
		res += " GroupName = " + Args.Group
	}
	if Args.Date != "" {
		count++
		if count == 1 {
			res += " set"
		}
		res += " ReleaseDate = " + Args.Group
	}
	if Args.Link != "" {
		count++
		if count == 1 {
			res += " set"
		}
		res += " Link = " + Args.Link
	}
	if Args.Text != "" {
		count++
		if count == 1 {
			res += " set"
		}
		res += " SongText = " + Args.Text
	}
	if count == 0 {
		return "", ErrNoParams
	}
	res += fmt.Sprintf("where id = %d", ID)
	return res, nil
}
