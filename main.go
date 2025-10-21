package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "saas/api/openapi"
	"saas/internal/comment"
	"saas/internal/common/logger"
	"saas/internal/common/metrics"
	"saas/internal/common/middleware/auth"
	"saas/internal/common/server"
	"saas/internal/common/uid"
	"saas/internal/common/utils"
	"saas/internal/img"
	"saas/internal/tenant"
	"saas/internal/user"
	"time"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/subosito/gotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setGDB() {
	host := utils.GetEnv("PSQL_HOST")
	port := utils.GetEnv("PSQL_PORT")
	username := utils.GetEnv("PSQL_USERNAME")
	password := utils.GetEnv("PSQL_PASSWORD")
	dbname := utils.GetEnv("PSQL_DB_NAME")
	sslmode := utils.GetEnv("PSQL_SSL_MODE")

	maxOpenConns := utils.GetEnvAsInt("DB_MAX_OPEN_CONNS")
	maxIdleConns := utils.GetEnvAsInt("DB_MAX_IDLE_CONNS")
	connMaxLifetime := time.Duration(utils.GetEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES")) * time.Minute
	connMaxIdleTime := time.Duration(utils.GetEnvAsInt("DB_CONN_MAX_IDLE_TIME_MINUTES")) * time.Minute

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, username, password, dbname, sslmode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	// 配置连接池
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	// 测试连接
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("无法连接到数据库: %v", err))
	}

	boil.SetDB(db)

	boil.DebugMode = true

	logMode := utils.GetEnv("LOG_MODE")
	if logMode != "dev" {
		if err := os.MkdirAll("./logs", 0755); err != nil {
			panic(fmt.Sprintf("创建日志目录失败:%v", err))
		}
		fh, err := os.OpenFile("./logs/sqlboiler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(fmt.Sprintf("打开debug日志错误:%v", err))
		}
		boil.DebugWriter = fh
	}
}

// 关闭http服务清理资源
func clear() {
	if err := metrics.StopPrometheusServer(); err != nil {
		log.Fatalf("Prometheus 服务成功关闭,err:%v", err)
	}
	log.Println("Prometheus 服务成功关闭")
}

// @title           自定义title
// @version         1.0
// @description     自定义描述
// @termsOfService  https://Lirous.com
// @contact.name   Lirous
// @contact.url    https://Lirous.com
// @contact.email  lirous@lirous.com
// @license.name  MIT
// @license.url   https://github.com/Lirous587/go-scaffold/main/LICENSE
// @host      localhost:8080
// @BasePath  /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
// swag init -g main.go -o ./api/openapi
func main() {
	var err error

	if err := gotenv.Load(); err != nil {
		panic(err)
	}

	uid.Init()

	setGDB()

	auth.Init()

	if err = logger.Init(); err != nil {
		panic(errors.WithMessage(err, "logger模块初始化失败"))
	}

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// 启动 Prometheus
	metricsClient := metrics.NewPrometheusClient()
	if err = metrics.StartPrometheusServer(); err != nil {
		panic(err)
	}

	// 启动 HTTP 服务器
	server.RunHttpServer(utils.GetEnv("SERVER_PORT"), metricsClient, func(r *gin.RouterGroup) {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
			ginSwagger.PersistAuthorization(true)))

		user.InitV1(r)
		img.InitV1(r)
		//captcha.InitV1(r)

		tenant.InitV1(r)
		comment.InitV1(r)
	},
		clear,
	)
}
