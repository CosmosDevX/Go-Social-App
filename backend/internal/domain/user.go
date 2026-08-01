// Package domain
package domain

type User struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}

/*
CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	username VARCHAR(60) UNIQUE NOT NULL,
	password VARCHAR(100) NOT NULL
);
*/
