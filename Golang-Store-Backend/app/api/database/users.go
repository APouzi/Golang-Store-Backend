package database

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var data sql.DB

type UserStatments struct{
	RegisterUser *sql.Stmt
	GetUser *sql.Stmt
	GetUserAndProfile *sql.Stmt
	GetGetUserProfileWishListWishListProduct *sql.Stmt
}

func InitUserStatments(db *sql.DB) *UserStatments {
	prep := &UserStatments{}
	var err error
	prep.RegisterUser, err = db.Prepare("INSERT INTO tblUser(PasswordHash,FirstName, LastName, Email) VALUES(?,?,?,?)")
	if err != nil{
		fmt.Println("Prepare statment broken", err)
	}

	return prep
}

func(stmt *UserStatments)  RegisterUserIntoDB(Password string, firstName string, lastName string, Email string) int{
	
	
	passByte, err := bcrypt.GenerateFromPassword([]byte(Password),bcrypt.DefaultCost)
	if err != nil{
		fmt.Println("Password Gen issue", err)
	}

	response, err := stmt.RegisterUser.Exec(passByte,firstName, lastName, Email)
	if err != nil{
		fmt.Println("Registering User Into DB Error:", err)
	}
	id, err := response.LastInsertId()
	if err != nil {
		// handle error
	}

	return int(id)
}
// func(models *Models) getAll() *Customer{
	
// } 