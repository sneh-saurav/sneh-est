package models

import "time"

type User struct {
	UserID    string    `db:"user_id" json:"user_id"`
	IsActive  bool      `db:"active" json:"active"`
	IsAdmin   bool      `db:"isAdmin" json:"isAdmin"`
	Email     string    `db:"email_id" json:"email_id"`
	Password  string    `db:"password" json:"password,omitempty"`
	UserToken string    `db:"user_token" json:"user_token"`
	UserName  string    `db:"username" json:"username"`
	UserRole  int32     `db:"user_role" json:"user_role"`
	CreatedOn time.Time `db:"added_date" json:"added_date"`
	UpdatedOn time.Time `db:"updated_date" json:"updated_date"`
}

type UserResponse struct {
	AuthToken string    `json:"authenticationToken"`
	ExpAt     time.Time `json:"expiresAt"`
	UserName  string    `json:"user_name"`
	UserRole  int32     `json:"role"`
}
