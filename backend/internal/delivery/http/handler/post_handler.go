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

func (h PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	if err := r.ParseMultipartForm(20 * 1024 * 1024); err != nil {
		WriteError(w, *domain.NewFileError("file too large"))
		return
	}

	file, header, _ := r.FormFile("post_image")
	var contentType string
	if header != nil {
		contentType = header.Header.Get("Content-Type")
	}
	if file != nil && !strings.HasPrefix(contentType, "image/") {
		WriteError(w, *domain.NewFileError("invalid file type"))
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
		WriteError(w, *domain.NewValidationError(validateErr.Error()))
		return
	}

	creatorID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	postID, domainErr := h.postService.CreatePost(ctx, postDTO, creatorID, file, header)
	if domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, map[string]int{"post_id": postID})
}

func (h PostHandler) HandleGetCurrentUserPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse creator id"))
		return
	}

	dtos, domainErr := h.postService.GetCurrentUserPosts(ctx, parsedUserID)
	if domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, dtos)
}

func (h PostHandler) HandleGetUserPostsByUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	dtos, domainErr := h.postService.GetUserPostsByUsername(ctx, r.PathValue("username"), parsedUserID)
	if domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, dtos)
}

func (h PostHandler) HandleGetPostFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	dtos, domainErr := h.postService.GetPostFeed(ctx, parsedUserID)
	if domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, dtos)
}

func (h PostHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	postID, parseErr := strconv.Atoi(r.PathValue("post_id"))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	if domainErr := h.postService.DeletePost(ctx, postID, userID); domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, map[string]string{"message": "post deleted"})
}
