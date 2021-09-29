package web

import (
	"echo-demo/model"
	"net/http"

	"github.com/labstack/echo"
)

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UpdateUserRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewUserResponse(user model.User) UserResponse {
	result := UserResponse{
		ID:   user.ID,
		Name: user.Username,
	}
	return result
}

func NewUsersResponse(users []model.User) []UserResponse {
	result := []UserResponse{}
	for _, u := range users {
		result = append(result, NewUserResponse(u))
	}
	return result
}

func (api *API) GetAllUser(ctx echo.Context) error {
	// TODO:query
	var users []model.User
	return ctx.JSON(http.StatusOK, NewUsersResponse(users))
}

// POST users
func (api *API) CreateUser(ctx echo.Context) error {
	var data CreateUserRequest
	if err := ctx.Bind(data); err != nil {
		return echo.ErrBadRequest
	}
	// TODO:query
	var u model.User
	return ctx.JSON(http.StatusCreated, NewUserResponse(u))
}

// GET users/{id}
func (api *API) GetUser(ctx echo.Context) error {

	id := ctx.Param("id")
	_ = id
	// TODO:search by userID
	var user model.User
	return ctx.JSON(http.StatusOK, NewUserResponse(user))
}

// PATCH users/{id}
func (api *API) UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")
	_ = id
	// TODO:get user query
	var user model.User
	return ctx.JSON(http.StatusOK, NewUserResponse(user))
}

// DELETE users/{id}
func (api *API) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")
	_ = id
	// TODO:query
	return ctx.NoContent(http.StatusNoContent)
}
