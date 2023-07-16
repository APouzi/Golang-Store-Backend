package testroutes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Apouzi/golang-shop/app/api/database"
	"github.com/Apouzi/golang-shop/app/api/helpers"
	"github.com/redis/go-redis/v9"
)

type ProductRoutes struct {
	DB           *sql.DB
	ProductQuery *database.PrepareStatmentsProducts
	Redis        *redis.Client
}

func InjectDBRef(db *sql.DB, redis *redis.Client) *ProductRoutes {
	r := &ProductRoutes{
		DB:           db,
		ProductQuery: database.InitPrepare(db),
		Redis:        redis,
	}
	return r
}

func (route *ProductRoutes) GetOneProductSQL(w http.ResponseWriter, r *http.Request) {
	rows := route.DB.QueryRow("SELECT Product_ID, Product_Name, Product_Description, PRIMARY_IMAGE, Date_Created, Modified_Date FROM tblProducts WHERE Product_ID = 1")
	prodJSON := ProductJSON{}

	err := rows.Scan(
		&prodJSON.Product_ID,
		&prodJSON.Product_Name,
		&prodJSON.Product_Description,
		&prodJSON.PRIMARY_IMAGE,
		&prodJSON.ProductDateAdded,
		&prodJSON.ModifiedDate,
	)
	if err == sql.ErrNoRows {
		helpers.ErrorJSON(w, errors.New("failed error"))
		return
	}
	if err != nil {
		fmt.Println("scanning error:", err)
	}

	helpers.WriteJSON(w, http.StatusAccepted, prodJSON)
}

type ProductJSON struct {
	Product_ID          int    `json:"Product_ID"`
	Product_Name        string `json:"Product_Name"`
	Product_Description string `json:"Product_Description"`
	PRIMARY_IMAGE       string `json:"PRIMARY_IMAGE,omitempty"`
	ProductDateAdded    string `json:"DateAdded"`
	ModifiedDate        string `json:"ModifiedDate"`
}

func (route *ProductRoutes) GetOneProductRedis(w http.ResponseWriter, r *http.Request) {

	result, err := route.Redis.Get(context.Background(), "products/1").Result()

	if err == nil {
		helpers.WriteJSON(w, 200, result)
		return
	}

	rows := route.DB.QueryRow("SELECT Product_ID, Product_Name, Product_Description, Date_Created, Modified_Date FROM tblProducts WHERE Product_ID = 1")
	prodJSON := ProductJSON{}

	err = rows.Scan(
		&prodJSON.Product_ID,
		&prodJSON.Product_Name,
		&prodJSON.Product_Description,
		// &prodJSON.PRIMARY_IMAGE,
		&prodJSON.ProductDateAdded,
		&prodJSON.ModifiedDate,
	)
	if err == sql.ErrNoRows {
		helpers.ErrorJSON(w, errors.New("failed error"))
		return
	}
	if err != nil {
		fmt.Println("scanning error:", err)
	}
	test, err := json.Marshal(prodJSON)
	err = route.Redis.Set(context.Background(), "products/1", test, 0).Err()
	if err != nil {
		fmt.Println(err)
		fmt.Println("failed to save redis")
	}

	helpers.WriteJSON(w, http.StatusAccepted, prodJSON)
}

func (route *ProductRoutes) PullTestCategory(w http.ResponseWriter, r *http.Request){
	// 	JOIN tblCategoriesFinal ON tblCategoriesFinal.Category_ID = tblCategoriesFinal.CatFinalID 
	// query := "SELECT tblProducts.Product_ID, tblProducts.Product_Name FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.Product_ID = tblProducts.Product_ID JOIN tblCategoriesFinal ON tblCategoriesFinal.Category_ID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.Category_ID JOIN tblCategoriesSub ON tblCategoriesSub.Category_ID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.Category_ID JOIN tblCategoriesPrime ON tblCategoriesPrime.Category_ID = tblCatPrimeSub.CatPrimeID WHERE tblProducts.Product_ID = ?"
	// query := "SELECT tblProducts.Product_ID, tblProducts.Product_Name FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.Product_ID = tblProducts.Product_ID JOIN tblCategoriesFinal ON tblCategoriesFinal.Category_ID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.Category_ID JOIN tblCategoriesSub ON tblCategoriesSub.Category_ID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.Category_ID JOIN tblCategoriesPrime ON tblCategoriesPrime.Category_ID = tblCatPrimeSub.CatPrimeID"
	// query := "SELECT tblProducts.Product_ID, tblProducts.Product_Name FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.Product_ID = tblProducts.Product_ID JOIN tblCategoriesFinal ON tblCategoriesFinal.Category_ID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.Category_ID JOIN tblCategoriesSub ON tblCategoriesSub.Category_ID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.Category_ID JOIN tblCategoriesPrime ON tblCategoriesPrime.Category_ID = tblCatPrimeSub.CatPrimeID WHERE tblCategoriesPrime.CategoryName = ?"
	query := "SELECT Product_ID, Product_Name FROM PrimeSubFinalCategoryProducts where CategoryName = ?"
	type RowReadTest struct{
		Product_ID int
		Product_Name string
	}
	// row := route.DB.QueryRow(query2)
	readinto := RowReadTest{}
	row, err := route.DB.Query(query, "PrimeTest")

	// err:= row.Scan(&readinto.Product_ID, &readinto.Product_Name)
	if err != nil{
		fmt.Println("err with row in PullTestCategory, error below")
		fmt.Println(err)
		return
	}
	for row.Next(){
		err := row.Scan(&readinto.Product_ID, &readinto.Product_Name)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(readinto.Product_ID, readinto.Product_Name)
	}
	fmt.Println("PullTestCAtegory result is:", readinto.Product_ID,readinto.Product_Name)
}



// Note, in admin we want to be able to query all products based on the hierarchy level, it would be more performant to write products and tables to every single level but only read based on that singular leve instead of joining like 3 different tables if we are getting the last level of tables. 
//An example: catPrime join catSub join catFinal then pull data based on that. 
// So instead, when we pull some data, we need to use a transaction to input into the 
func (route *ProductRoutes) CreateTestCategory(w http.ResponseWriter, r *http.Request){
	tx, err := route.DB.Begin()

	if err != nil{
		fmt.Println("TextCategories transaction intialization failed")
	}

	prime_category :=  "INSERT INTO tblCategoriesPrime(CategoryName, CategoryDescription) VALUES(?,?)"
	sub_cat := "INSERT INTO tblCategoriesSub(CategoryName, CategoryDescription) VALUES(?,?)"
	final_cat := "INSERT INTO tblCategoriesFinal(CategoryName, CategoryDescription) VALUES(?,?)"


	idPrime, err := tx.Exec(prime_category, "PrimeTest","this is a test category")
	if err != nil{
		fmt.Println("Issue with Prime transaction")
	}
	idSub, err := tx.Exec(sub_cat, "SubTest","this is a test category")
	if err != nil{
		fmt.Println("Issue with SubTest transaction")
	}
	idFinal, err := tx.Exec(final_cat, "FinalTest","this is a test category")
	if err != nil{
		fmt.Println("Issue with FinalTest transaction")
	}

	PrimeSub:= "INSERT INTO tblCatPrimeSub(CatPrimeID, CatSubID) VALUES(?,?)"
	FinalProd := "INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)"
	SubFinal := "INSERT INTO tblCatSubFinal(CatSubID, CatFinalID) VALUES(?,?)"
	
	idPrimeR,err := idPrime.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	idSubR, err := idSub.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	idFinalR, err := idFinal.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	tx.Exec(FinalProd, idFinalR, 1)
	tx.Exec(PrimeSub, idPrimeR,idSubR)
	tx.Exec(SubFinal, idSubR, idFinalR)
	tx.Commit()
}