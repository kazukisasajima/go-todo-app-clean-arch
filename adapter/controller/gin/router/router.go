package router

import (
	"encoding/json"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	ginMiddleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"

	"go-todo-app-clean-arch/adapter/controller/gin/handler"
	"go-todo-app-clean-arch/adapter/controller/gin/middleware"
	"go-todo-app-clean-arch/adapter/controller/gin/presenter"
	"go-todo-app-clean-arch/adapter/gateway"
	"go-todo-app-clean-arch/pkg"
	"go-todo-app-clean-arch/pkg/logger"
	"go-todo-app-clean-arch/usecase"
)

// Swaggerの設定
func setupSwagger(router *gin.Engine) (*openapi3.T, error) {
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
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	return swagger, nil
}

func NewGinRouter(db *gorm.DB, corsAllowOrigins []string) (*gin.Engine, error) {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	swagger, err := setupSwagger(router)
	if err != nil {
		logger.Warn(err.Error())
		return nil, err
	}

	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())

	// ViewのHTMLの設定
	router.LoadHTMLGlob("./adapter/presenter/html/*")
	router.GET("/", handler.Index)

	// Healthチェック用のAPI
	router.GET("/health", handler.Health)

	apiGroup := router.Group("/api")
	{
		apiGroup.Use(middleware.TimeoutMiddleware(2 * time.Second))
		v1 := apiGroup.Group("/v1")
		{
			v1.Use(ginMiddleware.OapiRequestValidator(swagger))
			taskRepository := gateway.NewTaskRepository(db)
			taskUseCase := usecase.NewTaskUseCase(taskRepository)
			taskHandler := handler.NewTaskHandler(taskUseCase)
			userRepository := gateway.NewUserRepository(db)
			userUseCase := usecase.NewUserUseCase(userRepository)
			userHandler := handler.NewUserHandler(userUseCase)
			presenter.RegisterHandlers(v1,
				handler.NewHandler().
					Register(taskHandler).
					Register(userHandler))
		}
	}
	return router, err
}
