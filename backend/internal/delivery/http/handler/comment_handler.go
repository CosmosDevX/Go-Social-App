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
	DeleteComment(commentID, userID int, ctx context.Context) *domain.DomainError
	CreateComment(commentDTO dto.CommentDTO, creatorID, postID int, ctx context.Context) (int, *domain.DomainError)
	GetAllCommentsByPostID(postID int, ctx context.Context) ([]dto.CommentDTO, *domain.DomainError)
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
		utils.WriteError(w, *domain.NewDeserializingError("error during deserializing comment"))
		return
	}

	if validateErr := commentDTO.Validate(); validateErr != nil {
		utils.WriteError(w, *domain.NewValidationError(validateErr.Error()))
		return
	}

	creatorID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse creator id"))
		return
	}

	postID, parseErr := strconv.ParseUint(r.PathValue("post_id"), 10, 64)
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	commentID, domainErr := h.commentService.CreateComment(commentDTO, int(creatorID), int(postID), ctx)
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]int{"comment_id": commentID})
}

func (h CommentHandler) HandleGetAllCommentsOnPost(w http.ResponseWriter, r *http.Request) {
	postID, parseErr := strconv.ParseUint(r.PathValue("post_id"), 10, 64)
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse post id"))
		return
	}

	dtos, domainErr := h.commentService.GetAllCommentsByPostID(int(postID), r.Context())
	if domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, dtos)
}

func (h CommentHandler) HandleDeleteComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse user id"))
		return
	}

	commentID, parseErr := strconv.ParseUint(r.PathValue("comment_id"), 10, 64)
	if parseErr != nil {
		utils.WriteError(w, *domain.NewParseError("error during parse comment id"))
		return
	}

	if domainErr := h.commentService.DeleteComment(int(commentID), int(userID), ctx); domainErr != nil {
		utils.WriteError(w, *domainErr)
		return
	}

	utils.WriteJSON(w, map[string]string{"message": "comment deleted"})
}
