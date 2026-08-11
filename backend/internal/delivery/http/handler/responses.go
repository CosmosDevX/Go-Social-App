package handler

import "myapp/internal/delivery/http/dto"

// ErrorResponse стандартная ошибка API
type ErrorResponse struct {
	Code    string `json:"code" example:"VALIDATION_ERROR"`
	Message string `json:"message" example:"username cannot be bigger than 60 and lower than 3"`
}

// AccessTokenResponse ответ с access token
type AccessTokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
}

// MessageResponse простой текстовый ответ
type MessageResponse struct {
	Message string `json:"message" example:"logout successful"`
}

// UserIDResponse ответ с id пользователя
type UserIDResponse struct {
	UserID int `json:"user_id" example:"1"`
}

// PostIDResponse ответ с id поста
type PostIDResponse struct {
	PostID int `json:"post_id" example:"42"`
}

// CommentIDResponse ответ с id комментария
type CommentIDResponse struct {
	CommentID int `json:"comment_id" example:"7"`
}

// LikesResponse ответ с количеством лайков
type LikesResponse struct {
	Likes int `json:"likes" example:"15"`
}

// UsernameResponse ответ с username
type UsernameResponse struct {
	Username string `json:"username" example:"john_doe"`
}

// PostsResponse ответ с постами пользователя(не лента)
type PostsResponse struct {
	Posts    []dto.PostDTO `json:"posts"`
	Page     int           `json:"page" example:"1"`
	PageSize int           `json:"page_size" example:"10"`
}
