package handler

import (
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/service"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type PostLikeHandler struct {
	postLikeService service.PostLikeServiceInterface
}

func NewPostLikeHandler(postLikeService service.PostLikeServiceInterface) PostLikeHandler {
	return PostLikeHandler{
		postLikeService: postLikeService,
	}
}

func (h PostLikeHandler) HandleLikePost(w http.ResponseWriter, r *http.Request) {
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

	likes, apiErr := h.postLikeService.LikePost(uint(postID), userID, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}

func (h PostLikeHandler) HandleDislikePost(w http.ResponseWriter, r *http.Request) {
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

	likes, apiErr := h.postLikeService.DislikePost(uint(postID), userID, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}
