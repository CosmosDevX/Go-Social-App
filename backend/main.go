package main

import (
	"context"
	"log/slog"
	"myapp/internal/config"
	"myapp/internal/delivery/http/handler"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/infrastructure"
	"myapp/internal/logger"
	"myapp/internal/repository"
	"myapp/internal/service"
	"myapp/internal/service/authorization"
	"myapp/internal/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis_rate/v10"
)

func main() {
	cfg := config.Config{}
	cfg.Load()

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "text"
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Setup(logFormat, logLevel)

	slog.Info("starting application")

	//initialize clients
	redisClient := infrastructure.NewRedisClient()
	sqlxClient := infrastructure.NewSQLxClient(cfg.DBConnectionString)

	rateLimiter := redis_rate.NewLimiter(redisClient.GetClient())

	//initialize managers
	fileManager := utils.NewFileManager()
	unitOfWork := repository.NewUnitOfWork(sqlxClient.GetDB())

	//initialize repositories
	userRepository := repository.NewUserRepository(sqlxClient.GetDB())
	postRepository := repository.NewPostRepository(sqlxClient.GetDB())
	postLikeRepository := repository.NewPostLikeRepository(sqlxClient.GetDB())
	commentRepository := repository.NewCommentRepository(sqlxClient.GetDB())
	refreshTokenRepository := repository.NewRefreshTokenRepository(redisClient.GetClient())

	//initialize services
	jwtService := authorization.NewJWTService(cfg.SecretKey)
	authService := authorization.NewAuthService(userRepository, refreshTokenRepository, jwtService)
	userService := service.NewUserService(userRepository)
	postService := service.NewPostService(unitOfWork, postRepository, postLikeRepository, commentRepository, fileManager, userRepository)
	postLikeService := service.NewPostLikeService(unitOfWork)
	commentService := service.NewCommentService(commentRepository, userRepository)

	//initialize handlers
	authHandler := handler.NewAuthHandler(authService, *rateLimiter)
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	postLikeHandler := handler.NewPostLikeHandler(postLikeService)
	commentHandler := handler.NewCommentHandler(commentService)

	//initialize middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth", authHandler.HandleAuth)
	mux.HandleFunc("POST /refresh", authHandler.HandleRefresh)
	mux.HandleFunc("POST /logout", authMiddleware.Protect(authHandler.HandleLogout))

	mux.HandleFunc("POST /user/create", userHandler.HandleCreateUser)
	mux.HandleFunc("GET /user/get_username_by_id/{user_id}", userHandler.HandleGetUsernameByID)
	mux.HandleFunc("GET /user/current/profile", authMiddleware.Protect(userHandler.HandleCurrentUserProfile))

	mux.HandleFunc("POST /post/create", authMiddleware.Protect(postHandler.HandleCreatePost))
	mux.HandleFunc("DELETE /post/{post_id}", authMiddleware.Protect(postHandler.HandleDeletePost))
	mux.HandleFunc("GET /post/current_user/all", authMiddleware.Protect(postHandler.HandleGetCurrentUserPosts))
	mux.HandleFunc("GET /post/{username}/all", authMiddleware.Protect(postHandler.HandleGetUserPostsByUsername))
	mux.HandleFunc("GET /post/feed", authMiddleware.Protect(postHandler.HandleGetPostFeed))

	mux.HandleFunc("POST /post/like/{post_id}", authMiddleware.Protect(postLikeHandler.HandleLikePost))
	mux.HandleFunc("POST /post/dislike/{post_id}", authMiddleware.Protect(postLikeHandler.HandleDislikePost))

	mux.HandleFunc("POST /comment/create/{post_id}", authMiddleware.Protect(commentHandler.HandleCreateComment))
	mux.HandleFunc("GET /comment/all/{post_id}", authMiddleware.Protect(commentHandler.HandleGetAllCommentsOnPost))
	mux.HandleFunc("DELETE /comment/{comment_id}", authMiddleware.Protect(commentHandler.HandleDeleteComment))

	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	handlerChain := http.TimeoutHandler(
		middleware.LoggingMiddleware(
			middleware.CorsMiddleware(mux),
		),
		time.Minute,
		"request timeout",
	)

	httpService := service.NewHTTPService(handlerChain)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("http server listening", "addr", ":8080")
		if err := httpService.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			quit <- syscall.SIGTERM
		}
	}()

	<-quit
	slog.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpService.Server.Shutdown(ctx); err != nil {
		slog.Error("error during server shutdown", "error", err)
	}

	if err := redisClient.Shutdown(); err != nil {
		slog.Error("error during redis shutdown", "error", err)
	}

	if err := sqlxClient.Shutdown(); err != nil {
		slog.Error("error during database shutdown", "error", err)
	}

	slog.Info("application stopped")
}
