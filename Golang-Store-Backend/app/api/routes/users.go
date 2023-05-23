package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Apouzi/golang-shop/app/api/helpers"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TestUser struct{
	PasswordHash string `json:"PasswordHash"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	Email string `json:"email"`

}

func (route *Routes) Register(w http.ResponseWriter, r *http.Request){
	db := route.DB

	user := TestUser{}
	helpers.ReadJSON(w, r, &user)
	fmt.Println("Register",user)
	// passByte, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash),bcrypt.DefaultCost)
	// if err != nil{
	// 	fmt.Println("Password Gen issue", err)
	// }
	id, err := route.UserQuery.RegisterUserIntoDB(db,user.PasswordHash,user.FirstName,user.LastName,user.Email)
	if err != nil{
		fmt.Println(err)
		return
	}

	fmt.Println("returned id", id)

	
}

type LoginUser struct{
	Email string `email:"Email"`
	Password string `json:"PasswordHash"`
}

type SendBackLogin struct{
	Email string `email:"Email"`
}
func (route *Routes) Login(w http.ResponseWriter, r *http.Request){
	db := route.DB
	login := LoginUser{}
	helpers.ReadJSON(w, r, &login)
	_, passwordStored, userID,err := route.UserQuery.LoginUserDB(db, login.Email)
	fmt.Println("sent in password",login.Password)
	if err != nil{
		fmt.Println(err)
	}


	err = bcrypt.CompareHashAndPassword([]byte(passwordStored), []byte(login.Password))

	if err !=nil{
		fmt.Println("password does not match")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":time.Now().Add(time.Minute * 10).Unix(),
		"admin":"False",
		"email":login.Email,
		"userId":userID,
	})
	// Remove the testing key for this
	tokenString, err := token.SignedString([]byte("Testing key"))
	sendBack := SendBackLogin{Email: tokenString}

	helpers.WriteJSON(w, http.StatusAccepted, &sendBack)
}

type JWTtest struct{
	Token string `json:"JWT"`
}

func (route *Routes) VerifyTest(w http.ResponseWriter, r *http.Request){
	fmt.Println("THIS IS HIT")
	
	// fmt.Println("verify email",ctx)
	// jwttest := &JWTtest{}
	// helpers.ReadJSON(w, r, &jwttest)
	// token, err := jwt.Parse(jwttest.Token, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte("Testing key"), nil
	// })
	// if err != nil{
	// 	fmt.Println("verify test error")
	// 	fmt.Println(err)
	// }
	// if token.Valid{
	// 	fmt.Println("token validated")
	// }
	// fmt.Println(token.Claims)
}

type UserProfile struct{
	Cell int `json:"Cell"`
	Home int `json:"Home"`
}

func (route *Routes) UserProfile(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("userId")
	UserProfile := &UserProfile{}
	cell, home, err := route.UserQuery.GetUserProfile(route.DB, userID)


	if err != nil{
		fmt.Println("Error with getting userprofile in users.go")
	}
	
	UserProfile.Cell = cell
	UserProfile.Home = home

	helpers.WriteJSON(w,http.StatusAccepted, &UserProfile)


}