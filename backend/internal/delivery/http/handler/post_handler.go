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
	CreatePost(postDTO dto.PostDTO, creatorID uint, file multipart.File, header *multipart.FileHeader, ctx context.Context) (uint, *domain.DomainError)
	GetCurrentUserPosts(userID uint, ctx context.Context) ([]dto.PostDTO, *domain.DomainError)
	GetUserPostsByUsername(username string, currentUserID uint, ctx context.Context) ([]dto.PostDTO, *domain.DomainError)
	DeletePost(postID, userID uint, ctx context.Context) *domain.DomainError
	GetPostFeed(currentUserID uint, ctx context.Context) ([]dto.PostDTO, *domain.DomainError)
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
		utils.WriteError(w, *domain.NewFileError("file too large"))
		return
	}

	file, header, err := r.FormFile("post_image")
	contentType := header.Header.Get("Content-Type")
	if err != nil {
		utils.WriteError(w, *domain.NewFileError("file not found"))
		return
	}
	if !strings.HasPrefix(contentType, "image/") {
		utils.WriteError(w, *domain.NewFileError("invalid file type"))
		return
	}

	defer file.Close()

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

	creatorID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	postID, domainErr := h.postService.CreatePost(postDTO, creatorID, file, header, ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]uint{"post_id": postID})
}

func (h PostHandler) HandleGetCurrentUserPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse creator id"))
		return
	}

	dtos, domainErr := h.postService.GetCurrentUserPosts(parsedUserID, ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, dtos)
}

func (h PostHandler) HandleGetUserPostsByUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	dtos, domainErr := h.postService.GetUserPostsByUsername(r.PathValue("username"), parsedUserID, ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, dtos)
}

func (h PostHandler) HandleGetPostFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	dtos, domainErr := h.postService.GetPostFeed(parsedUserID, ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, dtos)
}

func (h PostHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	postID, parseErr := strconv.ParseUint(r.PathValue("post_id"), 10, 64)
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	if domainErr := h.postService.DeletePost(uint(postID), userID, ctx); domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "post deleted"})
}
