package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-jose/go-jose/v3/json"
)

func (route *Routes) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {
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

func (route *Routes) GetOneProductsEndPoint(w http.ResponseWriter, r *http.Request){
	query :=  chi.URLParam(r,"ProductID")
	queryToInt, err := strconv.Atoi(query)
	if err != nil{
		fmt.Println("String to Int failed:", err)
	}
	ProdJSON := route.ProductQuery.GetOneProduct(route.DB,queryToInt)
	JSONWrite,err := json.Marshal(ProdJSON)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint("Failed")))
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(JSONWrite)
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


func (route *Routes) GetProductCategoryEndPointFinal(w http.ResponseWriter, r *http.Request){
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

func (route *Routes) CreateTestCategory(w http.ResponseWriter, r *http.Request){
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
	
	PrimeSub := "INSERT INTO tblCatFinalProd(CatFinalID, ProductID) VALUES(?,?)"
	SubFinal := "INSERT INTO tblCatPrimeSub(CatPrimeID, CatSubID) VALUES(?,?)"
	FinalProd := "INSERT INTO tblCatSubFinal(CatSubID, CatFinalID) VALUES(?,?)"
	idPrimeR,err := idPrime.LastInsertId()
	idSubR, err := idSub.LastInsertId()
	idFinalR, err := idFinal.LastInsertId()
	tx.Exec(PrimeSub, idFinalR, 1)
	tx.Exec(SubFinal, idPrimeR,idSubR)
	tx.Exec(FinalProd, idSubR, idFinalR)
	tx.Commit()
}

func (route *Routes) PullTestCategory(w http.ResponseWriter, r *http.Request){
	query := "SELECT * FROM tblProducts 
	JOIN tblCatFinalProd ON tblCatFinalProd.ProductID = tblProducts.ProductID 
	JOIN tblCategoriesFinal ON tblCategoriesFinal.CategoryID = tblCatFinalProd.CategoryID 
	JOIN tblCategoriesFinal ON tblCategoriesFinal.CategoryID = tblCategoriesFinal.CatFinalID 
	JOIN tblCategoriesSub ON tblCategoriesSub.CategoryID = tblCatSubFinal."
}