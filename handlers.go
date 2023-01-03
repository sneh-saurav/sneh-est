package est

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/globalsign/est/controller"
	"github.com/globalsign/est/db"
	"github.com/globalsign/est/models"
)

// Create a new admin user
func HandleCreateAdminUser(w http.ResponseWriter, r *http.Request) {
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
		err := controller.CreateAdminUser(&user, ldb)
		if err != nil {
			fmt.Println("ERROR WHILE CREATING ADMIN USER")
		}
		fmt.Fprint(w, "USER CREATED", user)
	}
}
