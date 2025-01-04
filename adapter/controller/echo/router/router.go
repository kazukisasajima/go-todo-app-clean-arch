package router

import (
	"encoding/json"
	"html/template"
	"io"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	"go-todo-app-clean-arch/adapter/controller/echo/handler"
	"go-todo-app-clean-arch/adapter/controller/echo/middleware"
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

	router.Use(middleware.EchoZap())
	router.Use(middleware.RecoveryWithZap())

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

	// 共通ハンドラーの登録
	taskRepository := gateway.NewTaskRepository(db)
	taskUseCase := usecase.NewTaskUseCase(taskRepository)
	taskHandler := handler.NewTaskHandler(taskUseCase)

	userRepository := gateway.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository)
	userHandler := handler.NewUserHandler(userUseCase)

	serverHandler := handler.NewHandler().
		Register(taskHandler).
		Register(userHandler)

	// ルート定義
	router.GET("/", handler.Index)
	router.GET("/health", handler.Health)

	// API エンドポイント
	apiGroup := router.Group("/api/v1")
	{
		presenter.RegisterHandlers(apiGroup, serverHandler)
	}

	return router
}
