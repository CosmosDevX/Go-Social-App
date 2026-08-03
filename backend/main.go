package main

import (
	"context"
	"log"
	"myapp/internal/config"
	"myapp/internal/delivery/http/handler"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/infrastructure"
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
	config := config.Config{}
	config.Load()

	//initialize clients
	redisClient := infrastructure.NewRedisClient()
	sqlxClient := infrastructure.NewSQLxClient(config.DBConnectionString)
	// base migration sqlxClient.CreateTables("CREATE TABLE posts(id SERIAL PRIMARY KEY,name VARCHAR(100),description VARCHAR(900),creator_id INTEGER,likes INTEGER DEFAULT 0,image_name VARCHAR(300)); CREATE TABLE post_likes(id SERIAL PRIMARY KEY,liked_user_id INTEGER,post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE); CREATE TABLE comments(id SERIAL PRIMARY KEY,text VARCHAR(250),post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,creator_id INTEGER); CREATE TABLE users(id SERIAL PRIMARY KEY,username VARCHAR(60) UNIQUE NOT NULL,password VARCHAR(100) NOT NULL);")

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
	jwtService := authorization.NewJWTService(config.SecretKey)
	authService := authorization.NewAuthService(userRepository, refreshTokenRepository, jwtService)
	userService := service.NewUserService(userRepository)
	postService := service.NewPostService(unitOfWork, postRepository, postLikeRepository, commentRepository, fileManager, userRepository)
	postLikeService := service.NewPostLikeService(unitOfWork, postRepository, postLikeRepository)
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

	if err := sqlxClient.Shutdown(); err != nil {
		log.Println(err)
	}
}
