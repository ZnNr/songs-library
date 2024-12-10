package models

import "time"

// Song представляет модель песни в базе данных.
type Song struct {
	ID          int       `json:"id" db:"id"`                     // Уникальный идентификатор песни
	GroupName   string    `json:"group_name" db:"group_name"`     // Название группы или исполнителя
	SongName    string    `json:"song_name" db:"song_name"`       // Название песни
	ReleaseDate time.Time `json:"release_date" db:"release_date"` // Дата выпуска песни
	Text        string    `json:"text" db:"text"`                 // Текст песни
	Link        string    `json:"link" db:"link"`                 // Ссылка на песню (например, на YouTube)
	CreatedAt   time.Time `json:"created_at" db:"created_at"`     // Дата и время создания записи
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`     // Дата и время последнего обновления записи
}

// SongRequest представляет структуру запроса для создания или обновления песни.
type SongRequest struct {
	GroupName string `json:"group" binding:"required"` // Название группы, обязательное поле
	SongName  string `json:"song" binding:"required"`  // Название песни, обязательное поле
	Text      string `json:"text"`                     // Текст песни, необязательное поле
	Link      string `json:"link"`                     // Ссылка на песню, необязательное поле
}

// SongFilter представляет структуру фильтрации песен.
type SongFilter struct {
	GroupName string     `json:"group_name"` // Название группы для фильтрации
	SongName  string     `json:"song_name"`  // Название песни для фильтрации
	FromDate  *time.Time `json:"from_date"`  // Дата начала фильтрации (включительно)
	ToDate    *time.Time `json:"to_date"`    // Дата окончания фильтрации (включительно)
	Text      string     `json:"text"`       // Текст песни для фильтрации
	Link      string     `json:"link"`       // Ссылка на песню для фильтрации
	Page      int        `json:"page"`       // Номер текущей страницы
	PageSize  int        `json:"page_size"`  // Размер страницы (количество элементов на странице)
}

// SongsResponse представляет структуру ответа со списком песен и информацией о пагинации.
type SongsResponse struct {
	Songs      []Song `json:"songs"`       // Список песен
	Page       int    `json:"page"`        // Номер текущей страницы
	TotalPages int    `json:"total_pages"` // Общее количество страниц
	TotalItems int    `json:"total_items"` // Общее количество песен
	PageSize   int    `json:"page_size"`   // Количество элементов на странице
}

// LyricsResponse представляет структуру ответа с текстом куплетов и информацией о пагинации.
type LyricsResponse struct {
	Text        string `json:"text"`         // Текст песни или куплетов
	CurrentPage int    `json:"current_page"` // Номер текущей страницы
	TotalPages  int    `json:"total_pages"`  // Общее количество страниц
	PageSize    int    `json:"page_size"`    // Количество элементов на странице
}
