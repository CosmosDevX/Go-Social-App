package domain

type PostLike struct {
	ID          int `db:"id"`
	LikedUserID int `db:"liked_user_id"`
	PostID      int `db:"post_id"`
}
