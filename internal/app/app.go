package app

import (
	"songlibtest/config"
	"songlibtest/internal/paginator"
	"songlibtest/internal/redactor"
	songdata "songlibtest/internal/song_data"
	"songlibtest/pkg/http_server"
	"songlibtest/pkg/logger"
	"songlibtest/pkg/postgres"
)

type App struct {
	Paginator paginator.PaginatorUseCase
	Song_Data songdata.SongDataUsecase
	Redactor  redactor.RedactorUsecase
	Server    *http_server.HttpServer
}

func New(Cfg config.Config, Logger *logger.Logger) (*App, error) {
	Logger.Info("Building app")
	DB, err := postgres.NewDB(Cfg)
	if err != nil {
		Logger.Info("Couldn't build App")
		Logger.Debug(err)
		return nil, err
	}
	Server := http_server.NewServer(&Cfg)
	Paginator := paginator.NewPaginator(Logger, Cfg.PageSize, Cfg.VerseString, DB)
	Song_Data := songdata.New(DB, Logger, Cfg.EnrichmentAddr)
	Redactor := redactor.New(DB, Logger)
	err = Server.Map(http_server.HTTP_DELETE, "/deleteSong", Redactor.Delete)
	if err != nil {
		return nil, err
	}
	err = Server.Map(http_server.HTTP_PATCH, "/changeSong", Redactor.Change)
	if err != nil {
		return nil, err
	}
	err = Server.Map(http_server.HTTP_GET, "/songText", Paginator.GetSongTextWithPagination)
	if err != nil {
		return nil, err
	}
	err = Server.Map(http_server.HTTP_GET, "/songInfo", Paginator.GetSongsWithPagination)
	if err != nil {
		return nil, err
	}
	err = Server.Map(http_server.HTTP_POST, "/addSong", Song_Data.AddSong)
	if err != nil {
		return nil, err
	}
	return &App{Paginator: Paginator, Song_Data: Song_Data, Redactor: Redactor, Server: Server}, nil
}

func (App *App) Run() {
	App.Server.Run()
}
