package web

import (
	"echo-demo/api/model"
	"net/http"

	"github.com/labstack/echo"
)

type CreateUserRequest struct {
	// allow empty
	UserName string `json:"username" validate:"required"`
}

type UpdateUserRequest struct {
	// allow empty
	UserName *string `json:"username"`
}

type UserResponse struct {
	ID       string `json:"id"`
	UserName string `json:"username"`
}

func NewUserResponse(user model.User) UserResponse {
	result := UserResponse{
		ID:       user.PublicID,
		UserName: user.Username,
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

// GET /users/
func (api *API) GetAllUser(ctx echo.Context) error {
	// TODO:query
	var users []model.User
	return ctx.JSON(http.StatusOK, NewUsersResponse(users))
}

// POST /users/
func (api *API) CreateUser(ctx echo.Context) error {
	var data CreateUserRequest
	if err := ctx.Bind(data); err != nil {
		return echo.ErrBadRequest
	}

	id, err := api.Repo.CreateUser(data.UserName)
	if err != nil {
		return echo.ErrInternalServerError
	}

	u, err := api.Repo.GetUserByID(id)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusCreated, NewUserResponse(u))
}

// GET /users/{id}
func (api *API) GetUser(ctx echo.Context) error {
	id, _ := ctx.Get(UserIDCtxKey).(string)

	u, err := api.Repo.GetUserByPublicID(id)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return ctx.JSON(http.StatusOK, NewUserResponse(u))
}

// PATCH /users/{id}
func (api *API) UpdateUser(ctx echo.Context) error {
	var data UpdateUserRequest
	if err := ctx.Bind(data); err != nil {
		return echo.ErrBadRequest
	}

	id, _ := ctx.Get(UserIDCtxKey).(string)
	u, err := api.Repo.GetUserByPublicID(id)
	if err != nil {
		return echo.ErrBadRequest
	}

	if data.UserName != nil {
		u.Username = *data.UserName
	}

	_, err = api.Repo.UpdateUserByID(u)
	if err != nil {
		return echo.ErrInternalServerError
	}

	u, err = api.Repo.GetUserByPublicID(id)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusOK, NewUserResponse(u))
}

// DELETE /users/{id}
func (api *API) DeleteUser(ctx echo.Context) error {
	id, _ := ctx.Get(UserIDCtxKey).(string)

	u, err := api.Repo.GetUserByPublicID(id)
	if err != nil {
		return echo.ErrBadRequest
	}

	_, err = api.Repo.DeleteUserByID(u.ID)
	if err != nil {
		return echo.ErrInternalServerError
	}
	return ctx.NoContent(http.StatusNoContent)
}

const UserIDCtxKey string = "userID"

func RequireParamUserID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id := ctx.Param("id")
		if id == "" {
			return echo.ErrBadRequest
		}
		ctx.Set(UserIDCtxKey, id)
		return nil
	}
}
