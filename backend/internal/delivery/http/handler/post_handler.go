package handler

import (
	"context"
	"mime/multipart"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

type PostService interface {
	CreatePost(ctx context.Context, postDTO dto.PostDTO, creatorID int, file multipart.File, header *multipart.FileHeader) (int, *domain.DomainError)
	GetCurrentUserPosts(ctx context.Context, userID int) ([]dto.PostDTO, *domain.DomainError)
	GetUserPostsByUsername(ctx context.Context, username string, currentUserID int) ([]dto.PostDTO, *domain.DomainError)
	DeletePost(ctx context.Context, postID, userID int) *domain.DomainError
	GetPostFeed(ctx context.Context, currentUserID int) ([]dto.PostDTO, *domain.DomainError)
}

type PostHandler struct {
	postService PostService
}

func NewPostHandler(postService PostService) PostHandler {
	return PostHandler{
		postService: postService,
	}
}

// HandleCreatePost godoc
// @Summary      Создать пост
// @Description  Создаёт пост с опциональным изображением (multipart/form-data). Макс. размер файла 10 МБ.
// @Tags         posts
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        post_name         formData  string  true   "Название поста (5-100 символов)"
// @Param        post_description  formData  string  true   "Описание (1-900 символов)"
// @Param        post_image        formData  file    false  "Изображение (только image/*)"
// @Success      200  {object}  PostIDResponse
// @Failure      400  {object}  ErrorResponse  "VALIDATION_ERROR / FILE_ERROR"
// @Failure      401  {object}  ErrorResponse
// @Router       /post/create [post]
func (h PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	if err := r.ParseMultipartForm(20 * 1024 * 1024); err != nil {
		utils.WriteError(w, *domain.NewFileError("file too large"))
		return
	}

	file, header, _ := r.FormFile("post_image")
	var contentType string
	if header != nil {
		contentType = header.Header.Get("Content-Type")
	}
	if file != nil && !strings.HasPrefix(contentType, "image/") {
		utils.WriteError(w, *domain.NewFileError("invalid file type"))
		return
	}
	if file != nil {
		defer file.Close()
	}

	postName := r.FormValue("post_name")
	postDescription := r.FormValue("post_description")

	postDTO := dto.PostDTO{
		PostName:        postName,
		PostDescription: postDescription,
	}

	if validateErr := postDTO.Validate(); validateErr != nil {
		utils.WriteError(w, *domain.NewValidationError(validateErr.Error()))
		return
	}

	creatorID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	postID, domainErr := h.postService.CreatePost(ctx, postDTO, creatorID, file, header)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"post_id": postID})
}

// HandleGetCurrentUserPosts godoc
// @Summary      Посты текущего пользователя
// @Description  Возвращает все посты авторизованного пользователя
// @Tags         posts
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.PostDTO
// @Failure      401  {object}  ErrorResponse
// @Router       /post/current_user/all [get]
func (h PostHandler) HandleGetCurrentUserPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse creator id"))
		return
	}

	dtos, domainErr := h.postService.GetCurrentUserPosts(ctx, parsedUserID)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, dtos)
}

// HandleGetUserPostsByUsername godoc
// @Summary      Посты пользователя по username
// @Description  Возвращает посты указанного пользователя (с флагом is_liked для текущего юзера)
// @Tags         posts
// @Produce      json
// @Security     BearerAuth
// @Param        username  path  string  true  "Username автора"
// @Success      200  {array}   dto.PostDTO
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /post/{username}/all [get]
func (h PostHandler) HandleGetUserPostsByUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	dtos, domainErr := h.postService.GetUserPostsByUsername(ctx, r.PathValue("username"), parsedUserID)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, dtos)
}

// HandleGetPostFeed godoc
// @Summary      Лента постов
// @Description  Возвращает ленту постов для текущего пользователя
// @Tags         posts
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   dto.PostDTO
// @Failure      401  {object}  ErrorResponse
// @Router       /post/feed [get]
func (h PostHandler) HandleGetPostFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	dtos, domainErr := h.postService.GetPostFeed(ctx, parsedUserID)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, dtos)
}

// HandleDeletePost godoc
// @Summary      Удалить пост
// @Description  Удаляет пост (только автор)
// @Tags         posts
// @Produce      json
// @Security     BearerAuth
// @Param        post_id  path  int  true  "ID поста"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /post/{post_id} [delete]
func (h PostHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	postID, parseErr := strconv.Atoi(r.PathValue("post_id"))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	if domainErr := h.postService.DeletePost(ctx, postID, userID); domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "post deleted"})
}
