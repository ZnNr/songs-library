package models

import "time"

// Song модель песни в базе данных
type Song struct {
	ID          int       `json:"id" db:"id"`
	GroupName   string    `json:"group_name" db:"group_name"`
	SongName    string    `json:"song_name" db:"song_name"`
	ReleaseDate time.Time `json:"release_date" db:"release_date"`
	Text        string    `json:"text" db:"text"`
	Link        string    `json:"link" db:"link"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SongRequest структура запроса для создания/обновления песни
type SongRequest struct {
	GroupName string `json:"group" binding:"required"`
	SongName  string `json:"song" binding:"required"`
	Text      string `json:"text"`
	Link      string `json:"link"`
}

// SongFilter структура фильтрации песен
type SongFilter struct {
	GroupName string     `json:"group_name"`
	SongName  string     `json:"song_name"`
	FromDate  *time.Time `json:"from_date"`
	ToDate    *time.Time `json:"to_date"`
	Text      string     `json:"text"`
	Link      string     `json:"link"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
}

// SongsResponse структура ответа со списком песен и информацией о пагинации
type SongsResponse struct {
	Songs       []Song `json:"songs"`
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
	TotalItems  int    `json:"total_items"`
	PageSize    int    `json:"page_size"`
}

// LyricsResponse структура ответа с куплетами и информацией о пагинации
type LyricsResponse struct {
	Text        string `json:"text"`
	CurrentPage int    `json:"current_page"`
	TotalPages  int    `json:"total_pages"`
	PageSize    int    `json:"page_size"`
}
