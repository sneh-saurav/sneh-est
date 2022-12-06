package controller

import (
	"fmt"
	"log"

	"github.com/globalsign/est/db"
	"github.com/globalsign/est/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdminUser(u *models.User) error {
	//Connecting to the db
	ldb, errs := db.Connect("")
	if errs != nil {
		log.Fatal(errs)
	}
	defer ldb.Close()

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPass = hashed

	u.IsActive = true

	str := `INSERT INTO user (active, added_date, email_id, password, updated_date, user_token, username)
	 VALUES(:active, :added_date, :email_id, :password, :updated_date, :user_token, :username);`

	//unhashed := u.Password

	u.Password = string(hashed)
	id := uuid.New()
	u.UserToken = id.String()

	_, err = ldb.NamedExec(str, u)
	if err != nil {
		return err
	}
	u.Password = ""

	fmt.Print("CREATE USER FLOW----------", u)
	return nil
}
