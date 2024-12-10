package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ZnNr/songs-library/internal/errors"
	"github.com/ZnNr/songs-library/internal/models"
	"github.com/ZnNr/songs-library/internal/repository"
)

// SQL Queries
const (
	addSongQuery = `
	INSERT INTO songs (group_name, song_name, release_date, text, link)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, group_name, song_name, release_date, text, link, created_at, updated_at`

	getAllSongsQuery = `
		SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at
		FROM songs
		WHERE ($1 = '' OR group_name ILIKE '%' || $1 || '%')
		AND ($2 = '' OR song_name ILIKE '%' || $2 || '%')
		AND ($3::timestamp IS NULL OR release_date >= $3)
		AND ($4::timestamp IS NULL OR release_date <= $4)
		AND ($5 = '' OR text ILIKE '%' || $5 || '%')
		AND ($6 = '' OR link ILIKE '%' || $6 || '%')
		ORDER BY created_at DESC
		LIMIT $7 OFFSET $8`

	// queries счетчик количества песен с фильтрами
	countSongsQuery = `
		SELECT COUNT(*)
		FROM songs
		WHERE ($1 = '' OR group_name ILIKE '%' || $1 || '%')
		AND ($2 = '' OR song_name ILIKE '%' || $2 || '%')
		AND ($3::timestamp IS NULL OR release_date >= $3)
		AND ($4::timestamp IS NULL OR release_date <= $4)
		AND ($5 = '' OR text ILIKE '%' || $5 || '%')
		AND ($6 = '' OR link ILIKE '%' || $6 || '%')`

	// queries получить песню по id
	getSongByIDQuery = `
		SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at
		FROM songs
		WHERE id = $1`

	// update обновить песню
	updateSongQuery = `
		UPDATE songs
		SET group_name = $1, 
			song_name = $2, 
			release_date = $3,
			text = $4,
			link = $5,
			updated_at = NOW()
		WHERE id = $6
		RETURNING id, group_name, song_name, release_date, text, link, created_at, updated_at`

	// delete удалить песню по id
	deleteSongQuery = `DELETE FROM songs WHERE id = $1`

	// queries проверить существование песни
	checkSongExistsQuery = `
		SELECT EXISTS(
			SELECT 1 FROM songs 
			WHERE group_name = $1 
			AND song_name = $2 
			AND id != $3
		)`
)

// PostgresSongRepository имплементирует SongRepository для PostgreSQL.
type PostgresSongRepository struct {
	db *sql.DB
}

func NewPostgresSongRepository(db *sql.DB) repository.SongRepository {
	return &PostgresSongRepository{db: db}
}

// CreateSong создает новую песню
func (r *PostgresSongRepository) CreateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	if exists, err := r.songExists(ctx, song.GroupName, song.SongName, 0); err != nil || exists {
		if err != nil {
			return nil, errors.NewInternal("failed to check song existence", err)
		}
		return nil, errors.NewAlreadyExists("song already exists", nil)
	}
	return song, r.insertSong(ctx, song)
}

// songExists проверяет, существует ли песня с указанным названием и группой.
func (r *PostgresSongRepository) songExists(ctx context.Context, groupName, songName string, songID int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, checkSongExistsQuery, groupName, songName, songID).Scan(&exists)
	if err != nil {
		return false, errors.NewInternal("failed to check song existence", err)
	}
	return exists, nil
}

// insertSong вставляет новую песню в базу данных.
func (r *PostgresSongRepository) insertSong(ctx context.Context, song *models.Song) error {
	return r.db.QueryRowContext(
		ctx,
		addSongQuery,
		song.GroupName,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(
		&song.ID,
		&song.GroupName,
		&song.SongName,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt,
	)
}

// GetSongs получает список песен с учетом фильтров и постраничной навигации
func (r *PostgresSongRepository) GetSongs(ctx context.Context, filter *models.SongFilter) (*models.SongsResponse, error) {
	// Устанавливаем значения по умолчанию
	setDefaultFilterValues(filter)

	// Получаем общее количество записей
	totalItems, err := r.countTotalSongs(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Вычисляем общее количество страниц
	totalPages := (totalItems + filter.PageSize - 1) / filter.PageSize

	// Проверка существования запрашиваемой страницы
	if filter.Page > totalPages {
		return nil, errors.NewNotFound(fmt.Sprintf("page %d does not exist, total pages: %d", filter.Page, totalPages), nil)
	}

	offset := (filter.Page - 1) * filter.PageSize

	// Получаем записи для текущей страницы
	songs, err := r.getSongsByPage(ctx, filter, offset)
	if err != nil {
		return nil, err
	}

	// Формируем ответ
	return &models.SongsResponse{
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		Songs:      songs,
	}, nil
}

// setDefaultFilterValues устанавливает значения по умолчанию для фильтра
func setDefaultFilterValues(filter *models.SongFilter) {
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}
}

// countTotalSongs получает общее количество песен, соответствующих фильтру
func (r *PostgresSongRepository) countTotalSongs(ctx context.Context, filter *models.SongFilter) (int, error) {
	var totalItems int
	err := r.db.QueryRowContext(
		ctx,
		countSongsQuery,
		filter.GroupName,
		filter.SongName,
		filter.FromDate,
		filter.ToDate,
		filter.Text,
		filter.Link).Scan(&totalItems)
	if err != nil {
		return 0, errors.NewInternal("failed to count songs", err)
	}
	return totalItems, nil
}

// getSongsByPage получает список песен для заданной страницы
func (r *PostgresSongRepository) getSongsByPage(ctx context.Context, filter *models.SongFilter, offset int) ([]models.Song, error) {
	rows, err := r.db.QueryContext(ctx, getAllSongsQuery,
		filter.GroupName,
		filter.SongName,
		filter.FromDate,
		filter.ToDate,
		filter.Text,
		filter.Link,
		filter.PageSize,
		offset,
	)
	if err != nil {
		return nil, errors.NewInternal("failed to query songs", err)
	}
	defer rows.Close()

	// Собираем список песен
	var songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(
			&song.ID,
			&song.GroupName,
			&song.SongName,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
			&song.CreatedAt,
			&song.UpdatedAt,
		); err != nil {
			return nil, errors.NewInternal("failed to scan song", err)
		}
		songs = append(songs, song)
	}

	// Проверка на ошибки после завершения перебора строк
	if err := rows.Err(); err != nil {
		return nil, errors.NewInternal("error occurred while iterating over songs", err)
	}

	return songs, nil
}

// GetSongByID запрашивает информацию о song по ее ID из базы данных PostgreSQL
func (r *PostgresSongRepository) GetSongByID(ctx context.Context, id int) (*models.Song, error) {
	var song models.Song
	err := r.db.QueryRowContext(ctx, getSongByIDQuery, id).Scan(&song.ID,
		&song.GroupName,
		&song.SongName,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
		&song.CreatedAt,
		&song.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFound("song not found", err)
	} else if err != nil {
		return nil, errors.NewInternal("failed to get song", err)
	}
	return &song, nil
}

// updateSong обновляет информацию о песне в базе данных.
func (r *PostgresSongRepository) updateSong(ctx context.Context, song *models.Song) error {
	err := r.db.QueryRowContext(ctx, updateSongQuery, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link, song.ID).Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Link, &song.CreatedAt, &song.UpdatedAt)
	if err == sql.ErrNoRows {
		return errors.NewNotFound("song not found", err)
	} else if err != nil {
		return errors.NewInternal("failed to update song", err)
	}
	return nil
}

// UpdateSong обновляет информацию о песне в базе данных.
func (r *PostgresSongRepository) UpdateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	if exists, err := r.songExists(ctx, song.GroupName, song.SongName, song.ID); err != nil || exists {
		if err != nil {
			return nil, err
		}
		return nil, errors.NewAlreadyExists("song already exists", nil)
	}
	if err := r.updateSong(ctx, song); err != nil {
		return nil, err
	}
	return song, nil
}

// DeleteSong удаляет песню по заданному идентификатору.
func (r *PostgresSongRepository) DeleteSong(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, deleteSongQuery, id)
	if err != nil {
		return errors.NewInternal("failed to execute delete query", err)
	}

	if rowsAffected, err := result.RowsAffected(); err != nil {
		return errors.NewInternal("failed to retrieve affected rows after delete", err)
	} else if rowsAffected == 0 {
		return errors.NewNotFound("song not found", nil)
	}
	return nil
}
