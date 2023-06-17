package main

import (
	"database/sql"
	"flag"
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
	// argsRetrieve := os.Args[0]
	// if argsRetrieve == "initProd"{
	// 	fmt.Println("It did it!")
	// }

	connection, models := initDB()

	// flags to initailize this
	var initializeDB, initailizeView string

	flag.StringVar(&initializeDB, "initdb","","Initalize Database")
	flag.StringVar(&initailizeView,"initView","","Intialize Views")
	flag.Parse()
	
	app := Config{
		DB: connection,
		Models: models,
	}
	
	if initializeDB == "t" || initializeDB == "T"{
		PopulateProductTables(app.DB)
		InitateAndPopulateUsers(app.DB)
	}
	if initailizeView == "t" || initailizeView == "T"{
		IntializeViews(app.DB)
	}
	// if TestInitCreateThenDelete(app.DB) == false{
	// 	log.Fatal("Connection Test had failed")
	// }

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
	
	query, err := ioutil.ReadFile("./sql/Products.sql")

	if err != nil{
		log.Fatal("Error when loading sql file",err)
	}

	_, err = db.Exec(string(query))

	if err != nil{
		log.Fatal("Couldn't complete the execution of the file", err)
	}

	for i := 0.00; i <= 10; i++{
		_,err = db.Exec("INSERT INTO tblProducts (Product_Name, Product_Description) VALUES(?,?)", "testProductPopulate","This is a description!")
		if err != nil{
			log.Fatal("Error with tblProducts")
		}
	}
	// _, err = db.Exec("INSERT INTO tblCategoriesPrime(CategoryName, CategoryDescription) VALUES(?,?)", "Test Category","This is a description category")
	// if err != nil{
	// 	log.Println("Error inserting tblCategoriesPrime")
	// }
	// _,err =  db.Exec("INSERT INTO tblProductsCategoriesPrime(Product_ID, Category_ID) VALUES(1,1)")
	// if err != nil{
	// 	log.Println("Error inserting into tblProductsCategoriesPrime", err)
	// }
	// _,err =  db.Exec("INSERT INTO tblProductsCategoriesPrime(Product_ID, Category_ID) VALUES(2,1)")
	// if err != nil{
	// 	log.Println("Error inserting into tblProductCategoies")
	// }

	resultProd := database.Product{}
	row := db.QueryRow("select Product_ID, Product_Name, Product_Description from tblProducts where Product_ID = ?",4)
	if row == nil{
		fmt.Println("Nothing returned!")
	}
	err = row.Scan(&resultProd.Product_ID,&resultProd.Product_Name,&resultProd.Product_Description)
	if err != nil {
		fmt.Println(err)
	}

	listPrint:= []database.Product{}
	// rows, err := db.Query("SELECT tblProducts.Product_ID, tblProducts.Product_Name, tblProducts.Product_Description, tblProducts.Product_Price FROM tblProducts JOIN tblProductsCategoriesPrime ON tblProducts.Product_ID = tblProductsCategoriesPrime.Product_ID JOIN tblCategoriesPrime ON tblProductsCategoriesPrime.Category_ID = tblCategoriesPrime.Category_ID WHERE tblCategoriesPrime.CategoryName = ?", "Test Category" )
	if err != nil{
		log.Fatal("Error with category", err)
	}
	// defer rows.Close()
	// for rows.Next(){
	// 	resultProd2 := database.Product{}
	// 	rows.Scan(&resultProd2.Product_ID, &resultProd2.Product_Name, &resultProd2.Product_Description, &resultProd2.Product_Price)
	// 	listPrint = append(listPrint, resultProd2)
	// }
	fmt.Println(resultProd)
	fmt.Println("Population of tables has been completed!")
	fmt.Println("Categories:",listPrint)
	for _,v := range listPrint{
		fmt.Println(v.Product_Name)
	}
	//Test this out!

}

func TestCategories( db *sql.DB){
	
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

func InitateAndPopulateUsers(db *sql.DB){
	query, err := ioutil.ReadFile("./sql/User.sql")

	if err != nil{
		log.Fatal("Error when loading sql file",err)
	}

	_, err = db.Exec(string(query))
	if err != nil{
		fmt.Println("SQL user failed to initalized")
		log.Fatalln(err)
	}

	_, err = db.Exec("INSERT INTO tblUser (FirstName, LastName, Email) VALUES(?,?,?)","TestFirstName","TestLastName", "TestEmail@email.com" )
	if err != nil{
		fmt.Println("Problem with query to insert in tblUser")
	}
	_, err = db.Exec("INSERT INTO tblUserProfile (UserID, PhoneNumberCell, PhoneNumberHome) VALUES(?,?,?)", 1, "6195555555","8585555555" )
	if err != nil{
		fmt.Println("Problem with query to insert in tblUserProfile")
	}
	row := db.QueryRow("SELECT UserID, FirstName, LastName, Email FROM tblUser WHERE UserId = 1")

	rowProfile := db.QueryRow("SELECT PhoneNumberCell, PhoneNumberHome FROM tblUserProfile JOIN tblUser ON tblUserProfile.UserID = tblUser.UserID where tblUser.UserID = 1")



	user := user{}
	row.Scan(&user.UserID, &user.FirstName,&user.LastName,&user.Email)
	userProf := userProfile{}
	rowProfile.Scan(&userProf.PhoneNumberCell, &userProf.PhoneNumberHome)
	fmt.Println("users:",user)
	fmt.Println("userProfile:",userProf)
}

func IntializeViews(db *sql.DB) {
	
	query, err := ioutil.ReadFile("./sql/Views.sql")

	if err != nil{
		log.Fatal("Error when loading sql file",err)
	}

	_, err = db.Exec(string(query))
	if err != nil{
		log.Fatal("IntializeViews query execution failed", err)
	}
}