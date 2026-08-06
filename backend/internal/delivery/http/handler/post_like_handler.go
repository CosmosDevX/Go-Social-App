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
	LikePost(ctx context.Context, postID, userID int) (int, *domain.DomainError)
	DislikePost(ctx context.Context, postID, userID int) (int, *domain.DomainError)
}

type PostLikeHandler struct {
	postLikeService PostLikeService
}

func NewPostLikeHandler(postLikeService PostLikeService) PostLikeHandler {
	return PostLikeHandler{
		postLikeService: postLikeService,
	}
}

// HandleLikePost godoc
// @Summary      Поставить лайк
// @Description  Ставит лайк посту. Возвращает актуальное количество лайков.
// @Tags         likes
// @Produce      json
// @Security     BearerAuth
// @Param        post_id  path  int  true  "ID поста"
// @Success      200  {object}  LikesResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /post/like/{post_id} [post]
func (h PostLikeHandler) HandleLikePost(w http.ResponseWriter, r *http.Request) {
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

	likes, domainErr := h.postLikeService.LikePost(ctx, postID, userID)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}

// HandleDislikePost godoc
// @Summary      Убрать лайк
// @Description  Убирает лайк с поста. Возвращает актуальное количество лайков.
// @Tags         likes
// @Produce      json
// @Security     BearerAuth
// @Param        post_id  path  int  true  "ID поста"
// @Success      200  {object}  LikesResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /post/dislike/{post_id} [post]
func (h PostLikeHandler) HandleDislikePost(w http.ResponseWriter, r *http.Request) {
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

	likes, domainErr := h.postLikeService.DislikePost(ctx, postID, userID)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}
