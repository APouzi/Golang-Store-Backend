package routes

import (
	"fmt"
	"net/http"

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
	_, passwordStored, err := route.UserQuery.LoginUserDB(db, login.Email)
	fmt.Println("sent in password",login.Password)
	if err != nil{
		fmt.Println(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordStored), []byte(login.Password))

	if err !=nil{
		fmt.Println("password does not match")
		return
	}
	fmt.Println("success")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin":"True",
		"Hello!":"This is a hello",
	})

	tokenString, err := token.SignedString([]byte("Testing key"))
	sendBack := SendBackLogin{Email: tokenString}

	helpers.WriteJSON(w, http.StatusAccepted, &sendBack)
}

type JWTtest struct{
	Token string `json:"JWT"`
}

func (route *Routes) VerifyTest(w http.ResponseWriter, r *http.Request){
	jwttest := &JWTtest{}
	helpers.ReadJSON(w, r, &jwttest)
	token, err := jwt.Parse(jwttest.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("Testing key"), nil
	})
	if err != nil{
		fmt.Println("verify test error")
		fmt.Println(err)
	}
	if token.Valid{
		fmt.Println("token validated")
	}
	fmt.Println(token.Claims)
}