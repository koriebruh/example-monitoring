package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"wyw/docs"

	"wyw/entity"
	"wyw/handler"
	"wyw/metric"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	Once       sync.Once
	InstanceDB *gorm.DB
)

func getInstance() *gorm.DB {
	Once.Do(func() {
		//dsn := fmt.Sprintf("root:rootpassword@tcp(mysql:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local")
		dsn := fmt.Sprintf("root:korie123@tcp(localhost:3306)/hehey?charset=utf8mb4&parseTime=True&loc=Local")
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: nil,
		})
		if err != nil {
			log.Fatalf("got error dial mysql %v", err)
		}

		_ = db.AutoMigrate(&entity.User{})
		// SETUP CONNECTION POOL
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed get instance connection %v", err)
		}
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)

		InstanceDB = db
	})

	return InstanceDB
}

func CustomCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-API-Key")
		if c.Request.Method == "OPTIONS" {
			fmt.Println("Handling OPTIONS method")
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

// @title 	Tag Example Monitoring Service
// @version	1.0
// @description A Tag service API in Go using Gin framework
// @host 	localhost:8080
// @BasePath /api/v1
func main() {
	db := getInstance()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	//SETUP CORS
	r.Use(CustomCORSMiddleware())
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Inisialisasi collector metrik, and implemtntation into midddleware
	metrics := metric.NewAppMetricsExporter()
	r.Use(metrics.GinMiddleware())

	//INJECT HANDLER
	userHandler := handler.NewUserHandler(db, metrics)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "ASUK")
			return
		})

		v1.POST("/login", userHandler.Login)
		v1.POST("/register", userHandler.Register)
		v1.GET("/users", userHandler.GetUser)
	}

	/// BUAT EXSKPORTER BUAT SEND KE PROMETHEUS
	promServer := &http.Server{
		Addr:    ":8081",
		Handler: metrics.MetricsHandler(),
	}

	go func() {
		slog.Info("Listening And Serve Prometheus Exporter", slog.String("port", promServer.Addr))
		if err := promServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed: %s", err)
			panic(err)
		}
	}()

	/// RUN SERVER WEB SERVER
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r.Handler(),
	}

	go func() {
		slog.Info("Listening And Server HTTP on ", slog.String("port", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed: %s", err)
			panic(err)
		}
	}()

	/////GRATEFULLY SHOWDOWN/////
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
