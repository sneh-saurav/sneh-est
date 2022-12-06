package est

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/globalsign/est/controller"
	"github.com/globalsign/est/models"
)

// Create a new admin user
func HandleCreateAdminUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	err := controller.CreateAdminUser(&user)
	if err != nil {
		fmt.Println("ERROR WHILE CREATING ADMIN USER")
	}

	fmt.Fprint(w, "CREATE USER FLOW", user)
}
