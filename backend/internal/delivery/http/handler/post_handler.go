package handler

import (
	"myapp/internal/delivery/http/dto"
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

	postID, apiErr := h.postService.CreatePost(postDTO, ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]uint{"post_id": postID})
}

func (h PostHandler) HandleGetUserPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dtos, apiErr := h.postService.GetAllUserPosts(ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, dtos)
}

func (h PostHandler) HandleLikePost(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID, parseErr := strconv.ParseUint(r.PathValue("post_id"), 10, 64)
	if parseErr != nil {
		http.Error(w, "error during parsing post id", http.StatusBadRequest)
		return
	}

	likes, apiErr := h.postService.LikePost(uint(postID), ctx)
	if apiErr != nil {
		http.Error(w, apiErr.Message, utils.IdentifyAPIError(apiErr.Code))
		return
	}

	utils.WriteJSON(w, map[string]int{"likes": likes})
}
