package router

import (
	"net/http"

	"github.com/ZnNr/songs-library/internal/handlers"
	"github.com/ZnNr/songs-library/internal/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// NewRouter создает новый маршрутизатор и регистрирует маршруты.
func NewRouter(handler *handlers.SongHandler, logger *zap.Logger) *mux.Router {
	r := mux.NewRouter()

	// Добавляем миддлвары для логирования
	r.Use(middleware.LoggingMiddleware(logger))

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/songs", handler.GetSongs).Methods(http.MethodGet)
	api.HandleFunc("/songs/{id}/lyrics", handler.GetLyrics).Methods(http.MethodGet)
	api.HandleFunc("/songs", handler.CreateSong).Methods(http.MethodPost)
	api.HandleFunc("/songs/{id}", handler.UpdateSong).Methods(http.MethodPut)
	api.HandleFunc("/songs/{id}", handler.DeleteSong).Methods(http.MethodDelete)

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	return r
}
