package usecase

import (
	"errors"
	"go-todo-app-clean-arch/adapter/gateway"
	"go-todo-app-clean-arch/entity"
	"go-todo-app-clean-arch/pkg/logger"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	GetCurrentUser(userId int) (*entity.User, error)
	DeleteUser(userId int) error
	Signup(user *entity.User) (*entity.User, error)
	Login(credentials *entity.Credentials) (string, error)
}

type userUseCase struct {
	userRepository gateway.UserRepository
}

func NewUserUseCase(userRepository gateway.UserRepository) *userUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

// func (u *userUseCase) Create(user *entity.User) (*entity.User, error) {
// 	return u.userRepository.Create(user)
// }

func (u *userUseCase) GetCurrentUser(userId int) (*entity.User, error) {
	return u.userRepository.GetCurrentUser(userId)
}

func (u *userUseCase) DeleteUser(userId int) error {
	return u.userRepository.DeleteUser(userId)
}

func (u *userUseCase) Signup(user *entity.User) (*entity.User, error) {
	// パスワードをハッシュ化
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	// ユーザー作成
	return u.userRepository.Signup(user)
}

func (u *userUseCase) Login(credentials *entity.Credentials) (string, error) {
	// メールアドレスでユーザーを検索
	// TODO: credentialsではなく普通にuserを使用した方が余計な処理が減るかも
	user, err := u.userRepository.FindByEmail(credentials.Email)
	if err != nil || !CheckPasswordHash(credentials.Password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// ペイロードの作成
	claims := &jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	// トークン生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("SECRET")
	logger.Info("os.Getenv(SECRET): " + secret)

	// トークンに署名を付与
	tokenString, err := token.SignedString([]byte(secret))
	logger.Info("tokenString: " + tokenString)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
