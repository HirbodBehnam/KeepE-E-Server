package database

import "time"

type SavedImage struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type Color struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

type Note struct {
	ID              int32        `json:"id"`
	Title           string       `json:"title" binding:"required"`
	Text            string       `json:"text"`
	CreatedAt       time.Time    `json:"created_at"`
	Deadline        *time.Time   `json:"deadline,omitempty"`
	BackgroundColor *Color       `json:"background_color,omitempty"`
	Images          []SavedImage `json:"images"`
}

type TodoItem struct {
	Text     string     `json:"text" binding:"required"`
	Deadline *time.Time `json:"deadline,omitempty"`
	Done     bool       `json:"done"`
}

type Todo struct {
	ID        int32      `json:"id"`
	Title     string     `json:"title" binding:"required"`
	CreatedAt time.Time  `json:"created_at"`
	Tasks     []TodoItem `json:"tasks"`
}
