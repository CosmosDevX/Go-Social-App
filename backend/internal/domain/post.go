package domain

import (
	"myapp/internal/delivery/http/dto"
)

type Post struct {
	ID              int    `db:"id"`
	PostName        string `db:"name"`
	PostDescription string `db:"description"`
	CreatorID       int    `db:"creator_id"`
	Likes           int    `db:"likes"`
	ImageName       string `db:"image_name"`
	CreatorUsername string `db:"creator_username"`
}

func (p Post) ToPostDTO() dto.PostDTO {
	return dto.PostDTO{
		PostID:          p.ID,
		PostName:        p.PostName,
		PostDescription: p.PostDescription,
		CreatorID:       p.CreatorID,
		Likes:           p.Likes,
		CreatorName:     p.CreatorUsername,
		ImageName:       p.ImageName,
	}
}

/*
CREATE TABLE posts(
	id SERIAL PRIMARY KEY,
	name VARCHAR(100),
	description VARCHAR(900),
	creator_id INTEGER,
	likes INTEGER DEFAULT 0,
	image_name VARCHAR(300)
);
*/
