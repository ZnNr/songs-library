package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ZnNr/songs-library/internal/errors"
	"github.com/ZnNr/songs-library/internal/models"
	"github.com/ZnNr/songs-library/internal/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type SongHandler struct {
	service *service.SongService
	logger  *zap.Logger
}

func NewSongHandler(service *service.SongService, logger *zap.Logger) *SongHandler {
	return &SongHandler{
		service: service,
		logger:  logger,
	}
}

func (h *SongHandler) handleError(w http.ResponseWriter, err error) {
	var status int
	var message string

	if appErr, ok := err.(*errors.Error); ok {
		status = appErr.Status()
		message = appErr.Message
	} else {
		status = http.StatusInternalServerError
		message = "Internal server error"
	}

	h.logger.Error("Request error", zap.Error(err), zap.Int("status", status), zap.String("message", message))

	http.Error(w, message, status)
}

// @Summary Get songs with filtering and pagination
// @Description Get list of songs with optional filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group_name query string false "Group name"
// @Param song_name query string false "Song name"
// @Param from_date query string false "From date (format: 2006-01-02)"
// @Param to_date query string false "To date (format: 2006-01-02)"
// @Param text query string false "Text content"
// @Param link query string false "Link"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} models.SongsResponse
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling GetSongs request")

	filter := &models.SongFilter{
		GroupName: r.URL.Query().Get("group_name"),
		SongName:  r.URL.Query().Get("song_name"),
		Text:      r.URL.Query().Get("text"),
		Link:      r.URL.Query().Get("link"),
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			filter.Page = page
		} else {
			h.handleError(w, errors.NewBadRequest("Invalid page number", err))
			return
		}
	}
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			filter.PageSize = pageSize
		} else {
			h.handleError(w, errors.NewBadRequest("Invalid page size", err))
			return
		}
	}

	if fromDateStr := r.URL.Query().Get("from_date"); fromDateStr != "" {
		if fromDate, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			filter.FromDate = &fromDate
		} else {
			h.handleError(w, errors.NewBadRequest("Invalid from_date format", err))
			return
		}
	}
	if toDateStr := r.URL.Query().Get("to_date"); toDateStr != "" {
		if toDate, err := time.Parse("2006-01-02", toDateStr); err == nil {
			filter.ToDate = &toDate
		} else {
			h.handleError(w, errors.NewBadRequest("Invalid to_date format", err))
			return
		}
	}

	response, err := h.service.GetSongs(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

func (h *SongHandler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.handleError(w, errors.NewValidation("json encode error", err))
	}
}

// @Summary Get song lyrics
// @Description Get song lyrics with pagination by verses
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} models.LyricsResponse
// @Router /songs/{id}/lyrics [get]
func (h *SongHandler) GetLyrics(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling GetLyrics request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid song ID", err))
		return
	}

	page, pageSize := 0, 10
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil {
			page = parsedPage
		} else {
			h.handleError(w, errors.NewBadRequest("Invalid page query", err))
			return
		}
	}
	if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil {
			pageSize = parsedPageSize
		} else {
			h.handleError(w, errors.NewBadRequest("Invalid page size query", err))
			return
		}
	}

	response, err := h.service.GetLyrics(r.Context(), id, page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, response)
}

// @Summary Create new song
// @Description Create a new song with information from external API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.SongRequest true "Song information"
// @Success 201 {object} models.Song
// @Router /songs [post]
func (h *SongHandler) CreateSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling CreateSong request")

	var req models.SongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid request body", err))
		return
	}

	song, err := h.service.CreateSong(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusCreated, song)
}

// @Summary Update song
// @Description Update existing song information
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.SongRequest true "Song information"
// @Success 200 {object} models.Song
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling UpdateSong request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid song ID", err))
		return
	}

	var req models.SongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid request body", err))
		return
	}

	song, err := h.service.UpdateSong(r.Context(), id, &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondWithJSON(w, http.StatusOK, song)
}

// @Summary Delete song
// @Description Delete a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("Handling DeleteSong request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.handleError(w, errors.NewBadRequest("Invalid song ID", err))
		return
	}

	if err := h.service.DeleteSong(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
