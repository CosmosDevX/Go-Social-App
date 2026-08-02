package handler

import (
	"context"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type PostLikeService interface {
	LikePost(postID, userID int, ctx context.Context) (int, *domain.DomainError)
	DislikePost(postID, userID int, ctx context.Context) (int, *domain.DomainError)
}

type PostLikeHandler struct {
	postLikeService PostLikeService
}

func NewPostLikeHandler(postLikeService PostLikeService) PostLikeHandler {
	return PostLikeHandler{
		postLikeService: postLikeService,
	}
}

func (h PostLikeHandler) HandleLikePost(w http.ResponseWriter, r *http.Request) {
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

	likes, domainErr := h.postLikeService.LikePost(int(postID), int(userID), ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}

func (h PostLikeHandler) HandleDislikePost(w http.ResponseWriter, r *http.Request) {
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

	likes, domainErr := h.postLikeService.DislikePost(int(postID), int(userID), ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}
