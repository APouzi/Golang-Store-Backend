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
	GetUserProfileStmt *sql.Stmt
	GetGetUserProfileWishListWishListProduct *sql.Stmt
}

func InitUserStatments(db *sql.DB) *UserStatments {
	prep := &UserStatments{}
	var err error
	prep.RegisterUser, err = db.Prepare("INSERT INTO tblUser(PasswordHash,FirstName, LastName, Email) VALUES(?,?,?,?)")
	if err != nil{
		fmt.Println("Prepare statment broken", err)
	}
	prep.GetUserProfileStmt, err = db.Prepare("SELECT UserProfileID, PhoneNumberCell, PhoneNumberHome FROM tblUserProfile WHERE UserID = ?")
	if err != nil{
		fmt.Println("Prepare statement err", err)
	}
	prep.GetUser, err = db.Prepare("SELECT UserID,Email, PasswordHash FROM tblUser where Email = ?")
	if err != nil{
		fmt.Println("Prepare statement err", err)
	}
	return prep
}

func(stmt *UserStatments)  RegisterUserIntoDB(db *sql.DB,Password string, firstName string, lastName string, Email string) (int64, int64, error){
	sqlStmt := "SELECT COUNT(email) FROM tblUser WHERE Email = ?"
	var accCheck int
	var failedtransactionerror error
	row := db.QueryRow(sqlStmt, Email)
	err := row.Scan(&accCheck)
	if err != nil{
		fmt.Println("scanning error in RegisterUSerIntoDB", err)
		failedtransactionerror = fmt.Errorf("this user already exists")
		return -1, -1,failedtransactionerror
	}
	
	if accCheck != 0{
		fmt.Println("This user already exists")
		failedtransactionerror = fmt.Errorf("this user already exists")
		return -1, -1,failedtransactionerror
	}
	passByte, err := bcrypt.GenerateFromPassword([]byte(Password),bcrypt.DefaultCost)
	if err != nil{
		fmt.Println("Password Gen issue", err)
	}
	tx, err := db.Begin()
	defer func(){
		if err != nil{
			fmt.Println("Transaction for profile failed", err)
			tx.Rollback()
			return
		}
	}()
	queryUser := "INSERT INTO tblUser(PasswordHash,FirstName, LastName, Email) VALUES(?,?,?,?)"
	response, err := tx.Exec(queryUser, passByte,firstName, lastName, Email)

	id, err := response.LastInsertId()
	if err != nil {
		return 0, -1,err
	}
	queryProfile := "INSERT INTO tblUserProfile(UserID, PhoneNumberCell, PhoneNumberHome) VALUES(?,?,?)"

	profRes,err := tx.Exec(queryProfile,id, "33333","44444")

	if err != nil{
		fmt.Println("Registering User Into DB Error:", err)
	}
	ProfID,err := profRes.LastInsertId()
	if err != nil{
		fmt.Println("Profile ID Error:", err)
	}
	tx.Commit()
	return id, ProfID,nil
}


func(stmt *UserStatments)  RegisterAdminIntoDB(db *sql.DB,Password string, firstName string, lastName string, Email string) (int64, error){
	sqlStmt := "SELECT COUNT(email) FROM tblUser WHERE Email = ?"
	var accCheck int
	var failedtransactionerror error
	row := db.QueryRow(sqlStmt, Email)
	err := row.Scan(&accCheck)
	if err != nil{
		fmt.Println("scanning error in RegisterUSerIntoDB", err)
		failedtransactionerror = fmt.Errorf("this user already exists")
		return -1, failedtransactionerror
	}
	
	if accCheck != 0{
		fmt.Println("This user already exists")
		failedtransactionerror = fmt.Errorf("this user already exists")
		return -1, failedtransactionerror
	}
	passByte, err := bcrypt.GenerateFromPassword([]byte(Password),bcrypt.DefaultCost)
	if err != nil{
		fmt.Println("Password Gen issue", err)
	}
	tx, err := db.Begin()
	defer func(){
		if err != nil{
			fmt.Println("Transaction for profile failed", err)
			tx.Rollback()
			return
		}
	}()
	queryUser := "INSERT INTO tblUser(PasswordHash,FirstName, LastName, Email) VALUES(?,?,?,?)"
	response, err := tx.Exec(queryUser, passByte,firstName, lastName, Email)

	id, err := response.LastInsertId()
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, nil
}


func (stmt *UserStatments) LoginUserDB(db *sql.DB, email string)(string, string, int64,error){
	var userID int64
	var emailTwo string
	var password string
	row := stmt.GetUser.QueryRow(email).Scan(&userID,&emailTwo,&password)

	if row == sql.ErrNoRows{
		fmt.Println("email doesn't exist")
		err := fmt.Errorf("email doesn't exist")
		return "", "", -1,err
	}

	return emailTwo, password,userID, nil

}

func (stmt *UserStatments) GetUserProfile(db *sql.DB, userProfileID any)(int, int, error){
	var UserProfileID int
	var phoneNumCell int
	var phoneNumHome int
	row := stmt.GetUserProfileStmt.QueryRow(userProfileID).Scan(&UserProfileID,&phoneNumCell, &phoneNumHome)
	if row == sql.ErrNoRows{
		fmt.Println("profile doesn't exist")
		err := fmt.Errorf("profile doesn't exist")
		return -1, -1, err
	}
	return phoneNumCell,phoneNumHome,nil
}