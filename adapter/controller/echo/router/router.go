package router

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	"go-todo-app-clean-arch/adapter/controller/echo/custommiddleware"
	"go-todo-app-clean-arch/adapter/controller/echo/handler"
	"go-todo-app-clean-arch/adapter/controller/echo/presenter"
	"go-todo-app-clean-arch/adapter/gateway"
	"go-todo-app-clean-arch/pkg"
	"go-todo-app-clean-arch/pkg/logger"
	"go-todo-app-clean-arch/usecase"
)

// TemplateRendererはEcho用のHTMLテンプレートレンダー
type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

// Swagger の設定
func setupSwagger(router *echo.Echo) (*openapi3.T, error) {
	swagger, err := presenter.GetSwagger()
	if err != nil {
		return nil, err
	}

	env := pkg.GetEnvDefault("APP_ENV", "development")
	if env == "development" {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		router.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	return swagger, nil
}

// Echo 用のルータを作成。
func NewEchoRouter(db *gorm.DB) *echo.Echo {
	router := echo.New()

	// ミドルウェア設定
	router.Use(custommiddleware.CustomRequestLogger())
	router.Use(custommiddleware.CustomRecovery())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		// AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowCredentials: true,
	}))
	router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteNoneMode,
		// CookieSameSite: http.SameSiteDefaultMode,
		// CookieMaxAge:   60,
	}))

	// Swagger の設定
	_, err := setupSwagger(router)
	if err != nil {
		logger.Warn("Swagger setup error: " + err.Error())
	}

	// テンプレート設定
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./adapter/presenter/html/*")),
	}
	router.Renderer = renderer

	// リポジトリとユースケースの設定
	taskRepository := gateway.NewTaskRepository(db)
	taskUseCase := usecase.NewTaskUseCase(taskRepository)
	taskHandler := handler.NewTaskHandler(taskUseCase)

	userRepository := gateway.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)

	// ユーザー用エンドポイント
	users := router.Group("/api/v1/users")
	users.Use(custommiddleware.JWTMiddleware())
	users.GET("", userHandler.GetCurrentUser)
	users.DELETE("", userHandler.DeleteUser)

	// 認証用エンドポイント
	auth := router.Group("/api/v1/auth")
	auth.POST("/login", userHandler.Login)
	auth.POST("/signup", userHandler.Signup)
	auth.POST("/logout", userHandler.Logout)
	auth.GET("/csrf", userHandler.CsrfToken)

	// 認証が必要なタスク用エンドポイント
	tasks := router.Group("/api/v1/tasks")
	tasks.Use(custommiddleware.JWTMiddleware())	
	tasks.POST("", taskHandler.CreateTask)
	tasks.GET("", taskHandler.GetAllTasks)
	tasks.GET("/:id", taskHandler.GetTaskById)
	tasks.PUT("/:id", taskHandler.UpdateTaskById)
	tasks.DELETE("/:id", taskHandler.DeleteTaskById)

	// Swagger やその他のルート
	router.GET("/", handler.Index)
	router.GET("/health", handler.Health)

	return router
}
