package model

import "time"

type CreateRunRequest struct {
	Topic string `json:"topic" binding:"required,min=1,max=500"`
}

type RunResponse struct {
	ID         string      `json:"id"`
	Topic      string      `json:"topic"`
	Status     string      `json:"status"`
	Result     interface{} `json:"result"`
	Validation interface{} `json:"validation"`
	Error      string      `json:"error,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type RunSummary struct {
	ID        string    `json:"id"`
	Topic     string    `json:"topic"`
	Status    string    `json:"status"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}