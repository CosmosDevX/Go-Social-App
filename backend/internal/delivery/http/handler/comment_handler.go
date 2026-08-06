package handler

import (
	"context"
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/domain"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type CommentService interface {
	DeleteComment(ctx context.Context, commentID, userID int) *domain.DomainError
	CreateComment(ctx context.Context, commentDTO dto.CommentDTO, creatorID, postID int) (int, *domain.DomainError)
	GetAllCommentsByPostID(ctx context.Context, postID int) ([]dto.CommentDTO, *domain.DomainError)
}

type CommentHandler struct {
	commentService CommentService
}

func NewCommentHandler(commentService CommentService) CommentHandler {
	return CommentHandler{
		commentService: commentService,
	}
}

// HandleCreateComment godoc
// @Summary      Создать комментарий
// @Description  Добавляет комментарий к посту (текст 1-250 символов)
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        post_id  path  int              true  "ID поста"
// @Param        comment  body  dto.CommentDTO   true  "Текст комментария"
// @Success      200  {object}  CommentIDResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /comment/create/{post_id} [post]
func (h CommentHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var commentDTO dto.CommentDTO
	if err := utils.Deserialize(r.Body, &commentDTO); err != nil {
		utils.WriteError(w, *domain.NewDeserializingError("error during deserializing comment"))
		return
	}

	if validateErr := commentDTO.Validate(); validateErr != nil {
		utils.WriteError(w, *domain.NewValidationError(validateErr.Error()))
		return
	}

	creatorID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse creator id"))
		return
	}

	postID, parseErr := strconv.Atoi(r.PathValue("post_id"))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	commentID, domainErr := h.commentService.CreateComment(ctx, commentDTO, creatorID, postID)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"comment_id": commentID})
}

// HandleGetAllCommentsOnPost godoc
// @Summary      Все комментарии поста
// @Description  Возвращает список комментариев к посту
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        post_id  path  int  true  "ID поста"
// @Success      200  {array}   dto.CommentDTO
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Router       /comment/all/{post_id} [get]
func (h CommentHandler) HandleGetAllCommentsOnPost(w http.ResponseWriter, r *http.Request) {
	postID, parseErr := strconv.Atoi(r.PathValue("post_id"))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	dtos, domainErr := h.commentService.GetAllCommentsByPostID(r.Context(), postID)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, dtos)
}

// HandleDeleteComment godoc
// @Summary      Удалить комментарий
// @Description  Удаляет комментарий (только автор)
// @Tags         comments
// @Produce      json
// @Security     BearerAuth
// @Param        comment_id  path  int  true  "ID комментария"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      403  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /comment/{comment_id} [delete]
func (h CommentHandler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	commentID, parseErr := strconv.Atoi(r.PathValue("comment_id"))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse comment id"))
		return
	}

	if domainErr := h.commentService.DeleteComment(ctx, commentID, userID); domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "comment deleted"})
}
