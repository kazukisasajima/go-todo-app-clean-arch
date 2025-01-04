package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"go-todo-app-clean-arch/pkg/logger"
)

// カラーコード
var (
	green  = "\033[32m"
	blue   = "\033[34m"
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

// リクエスト情報をログ出力するミドルウェア
func EchoZap() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			// ステータスコードに応じた色付け
			status := c.Response().Status
			statusColor := getStatusColor(status)

			// HTTPメソッドに応じた色付け
			method := c.Request().Method
			methodColor := getMethodColor(method)

			// ログの出力
			logger.ZapLogger.Info("Request",
				zap.String("method", method),
				zap.String("path", c.Request().URL.Path),
				zap.Int("status", status),
				zap.String("latency", stop.Sub(start).String()),
				zap.String("client_ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			)

			// カラフルな出力
			fmt.Printf("[%s%s%s] %s%3d%s | %13v | %15s | %s%-7s%s %s\n",
				blue, time.Now().Format("2006/01/02 - 15:04:05"), reset, // 日時
				statusColor, status, reset,                             // ステータスコードと色
				stop.Sub(start),                                        // レイテンシ
				c.RealIP(),                                            // クライアントIP
				methodColor, method, reset,                            // HTTPメソッドと色
				c.Request().URL.Path,                                  // リクエストパス
			)

			return err
		}
	}
}

// RecoveryWithZap はパニックをキャッチしてログを出力するミドルウェア
func RecoveryWithZap() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.ZapLogger.Error("Panic recovered",
						zap.Any("error", r),
					)
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

// HTTPステータスコードに応じた色を返却。
func getStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return green
	case status >= 300 && status < 400:
		return blue
	case status >= 400 && status < 500:
		return yellow
	default:
		return red
	}
}

// HTTPメソッドに応じた色を返却
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return green
	case "PUT":
		return yellow
	case "DELETE":
		return red
	default:
		return reset
	}
}
