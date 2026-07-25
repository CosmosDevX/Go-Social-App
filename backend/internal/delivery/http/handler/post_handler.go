package handler

import (
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/service"
	"myapp/internal/utils"
	"net/http"
	"strconv"
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
	var postDTO dto.PostDTO
	if err := utils.Deserialize(r.Body, &postDTO); err != nil {
		http.Error(w, "error during deserializing post dto", http.StatusBadRequest)
		return
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

	postID, apiErr := h.postService.CreatePost(postDTO, creatorID, ctx)
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

func (h PostHandler) HandleLikePost(w http.ResponseWriter, r *http.Request) {
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

	likes, apiErr := h.postService.LikePost(uint(postID), userID, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}

func (h PostHandler) HandleDislikePost(w http.ResponseWriter, r *http.Request) {
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

	likes, apiErr := h.postService.DislikePost(uint(postID), userID, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}
