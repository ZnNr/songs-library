// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "Get list of songs with optional filtering and pagination",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get songs with filtering and pagination",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Group name",
                        "name": "group_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Song name",
                        "name": "song_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "From date (format: 2006-01-02)",
                        "name": "from_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "To date (format: 2006-01-02)",
                        "name": "to_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Text content",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Link",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.SongsResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new song with information from external API",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Create new song",
                "parameters": [
                    {
                        "description": "Song information",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SongRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                }
            }
        },
        "/songs/{id}": {
            "put": {
                "description": "Update existing song information",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Update song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Song information",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SongRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a song by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Delete song",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/songs/{id}/lyrics": {
            "get": {
                "description": "Get song lyrics with pagination by verses",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Get song lyrics",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Song ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Page number",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page size",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.LyricsResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.LyricsResponse": {
            "type": "object",
            "properties": {
                "current_page": {
                    "description": "Номер текущей страницы",
                    "type": "integer"
                },
                "page_size": {
                    "description": "Количество элементов на странице",
                    "type": "integer"
                },
                "text": {
                    "description": "Текст песни или куплетов",
                    "type": "string"
                },
                "total_pages": {
                    "description": "Общее количество страниц",
                    "type": "integer"
                }
            }
        },
        "models.Song": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "Дата и время создания записи",
                    "type": "string"
                },
                "group_name": {
                    "description": "Название группы или исполнителя",
                    "type": "string"
                },
                "id": {
                    "description": "Уникальный идентификатор песни",
                    "type": "integer"
                },
                "link": {
                    "description": "Ссылка на песню (например, на YouTube)",
                    "type": "string"
                },
                "release_date": {
                    "description": "Дата выпуска песни",
                    "type": "string"
                },
                "song_name": {
                    "description": "Название песни",
                    "type": "string"
                },
                "text": {
                    "description": "Текст песни",
                    "type": "string"
                },
                "updated_at": {
                    "description": "Дата и время последнего обновления записи",
                    "type": "string"
                }
            }
        },
        "models.SongRequest": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "description": "Название группы, обязательное поле",
                    "type": "string"
                },
                "link": {
                    "description": "Ссылка на песню, необязательное поле",
                    "type": "string"
                },
                "song": {
                    "description": "Название песни, обязательное поле",
                    "type": "string"
                },
                "text": {
                    "description": "Текст песни, необязательное поле",
                    "type": "string"
                }
            }
        },
        "models.SongsResponse": {
            "type": "object",
            "properties": {
                "page": {
                    "description": "Номер текущей страницы",
                    "type": "integer"
                },
                "page_size": {
                    "description": "Количество элементов на странице",
                    "type": "integer"
                },
                "songs": {
                    "description": "Список песен",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Song"
                    }
                },
                "total_items": {
                    "description": "Общее количество песен",
                    "type": "integer"
                },
                "total_pages": {
                    "description": "Общее количество страниц",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Music Library API",
	Description:      "API for managing music library",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
