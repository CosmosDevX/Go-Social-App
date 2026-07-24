package main

import (
	"context"
	"log"
	"myapp/internal/delivery/http/handler"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/infrastructure"
	"myapp/internal/repository"
	"myapp/internal/service"
	"myapp/internal/service/authorization"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//initialize clients
	gormClient := infrastructure.GormClient{}
	gormClient.Setup()
	redisClient := infrastructure.NewRedisClient()

	//initialize repositories
	userRepository := repository.NewUserRepository(gormClient.GetDB())
	refreshTokenRepository := repository.NewRefreshTokenRepository(redisClient.GetClient())

	//initialize services
	jwtService := authorization.NewJWTService()
	authService := authorization.NewAuthService(userRepository, refreshTokenRepository, jwtService)
	userService := service.NewUserService(userRepository)

	//initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	//initialize middlewares
	//authMiddleware := middleware.NewAuthMiddleware(jwtService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth", authHandler.AuthHandler)
	mux.HandleFunc("POST /refresh", authHandler.RefreshHandler)

	mux.HandleFunc("POST /user/create", userHandler.CreateUserHandler)

	httpService := service.NewHTTPService(http.TimeoutHandler(middleware.CorsMiddleware(mux), time.Minute, "request timeout"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpService.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpService.Server.Shutdown(ctx); err != nil {
		log.Println("error during shutdown server")
	}

	if err := redisClient.Shutdown(); err != nil {
		log.Println(err)
	}

	if err := gormClient.Shutdown(); err != nil {
		log.Println(err)
	}
}
