package handler

import (
	"fmt"
	"io"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/service"
	"myapp/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type PostHandler struct {
	postService service.PostServiceInterface
}

func NewPostHandler(postService service.PostServiceInterface) PostHandler {
	return PostHandler{
		postService: postService,
	}
}

func (h PostHandler) HandleCreatePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	if err := r.ParseMultipartForm(20 * 1024 * 1024); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("post_image")
	if err != nil {
		http.Error(w, "file not found", http.StatusBadRequest)
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
		http.Error(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	creatorID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		http.Error(w, "error during parsing creator id", http.StatusBadRequest)
		return
	}

	fileExt := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)
	savePath := filepath.Join("uploads", filename)

	dst, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "write error", http.StatusInternalServerError)
		return
	}

	postID, apiErr := h.postService.CreatePost(postDTO, creatorID, filename, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]uint{"post_id": postID})
}

func (h PostHandler) HandleGetCurrentUserPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		http.Error(w, "error during parsing creator id", http.StatusBadRequest)
		return
	}

	dtos, apiErr := h.postService.GetCurrentUserPosts(parsedUserID, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, dtos)
}

func (h PostHandler) HandleGetUserPostsByUsername(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		http.Error(w, "error during parsing user id", http.StatusBadRequest)
		return
	}

	dtos, apiErr := h.postService.GetUserPostsByUsername(r.PathValue("username"), parsedUserID, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, dtos)
}

func (h PostHandler) HandleGetPostFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	parsedUserID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		http.Error(w, "error during parsing user id", http.StatusBadRequest)
		return
	}

	dtos, apiErr := h.postService.GetPostFeed(parsedUserID, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, dtos)
}

func (h PostHandler) HandleDeletePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		http.Error(w, "error during parsing user id", http.StatusBadRequest)
		return
	}

	postID, parseErr := strconv.ParseUint(r.PathValue("post_id"), 10, 64)
	if parseErr != nil {
		http.Error(w, "error during parsing post id", http.StatusBadRequest)
		return
	}

	if apiErr := h.postService.DeletePost(uint(postID), userID, ctx); apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "post deleted"})
}
