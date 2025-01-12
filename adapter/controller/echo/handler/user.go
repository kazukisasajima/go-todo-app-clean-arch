package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"go-todo-app-clean-arch/adapter/controller/echo/presenter"
	"go-todo-app-clean-arch/entity"
	"go-todo-app-clean-arch/pkg/logger"
	"go-todo-app-clean-arch/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func userToResponse(user *entity.User) *presenter.UserResponse {
	return &presenter.UserResponse{
		Id:    user.ID,
		Email: user.Email,
	}
}

func (u *UserHandler) GetCurrentUser(c echo.Context) error {
	// JWTトークンからユーザーIDを取得
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64)) // トークンから取得した user_id

	// ユースケースからユーザー情報を取得
	userEntity, err := u.userUseCase.GetCurrentUser(userId)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: "Failed to retrieve user"})
	}

	// レスポンスとしてユーザー情報を返す
	return c.JSON(http.StatusOK, userToResponse(userEntity))
}

func (u *UserHandler) DeleteUser(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := int(claims["user_id"].(float64))

	if err := u.userUseCase.DeleteUser(userId); err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

func (u *UserHandler) Signup(c echo.Context) error {
	var requestBody presenter.CreateUserJSONRequestBody
	if err := c.Bind(&requestBody); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	user := &entity.User{
		Email:    string(requestBody.Email),
		Password: requestBody.Password,
	}

	createdUser, err := u.userUseCase.Signup(user)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, &presenter.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, userToResponse(createdUser))
}

func (u *UserHandler) Login(c echo.Context) error {
	var credentials entity.Credentials
	if err := c.Bind(&credentials); err != nil {
		logger.Warn(err.Error())
		return c.JSON(http.StatusBadRequest, &presenter.ErrorResponse{Message: err.Error()})
	}

	tokenString, err := u.userUseCase.Login(&credentials)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Set JWT token as a secure cookie
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	// return c.JSON(http.StatusOK, map[string]string{"message": "login successful"})
	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	// cookie.Secure = true
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteNoneMode
	c.SetCookie(cookie)

	// return c.JSON(http.StatusOK, map[string]string{"message": "logout successful"})
	return c.NoContent(http.StatusOK)
}

func (u *UserHandler) CsrfToken(c echo.Context) error {
	secret := os.Getenv("SECRET")
	logger.Info("secret: " + secret)

	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
