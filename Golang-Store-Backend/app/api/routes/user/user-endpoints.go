package userendpoints

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Apouzi/golang-shop/app/api/database"
	"github.com/Apouzi/golang-shop/app/api/helpers"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRoutes struct{
	DB *sql.DB
	UserQuery *database.UserStatments
}

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

func InstanceUserRoutes(db *sql.DB ) *UserRoutes {
	r := &UserRoutes{
		DB: db,
		UserQuery: database.InitUserStatments(db),
	}
	return r
}

func (route *UserRoutes) AdminSuperUserCreation(w http.ResponseWriter, r *http.Request){
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

func (route *UserRoutes) Register(w http.ResponseWriter, r *http.Request){
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

func (route *UserRoutes) Login(w http.ResponseWriter, r *http.Request){
	db := route.DB
	login := LoginUser{}
	helpers.ReadJSON(w, r, &login)
	_, passwordStored, userID,err := route.UserQuery.LoginUserDB(db, login.Email)
	var errRet error
	if err != nil{
		errRet = fmt.Errorf("User does not exist")
		helpers.ErrorJSON(w,errRet,http.StatusBadRequest)
		return
	}


	err = bcrypt.CompareHashAndPassword([]byte(passwordStored), []byte(login.Password))

	if err !=nil{
		errRet = fmt.Errorf("password does not match")
		helpers.ErrorJSON(w,errRet,http.StatusBadRequest)
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
		errRet = fmt.Errorf("server issue, cannot send token")
		helpers.ErrorJSON(w,errRet,http.StatusBadRequest)
		return
	}

	helpers.WriteJSON(w, http.StatusAccepted, &sendBack)
}

type UserProfile struct{
	Cell int `json:"Cell"`
	Home int `json:"Home"`
}

func (route *UserRoutes) UserProfile(w http.ResponseWriter, r *http.Request){
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