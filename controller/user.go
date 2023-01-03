package controller

import (
	"fmt"

	"github.com/globalsign/est/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdminUser(u *models.User, ldb *sqlx.DB) error {

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	str := `INSERT INTO user (active, added_date, email_id, password, updated_date, user_token, username)
	 VALUES(:active, :added_date, :email_id, :password, :updated_date, :user_token, :username);`

	u.IsActive = true

	u.Password = string(hashed)

	id := uuid.New()
	u.UserToken = id.String()

	_, err = ldb.NamedExec(str, u)
	if err != nil {
		return err
	}
	u.Password = ""

	return nil
}

func CheckIfUserExist(u string, ldb *sqlx.DB) bool {
	userCount := 0
	str := `SELECT count(*) FROM user WHERE username LIKE '%%%s%%'`

	str = fmt.Sprintf(str, u)

	err := ldb.Get(&userCount, str)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	if userCount > 0 {
		return true
	} else {
		return false
	}
}
