package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Apouzi/golang-shop/app/api/helpers"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct{
	Email string `json:"Email"`
	Password string `json:"Password"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	

}

type AdminReturn struct{
	ID int64 `json:"ID"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	Email string `json:"Email"`

}

func (route *Routes) AdminSuperUserCreation(w http.ResponseWriter, r *http.Request){
	query := "SELECT COUNT(UserID) FROM tblUser"
	sqlRes := route.DB.QueryRow(query)
	if sqlRes.Err()!= nil{
		fmt.Println("Error in AdminSuperUserCreation Count check", sqlRes.Err().Error())
	}
	var rowCount int
	sqlRes.Scan(&rowCount)
	if rowCount != 0{
		fmt.Println("Can't create super user, users already exist", rowCount)
		return
	}
	user := User{}
	helpers.ReadJSON(w, r, &user)
	id, err := route.UserQuery.RegisterAdminIntoDB(route.DB,user.Password,user.FirstName,user.LastName,user.Email)
	if err != nil{
		fmt.Println(err)
		helpers.ErrorJSON(w,err,http.StatusBadRequest)
		return
	}
	userRet := AdminReturn{ID:id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email }
	helpers.WriteJSON(w,http.StatusAccepted,userRet)

}

type UserReturn struct{
	ID int64 `json:"ID"`
	ProfileID int64 `json:"ProfileID"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	Email string `json:"Email"`

}

func (route *Routes) Register(w http.ResponseWriter, r *http.Request){
	db := route.DB

	user := User{}
	helpers.ReadJSON(w, r, &user)
	// passByte, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash),bcrypt.DefaultCost)
	// if err != nil{
	// 	fmt.Println("Password Gen issue", err)
	// }
	id,profId, err := route.UserQuery.RegisterUserIntoDB(db,user.Password,user.FirstName,user.LastName,user.Email)
	if err != nil{
		fmt.Println(err)
		helpers.ErrorJSON(w,err,http.StatusBadRequest)
		return
	}

	userRet := UserReturn{ID:id, ProfileID:profId, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email }
	helpers.WriteJSON(w,http.StatusAccepted,userRet)
	
}

type LoginUser struct{
	Email string `email:"Email"`
	Password string `json:"PasswordHash"`
}

type SendBackLogin struct{
	Token string `jwt:"Email"`
}
func (route *Routes) Login(w http.ResponseWriter, r *http.Request){
	db := route.DB
	login := LoginUser{}
	helpers.ReadJSON(w, r, &login)
	_, passwordStored, userID,err := route.UserQuery.LoginUserDB(db, login.Email)
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
		"iat":time.Now().Unix(),
		"admin":"False",
		"email":login.Email,
		"userId":userID,
	})
	// Remove the testing key for this
	tokenString, err := token.SignedString([]byte("Testing key"))
	sendBack := SendBackLogin{Token: tokenString}
	if err != nil{
		fmt.Println("signed token error")
		fmt.Println(err)
	}

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
	userID := r.Context().Value("userid")
	UserProfile := &UserProfile{}
	cell, home, err := route.UserQuery.GetUserProfile(route.DB, userID)


	if err != nil{
		fmt.Println("Error with getting userprofile in users.go")
	}
	
	UserProfile.Cell = cell
	UserProfile.Home = home

	helpers.WriteJSON(w,http.StatusAccepted, &UserProfile)


}