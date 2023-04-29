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

func(stmt *UserStatments)  RegisterUserIntoDB(db *sql.DB,Password string, firstName string, lastName string, Email string) (int64, error){
	sqlStmt := "SELECT email FROM tblUser WHERE user = email"
	row := db.QueryRow(sqlStmt, Email).Scan(&Email)
	if row == sql.ErrNoRows{
		fmt.Println("This user already exists")
		return -1, nil
	}
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
		return 0, err
	}

	// if 

	return id, nil
}
// func(models *Models) getAll() *Customer{
	
// } 