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

func (h CommentHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var commentDTO dto.CommentDTO
	if err := utils.Deserialize(r.Body, &commentDTO); err != nil {
		WriteError(w, *domain.NewDeserializingError("error during deserializing comment"))
		return
	}

	if validateErr := commentDTO.Validate(); validateErr != nil {
		WriteError(w, *domain.NewValidationError(validateErr.Error()))
		return
	}

	creatorID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse creator id"))
		return
	}

	postID, parseErr := strconv.Atoi(r.PathValue("post_id"))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	commentID, domainErr := h.commentService.CreateComment(ctx, commentDTO, creatorID, postID)
	if domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, map[string]int{"comment_id": commentID})
}

func (h CommentHandler) HandleGetAllCommentsOnPost(w http.ResponseWriter, r *http.Request) {
	postID, parseErr := strconv.Atoi(r.PathValue("post_id"))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	dtos, domainErr := h.commentService.GetAllCommentsByPostID(r.Context(), postID)
	if domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, dtos)
}

func (h CommentHandler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserIDContextKey{}))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	commentID, parseErr := strconv.Atoi(r.PathValue("comment_id"))
	if parseErr != nil {
		WriteError(w, *domain.NewParseError("error during parse comment id"))
		return
	}

	if domainErr := h.commentService.DeleteComment(ctx, commentID, userID); domainErr != nil {
		WriteError(w, *domainErr)
		return
	}

	WriteJSON(w, map[string]string{"message": "comment deleted"})
}
