package productendpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Apouzi/golang-shop/app/api/database"
	"github.com/Apouzi/golang-shop/app/api/helpers"
	"github.com/go-chi/chi"
)

type ProductRoutes struct{
	DB *sql.DB
	ProductQuery *database.PrepareStatmentsProducts
}

func InstanceProductsRoutes(db *sql.DB ) *ProductRoutes {
	r := &ProductRoutes{
		DB: db,
		ProductQuery: database.InitPrepare(db),
	}
	return r
}

func (route *ProductRoutes) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {
	ProdJSON := route.ProductQuery.GetAllProducts(route.DB)
	JSONWrite,err := json.Marshal(ProdJSON)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Failed")))
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(JSONWrite)
	
}


func (route *ProductRoutes) GetOneProductsEndPoint(w http.ResponseWriter, r *http.Request){
	query, err :=  strconv.Atoi(chi.URLParam(r,"Product_ID"))
	if err != nil{
		fmt.Println("String to Int failed:", err)
	}
	ProdJSON, err := route.ProductQuery.GetOneProduct(route.DB,query)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted,ProdJSON)

}

// func (route *Routes) GetProductCategoryEndPoint(w http.ResponseWriter, r *http.Request){
// 	category, err := strconv.Atoi(chi.URLParam(r, "CategoryName"))
	
// 	if err != nil{
// 		fmt.Println("Get Product Category ")
// 	}

// 	ProdJSON := route.ProductQuery.GetProductCategoryFinal(route.DB,category)
// 	JSONWrite, err := json.Marshal(ProdJSON)

// 	if err != nil{
// 		fmt.Println(err)
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(fmt.Sprint("Failed")))
// 	}

// 	w.WriteHeader((http.StatusAccepted))
// 	w.Write(JSONWrite)

// }


func (route *ProductRoutes) GetProductCategoryEndPointFinal(w http.ResponseWriter, r *http.Request){
	category := chi.URLParam(r, "CategoryName")

	// if err != nil{
	// 	fmt.Println("Get Product Category ")
	// }
// TODO needs error handling for none existent categories!
	ProdJSON := route.ProductQuery.GetProductCategoryFinal(route.DB,category)
	JSONWrite, err := json.Marshal(ProdJSON)

	if err != nil{
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Failed")))
	}

	w.WriteHeader((http.StatusAccepted))
	w.Write(JSONWrite)

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




