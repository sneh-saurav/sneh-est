package est

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/globalsign/est/controller"
	"github.com/globalsign/est/db"
	"github.com/globalsign/est/models"
	"github.com/globalsign/est/security"
)

// Create a new admin user
func HandleCreateAdmin(w http.ResponseWriter, r *http.Request) {
	//Connecting to the db
	ldb, errs := db.Connect("")
	if errs != nil {
		log.Fatal(errs)
	}
	defer ldb.Close()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintln(w, "ERROR WHILE UNMARSHALLING")
	}

	IsUserExist := controller.CheckIfUserExist(user.UserName, ldb)

	if IsUserExist {
		fmt.Fprintln(w, "USERNAME ALREADY EXISTS")
	} else {
		err := controller.Create(&user, ldb)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println("ERROR WHILE CREATING ADMIN USER", err)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "USER CREATED", user)
		}
	}
}

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	tok := security.ExtractToken(r)
	err := security.VerifyToken(tok)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Connecting to the db
	ldb, errs := db.Connect("")
	if errs != nil {
		log.Fatal(errs)
	}
	defer ldb.Close()

	var user models.User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Fprintln(w, "ERROR WHILE UNMARSHALLING")
	}

	IsUserExist := controller.CheckIfUserExist(user.UserName, ldb)

	if IsUserExist {
		fmt.Fprintln(w, "USERNAME ALREADY EXISTS")
	} else {
		err := controller.Create(&user, ldb)
		if err != nil {
			fmt.Println("ERROR WHILE CREATING ADMIN USER")
		}
		fmt.Fprint(w, "USER CREATED", user)
	}
}

// Login a User
func HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	//Connecting to the db
	ldb, errs := db.Connect("")
	if errs != nil {
		log.Fatal(errs)
	}
	defer ldb.Close()

	p := struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		fmt.Fprintln(w, "ERROR WHILE UNMARSHALLING")
	}

	u, _, err := controller.Login(p.UserName, p.Password, ldb)
	res, _ := json.Marshal(u)
	if err != nil {
		fmt.Fprint(w, "ERROR WHILE LOGIN - ", err)
	} else {
		fmt.Fprint(w, "LOGIN SUCCESSFUL \n")
		w.Write(res)
	}
}
