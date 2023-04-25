package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Apouzi/golang-shop/app/api/database"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const webport = 8000

type Config struct{
	DB *sql.DB
	Models *database.Models
}


func main(){
	connection, models := initDB()
	app := Config{
		DB: connection,
		Models: models,
	}
	fmt.Printf("Starting Store Backend on port %d \n", webport)

	serve := &http.Server{
		Addr: fmt.Sprintf(":%d", webport),
		Handler: app.StartRouter(),
	}

	err := serve.ListenAndServe()
	if err != nil{
		log.Panic(err)
	}
	
	// fmt.Println("test", reflect.TypeOf(router))
}


// Initializing the environment variables to run. 
func init() {
    // Get the absolute path to the directory where the executable is located
    exeDir, err := filepath.Abs("./")
    if err != nil {
        log.Fatal(err)
    }

    // Load the .env file from the directory
    err = godotenv.Load(filepath.Join(exeDir, ".env"))
    if err != nil {
        log.Fatal("Error loading .env file\n","exeDir: ", exeDir)
    }
}

func initDB() (*sql.DB,*database.Models){
	cfg := mysql.Config{
		User:   "user",
		Passwd: "example",
		Net:    "tcp",
		Addr:   "mysql:3306",
		DBName: "database",
		MultiStatements: true,
	}
	var db *sql.DB
	var err error
	count := 0

	
	for count < 11{
		db, err = sql.Open("mysql", cfg.FormatDSN(),)
		count++
		if err != nil{
			fmt.Printf("MySQL is still waiting to connect, trying to connect again. Attempt: %d \n", count)
		} else if err = db.Ping(); err == nil {
			fmt.Println("MySQL server connected confirmation")
				break
		}
		fmt.Printf("Attempt: %d connecting to MySQL server again",count)
		
		time.Sleep(2 * time.Second)
		
	}
	
	if TestInitCreateThenDelete(db) == false{
		log.Fatal("Connection Test had failed")
	}

	PopulateProductTables(db)
	PopulateTestUsers(db)
	

	database := &database.Models{}
	return db , database
}


func TestInitCreateThenDelete(db *sql.DB) bool{
	fmt.Println("Starting connection query test to MySQL")
	fmt.Println("Creating Table started")
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS tblTEST(id INT AUTO_INCREMENT PRIMARY KEY,name VARCHAR(14))")
	if err != nil{
		fmt.Println("Create fail", err)
		return false
	}

	fmt.Println("Created Table. Starting Insertion to table")
	_, err = db.Exec("INSERT INTO tblTEST(id, name) VALUES (1, 'alex')")
	if err != nil{
		fmt.Println("Insert fail", err)
		return false
	}

	fmt.Println("Table Insertion Completed. Starting Query")
	var name string
	var id int
	err = db.QueryRow("SELECT * FROM tblTEST WHERE id = ?",1).Scan(&id,&name)
	if err != nil{
		fmt.Println("Error at QueryRow",err)
		return false
	}

	fmt.Println("Query Completed. Expectation is: alex 1.\n The results:",name,id)
	if err != nil{
		fmt.Println(err)
		return false
	}

	fmt.Println("Deletion of table to complete test")
	_, err = db.Exec("DELETE FROM tblTEST WHERE id = 1")
	if err != nil{
		fmt.Println("Table deletion fail", err)
		return false
	}

	return true
}

func PopulateProductTables(db *sql.DB) {
	
	query, err := ioutil.ReadFile("./sql/CreateTables&Rows.sql")

	if err != nil{
		log.Fatal("Error when loading sql file",err)
	}

	_, err = db.Exec(string(query))

	if err != nil{
		log.Fatal("Couldn't complete the execution of the file", err)
	}

	for i := 0.00; i <= 10; i++{
		_,err = db.Exec("INSERT INTO tblProducts (ProductName, ProductDescription, ProductPrice, SKU, UPC) VALUES(?,?,?,?,?)", "testProductPopulate","This is a description!",10.85+i,"SKUABC123","21124214311A")
		if err != nil{
			log.Fatal("Error with tblProducts")
		}
	}

	_, err = db.Exec("INSERT INTO tblCategoriesPrime(CategoryName, CategoryDescription) VALUES(?,?)", "Test Category","This is a description category")
	if err != nil{
		log.Fatal("Error inserting tblCategoriesPrime")
	}
	_,err =  db.Exec("INSERT INTO tblProductsCategoriesPrime(ProductID, CategoryID) VALUES(1,1)")
	if err != nil{
		log.Fatal("Error inserting into tblProductsCategoriesPrime", err)
	}
	_,err =  db.Exec("INSERT INTO tblProductsCategoriesPrime(ProductID, CategoryID) VALUES(2,1)")
	if err != nil{
		log.Fatal("Error inserting into tblProductCategoies")
	}

	resultProd := database.Product{}
	row := db.QueryRow("select ProductID, ProductName, ProductDescription, ProductPrice from tblProducts where ProductID = ?",4)
	if row == nil{
		fmt.Println("Nothing returned!")
	}
	err = row.Scan(&resultProd.ProductID,&resultProd.ProductName,&resultProd.ProductDescription, &resultProd.ProductPrice)
	if err != nil {
		fmt.Println(err)
	}

	listPrint:= []database.Product{}
	rows, err := db.Query("SELECT tblProducts.ProductID, tblProducts.ProductName, tblProducts.ProductDescription, tblProducts.ProductPrice FROM tblProducts JOIN tblProductsCategoriesPrime ON tblProducts.ProductID = tblProductsCategoriesPrime.ProductID JOIN tblCategoriesPrime ON tblProductsCategoriesPrime.CategoryID = tblCategoriesPrime.CategoryID WHERE tblCategoriesPrime.CategoryName = ?", "Test Category" )
	if err != nil{
		log.Fatal("Error with category", err)
	}
	defer rows.Close()
	for rows.Next(){
		resultProd2 := database.Product{}
		rows.Scan(&resultProd2.ProductID, &resultProd2.ProductName, &resultProd2.ProductDescription, &resultProd2.ProductPrice)
		listPrint = append(listPrint, resultProd2)
	}
	fmt.Println(resultProd)
	fmt.Println("Population of tables has been completed!")
	fmt.Println("Categories:",listPrint)
	for _,v := range listPrint{
		fmt.Println(v.ProductName)
	}
	//Test this out!

}
type user struct{
	UserID int
	FirstName string
	LastName string
	Email string
}

type userProfile struct{
	PhoneNumberCell string
	PhoneNumberHome string
}

func PopulateTestUsers(db *sql.DB){
	query, err := ioutil.ReadFile("./sql/User.sql")

	if err != nil{
		log.Fatal("Error when loading sql file",err)
	}

	_, err = db.Exec(string(query))

	_, err = db.Exec("INSERT INTO tblUser (FirstName, LastName, Email) VALUES(?,?,?)","TestFirstName","TestLastName", "TestEmail@email.com" )

	_, err = db.Exec("INSERT INTO tblUserProfile (UserID, PhoneNumberCell, PhoneNumberHome) VALUES(?,?,?)", 1, "6195555555","8585555555" )

	row := db.QueryRow("SELECT UserID, FirstName, LastName, Email FROM tblUser WHERE UserId = 1")

	rowProfile := db.QueryRow("SELECT PhoneNumberCell, PhoneNumberHome FROM tblUserProfile JOIN tblUser ON tblUserProfile.UserID = tblUser.UserID where tblUser.UserID = 1")



	user := user{}
	row.Scan(&user.UserID, &user.FirstName,&user.LastName,&user.Email)
	userProf := userProfile{}
	rowProfile.Scan(&userProf.PhoneNumberCell, &userProf.PhoneNumberHome)
	fmt.Println("users:",user)
	fmt.Println("userProfile:",userProf)
}