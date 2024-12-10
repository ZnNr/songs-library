package service

import (
	"context"
	"strings"
	"time"

	"github.com/ZnNr/songs-library/internal/errors"
	"github.com/ZnNr/songs-library/internal/models"
	"github.com/ZnNr/songs-library/internal/repository"
	"go.uber.org/zap"
)

type SongService struct {
	repo   repository.SongRepository
	logger *zap.Logger
}

func NewSongService(repo repository.SongRepository, logger *zap.Logger) *SongService {
	return &SongService{
		repo:   repo,
		logger: logger,
	}
}

// GetSongs получает список песен с опциональным фильтром.
func (s *SongService) GetSongs(ctx context.Context, filter *models.SongFilter) (*models.SongsResponse, error) {
	s.logger.Info("Getting songs with filter",
		zap.String("group", filter.GroupName),
		zap.String("song", filter.SongName),
		zap.Any("fromDate", filter.FromDate),
		zap.Any("toDate", filter.ToDate),
		zap.String("text", filter.Text),
		zap.String("link", filter.Link),
		zap.Int("page", filter.Page),
		zap.Int("pageSize", filter.PageSize))

	// Валидация параметров фильтра
	if err := validateFilter(filter); err != nil {
		s.logger.Warn("Invalid filter", zap.Error(err))
		return nil, err
	}

	return s.repo.GetSongs(ctx, filter)
}

// GetLyrics получает текст песни с опциональным фильтром.
func (s *SongService) GetLyrics(ctx context.Context, id, page, pageSize int) (*models.LyricsResponse, error) {
	s.logger.Info("Getting lyrics",
		zap.Int("songId", id),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	song, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFound("song not found", err)
	}

	if song.Text == "" {
		return nil, errors.NewNotFound("lyrics not found", nil)
	}

	verses := strings.Split(song.Text, "\n\n")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Подсчет общих страниц
	totalPages := (len(verses) + pageSize - 1) / pageSize
	if page > totalPages {
		return nil, errors.NewNotFound("page out of range", nil)
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(verses) {
		end = len(verses)
	}

	return &models.LyricsResponse{
		Text:        strings.Join(verses[start:end], "\n\n"),
		CurrentPage: page,
		TotalPages:  totalPages,
		PageSize:    pageSize,
	}, nil
}

// CreateSong создает новую песню.
func (s *SongService) CreateSong(ctx context.Context, req *models.SongRequest) (*models.Song, error) {
	s.logger.Info("Creating new song",
		zap.String("group", req.GroupName),
		zap.String("song", req.SongName))

	if err := validateSongRequest(req); err != nil {
		return nil, err
	}

	song := &models.Song{
		GroupName:   req.GroupName,
		SongName:    req.SongName,
		ReleaseDate: time.Now(),
		Text:        req.Text,
		Link:        req.Link,
	}

	return s.repo.CreateSong(ctx, song)
}

// UpdateSong обновляет существующую песню.
func (s *SongService) UpdateSong(ctx context.Context, id int, req *models.SongRequest) (*models.Song, error) {
	s.logger.Info("Updating song",
		zap.Int("id", id),
		zap.String("group", req.GroupName),
		zap.String("song", req.SongName))

	song, err := s.repo.GetSongByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFound("song not found", err)
	}

	// Обновляем только предоставленные поля
	updateSongFields(song, req)
	song.UpdatedAt = time.Now()

	return s.repo.UpdateSong(ctx, song)
}

// DeleteSong удаляет существующую песню.
func (s *SongService) DeleteSong(ctx context.Context, id int) error {
	s.logger.Info("Deleting song", zap.Int("id", id))
	return s.repo.DeleteSong(ctx, id)
}

// validateFilter выполняет проверку валидации фильтра песен.
func validateFilter(filter *models.SongFilter) error {
	// Здесь можно добавить логику валидации
	return nil // Замените на реальную логику валидации
}

// validateSongRequest выполняет проверку валидности запроса на создание или обновление песни.
func validateSongRequest(req *models.SongRequest) error {
	if req.SongName == "" {
		return errors.NewValidation("song name cannot be empty", nil)
	}
	// Добавить дополнительные проверки при необходимости
	return nil
}

// updateSongFields обновляет поля песни на основании запроса.
func updateSongFields(song *models.Song, req *models.SongRequest) {
	if req.GroupName != "" {
		song.GroupName = req.GroupName
	}
	if req.SongName != "" {
		song.SongName = req.SongName
	}
	if req.Text != "" {
		song.Text = req.Text
	}
	if req.Link != "" {
		song.Link = req.Link
	}
}
