package database

import (
	"database/sql"

	"fmt"
	"github.com/ZnNr/songs-library/internal/errors"
	"github.com/ZnNr/songs-library/internal/models"
)

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

	// delete удалить песню
	deleteSongQuery = `DELETE FROM songs WHERE id = $1`

	// queries проверить существование песни
	checkSongExistsQuery = `
		SELECT EXISTS(
			SELECT 1 FROM songs 
			WHERE group_name = $1 
			AND song_name = $2 
			AND id != $3
		)`

	// queries проверить существование песни
	checkSongExistsForCreateQuery = `
		SELECT EXISTS(
			SELECT 1 FROM songs 
			WHERE group_name = $1 
			AND song_name = $2
		)`
)

// AddSong сохраняет список элементов заказа в БД, пропуская существующие элементы
func AddSongs(db *sql.DB, song *[]models.Song) (*[]models.Song, error) {
	for _, song := range song {
		exists, err := SongExists(db, song.GroupName, song.SongName) // Проверка существования
		if err != nil {
			return nil, errors.NewInternal("failed to check song existence", err)
		}
		if exists {
			return nil, errors.NewAlreadyExists("song with this group name and song name already exists", nil)
		}
		if !exists {
			err = AddSong(db, song)
			if err != nil {
				return nil, errors.NewInternal("failed to create song", err)
			}
		}
	}
	return song, nil
}

// ItemExists проверяет, существует ли элемент в БД
func SongExists(db *sql.DB, GroupName string, SongName string) (bool, error) {
	var exists bool
	err := db.QueryRow(checkSongExistsForCreateQuery, GroupName, SongName).Scan(&exists)
	if err != nil {
		return false, errors.NewInternal("failed to check song existence", err)
	}
	if exists {
		return true, errors.NewAlreadyExists("song with this group name and song name already exists", nil)
	}
	return exists, nil
}

// AddItem добавляет новый элемент в БД
func AddSong(db *sql.DB, song models.Song) error {
	_, err := db.Exec(
		addSongQuery,
		song.GroupName,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
	)
	if err != nil {
		return fmt.Errorf("failed to execute add item query: %w", err)
	}
	return nil
}

// GetItems получает все элементы из БД по идентификатору заказа
func GetItems(db *sql.DB, orderUID string) ([]models.Item, error) {
	rows, err := db.Query(getAllItemsQuery, orderUID)
	if err != nil {
		return nil, fmt.Errorf("get items failed: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iteration over rows failed: %w", err)
	}

	if len(items) == 0 {
		// Логируем ненайденный orderUID
		fmt.Printf("No items found for order UID %s\n", orderUID)
		return nil, fmt.Errorf("items not found for order UID %s", orderUID)
	}

	return items, nil
}
