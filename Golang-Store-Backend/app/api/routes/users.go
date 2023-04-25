package routes

import (
	"fmt"
	"net/http"

	"github.com/Apouzi/golang-shop/app/api/helpers"
)

type TestUser struct{
	PasswordHash string `json:"PasswordHash"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	Email string `json:"email"`

}

func (route *Routes) Register(w http.ResponseWriter, r *http.Request){
	// db := route.DB

	user := TestUser{}
	helpers.ReadJSON(w, r, &user)
	fmt.Println("Register",user)
	id := route.UserQuery.RegisterUserIntoDB(user.PasswordHash,user.FirstName,user.LastName,user.Email)



	fmt.Println("returned id", id)

	
}