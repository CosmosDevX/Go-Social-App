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
	userRepository := repository.UserRepository{}
	postRepository := repository.PostRepository{}
	postLikeRepository := repository.PostLikeRepository{}
	commentRepository := repository.CommentRepository{}
	refreshTokenRepository := repository.NewRefreshTokenRepository(redisClient.GetClient())

	//initialize services
	jwtService := authorization.NewJWTService()
	authService := authorization.NewAuthService(userRepository, refreshTokenRepository, jwtService, gormClient.GetDB())
	userService := service.NewUserService(userRepository, gormClient.GetDB())
	postService := service.NewPostService(postRepository, postLikeRepository, commentRepository, gormClient.GetDB())
	postLikeService := service.NewPostLikeService(postRepository, postLikeRepository, gormClient.GetDB())
	commentService := service.NewCommentService(commentRepository, gormClient.GetDB())

	//initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	postLikeHandler := handler.NewPostLikeHandler(postLikeService)
	commentHandler := handler.NewCommentHandler(commentService)

	//initialize middlewares
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /auth", authHandler.AuthHandler)
	mux.HandleFunc("POST /refresh", authHandler.RefreshHandler)

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
