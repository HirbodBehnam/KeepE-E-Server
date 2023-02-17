package api

const internalServerErr = "internal error"
const bodyParseErr = "cannot parse body: "
const userExistsErr = "user already exists"

type errorResponse struct {
	Error string `json:"error"`
}

type signupLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}
