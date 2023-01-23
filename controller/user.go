package controller

import (
	"bytes"
	"fmt"
	_ "net/smtp"
	"text/template"

	"github.com/globalsign/est/models"
	"github.com/globalsign/est/security"
	smtpservice "github.com/globalsign/est/smtpservice"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func Create(u *models.User, ldb *sqlx.DB) error {

	if u.IsAdmin {
		u.UserRole = 1
	} else {
		u.UserRole = 2
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashed)

	id := uuid.New()
	u.UserToken = id.String()

	tx, _ := ldb.Beginx()

	uID, err := createUser(tx, u)

	if err == nil {
		err = SaveDeviceRegToken(tx, uID, id)
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()

		u.Password = ""
	}
	err = SendEmail(u)
	if err != nil {
		return err
	}

	return nil
}

func createUser(ldb *sqlx.Tx, u *models.User) (int64, error) {
	str := `INSERT INTO user (active, email_id, password, user_token, username, user_role)
	VALUES(:active, :email_id, :password, :user_token, :username, :user_role);`
	res, err := ldb.NamedExec(str, u)
	if err != nil {
		return 0, err
	}
	uID, _ := res.LastInsertId()

	return uID, nil
}

func Login(user_name string, pass string, ldb *sqlx.DB) (*models.UserResponse, string, error) {

	var res models.UserResponse
	usr, err := GetUserByUsername(user_name, ldb)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(pass))
	if err != nil {
		return nil, "", fmt.Errorf("your email and password didn’t match our records. please try again")
	}
	tok, err := security.CreateToken(user_name)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return nil, "", err
	}

	res.AuthToken = tok
	res.UserName = usr.UserName
	res.UserRole = usr.UserRole

	return &res, "", nil
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
func GetUserByUsername(user_name string, ldb *sqlx.DB) (*models.User, error) {

	var usr models.User
	str := `SELECT * FROM user WHERE username = '%s'`
	str = fmt.Sprintf(str, user_name)

	fmt.Println(str)
	err := ldb.Get(&usr, str)
	if err != nil {
		return nil, err
	}

	return &usr, err
}

func SaveDeviceRegToken(tx *sqlx.Tx, userId int64, token uuid.UUID) error {
	str := `INSERT INTO device_registration_token (user_id, registration_token)
	Values(?,?)`
	_, err := tx.Exec(str, userId, token)
	if err != nil {
		return err
	}
	return nil
}

func SendEmail(u *models.User) error {
	mailer := smtpservice.SMTP{}
	fname := fmt.Sprintf("%v/%v", "../.././templates", "sendtoken.html")
	tmp, err := template.ParseFiles(fname)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	err = tmp.Execute(&b, u)
	if err != nil {
		return err
	}

	err = mailer.Send("Here's the token to enroll your device", b.String(), nil, u.Email)
	if err != nil {
		return err
	}
	return nil
}
