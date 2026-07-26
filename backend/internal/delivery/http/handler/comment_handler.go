package handler

import (
	"myapp/internal/delivery/http/dto"
	"myapp/internal/delivery/http/middleware"
	"myapp/internal/service"
	"myapp/internal/utils"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	commentService service.CommentServiceInterface
}

func NewCommentHandler(commentService service.CommentServiceInterface) CommentHandler {
	return CommentHandler{
		commentService: commentService,
	}
}

func (h CommentHandler) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var commentDTO dto.CommentDTO
	if err := utils.Deserialize(r.Body, &commentDTO); err != nil {
		http.Error(w, "error during deserializing comment dto", http.StatusBadRequest)
		return
	}

	if validateErr := commentDTO.Validate(); validateErr != nil {
		http.Error(w, validateErr.Error(), http.StatusBadRequest)
		return
	}

	creatorID, parseErr := utils.ParseUserID(ctx.Value(middleware.UserContextKey{}))
	if parseErr != nil {
		http.Error(w, "error during parsing creator id", http.StatusBadRequest)
		return
	}

	postID, parseErr := strconv.ParseUint(r.PathValue("post_id"), 10, 64)
	if parseErr != nil {
		http.Error(w, "error during parsing post id", http.StatusBadRequest)
		return
	}

	commentID, apiErr := h.commentService.CreateComment(commentDTO, creatorID, uint(postID), ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]uint{"comment_id": commentID})
}

func (h CommentHandler) HandleGetAllCommentsOnPost(w http.ResponseWriter, r *http.Request) {
	postID, parseErr := strconv.ParseUint(r.PathValue("post_id"), 10, 64)
	if parseErr != nil {
		http.Error(w, "error during parsing post id", http.StatusBadRequest)
		return
	}

	dtos, apiErr := h.commentService.GetAllCommentsByPostID(uint(postID), r.Context())
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, dtos)
}
