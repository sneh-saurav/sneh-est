package models

import "time"

type User struct {
	UserID     string    `db:"user_id" json:"user_id"`
	UserName   string    `db:"username" json:"username"`
	IsActive   bool      `db:"active" json:"active"`
	Email      string    `db:"email_id" json:"email_id"`
	Password   string    `db:"-" json:"password,omitempty"`
	HashedPass []byte    `db:"password" json:"-"`
	UserToken  string    `db:"user_token" json:"user_token"`
	CreatedOn  time.Time `db:"added_date" json:"added_date"`
	UpdatedOn  time.Time `db:"updated_date" json:"updated_date"`
}
