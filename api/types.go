package api

const internalServerError = "internal error"

type ErrorResponse struct {
	Error string `json:"error"`
}
