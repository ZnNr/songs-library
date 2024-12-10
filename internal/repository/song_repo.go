package repository

import (
	"context"
	"github.com/ZnNr/songs-library/internal/models"
)

// SongRepository  методы для взаимодействия с данными песен в базе данных
type SongRepository interface {
	GetSongs(ctx context.Context, filter *models.SongFilter) (*models.SongsResponse, error)
	GetSongByID(ctx context.Context, id int) (*models.Song, error)
	CreateSong(ctx context.Context, song *models.Song) (*models.Song, error)
	UpdateSong(ctx context.Context, song *models.Song) (*models.Song, error)
	DeleteSong(ctx context.Context, id int) error
}
