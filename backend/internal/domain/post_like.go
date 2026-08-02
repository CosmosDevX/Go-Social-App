package domain

type PostLike struct {
	ID          int `db:"id"`
	LikedUserID int `db:"liked_user_id"`
	PostID      int `db:"post_id"`
}

/*
CREATE TABLE post_likes (
	id SERIAL PRIMARY KEY,
	liked_user_id INTEGER,
	post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE
);
*/
