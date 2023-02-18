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

type getDeleteEntityRequest struct {
	ID int32 `form:"id" binding:"required"`
}

type reOrderRequest struct {
	ID1 int32 `json:"id1" binding:"required"`
	ID2 int32 `json:"id2" binding:"required"`
}
