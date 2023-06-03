package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Apouzi/golang-shop/app/api/helpers"
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

func(route *Routes) CreateProduct(w http.ResponseWriter, r *http.Request){

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

	PrimeSub:= "INSERT INTO tblCatPrimeSub(CatPrimeID, CatSubID) VALUES(?,?)"
	FinalProd := "INSERT INTO tblCatFinalProd(CatFinalID, ProductID) VALUES(?,?)"
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

func (route *Routes) PullTestCategory(w http.ResponseWriter, r *http.Request){
	// 	JOIN tblCategoriesFinal ON tblCategoriesFinal.CategoryID = tblCategoriesFinal.CatFinalID 
	// query := "SELECT tblProducts.ProductID, tblProducts.ProductName FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.ProductID = tblProducts.ProductID JOIN tblCategoriesFinal ON tblCategoriesFinal.CategoryID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.CategoryID JOIN tblCategoriesSub ON tblCategoriesSub.CategoryID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.CategoryID JOIN tblCategoriesPrime ON tblCategoriesPrime.CategoryID = tblCatPrimeSub.CatPrimeID WHERE tblProducts.ProductID = ?"
	// query := "SELECT tblProducts.ProductID, tblProducts.ProductName FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.ProductID = tblProducts.ProductID JOIN tblCategoriesFinal ON tblCategoriesFinal.CategoryID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.CategoryID JOIN tblCategoriesSub ON tblCategoriesSub.CategoryID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.CategoryID JOIN tblCategoriesPrime ON tblCategoriesPrime.CategoryID = tblCatPrimeSub.CatPrimeID"
	// query := "SELECT tblProducts.ProductID, tblProducts.ProductName FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.ProductID = tblProducts.ProductID JOIN tblCategoriesFinal ON tblCategoriesFinal.CategoryID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.CategoryID JOIN tblCategoriesSub ON tblCategoriesSub.CategoryID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.CategoryID JOIN tblCategoriesPrime ON tblCategoriesPrime.CategoryID = tblCatPrimeSub.CatPrimeID WHERE tblCategoriesPrime.CategoryName = ?"
	query := "SELECT ProductID, ProductName FROM PrimeSubFinalCategoryProducts where CategoryName = ?"
	type RowReadTest struct{
		ProductID int
		ProductName string
	}
	// row := route.DB.QueryRow(query2)
	readinto := RowReadTest{}
	row, err := route.DB.Query(query, "PrimeTest")

	// err:= row.Scan(&readinto.ProductID, &readinto.ProductName)
	if err != nil{
		fmt.Println("err with row in PullTestCategory, error below")
		fmt.Println(err)
		return
	}
	for row.Next(){
		err := row.Scan(&readinto.ProductID, &readinto.ProductName)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(readinto.ProductID, readinto.ProductName)
	}
	fmt.Println("PullTestCAtegory result is:", readinto.ProductID,readinto.ProductName)
}



// Admin functionality

type CategoryInsert struct{
	CategoryName string `json:"CategoryName"`
	CategoryDescription string `json:"CategoryDescription"`
}
func (route *Routes) CreatePrimeCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCategoriesPrime(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription )
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	
	helpers.WriteJSON(w, http.StatusAccepted, resultID)
}

func (route *Routes) CreateSubCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCategoriesSub(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription )
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	
	helpers.WriteJSON(w, http.StatusAccepted, resultID)
}

func (route *Routes) CreateFinalCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCategoriesFinal(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription )
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	
	helpers.WriteJSON(w, http.StatusAccepted, resultID)
}

type CatToCat struct {
	CatStart int `json:"CategoryStart"`
	CatEnd int `json:"CategoryEnd"`
}


func (route *Routes) ConnectPrimeToSubCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w,r, &FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCatPrimeSub(CatPrimeID,  CatSubID) VALUES(?,?)", FinalSub.CatStart, FinalSub.CatEnd)
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	
	helpers.WriteJSON(w, http.StatusAccepted, resultID)
}

func (route *Routes) ConnectSubToFinalCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w,r, &FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCatSubFinal(CatSubID, CatFinalID) VALUES(?,?)", FinalSub.CatStart, FinalSub.CatEnd)
	if err != nil{
		fmt.Println(err)
	}

	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	helpers.WriteJSON(w, http.StatusAccepted, resultID)
}

type CatToProd struct {
	Cat int `json:"Category"`
	Prod int `json:"Product"`
}
func (route *Routes) ConnectFinalToProdCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalProd := CatToProd{}
	err := helpers.ReadJSON(w,r, &FinalProd)
	if err != nil{
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCatFinalProd(CatFinalID, ProductID) VALUES(?,?)", FinalProd.Cat, FinalProd.Prod)

	if err != nil{
		fmt.Println(err)
	}

	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	helpers.WriteJSON(w, http.StatusAccepted, resultID)
}


type ReadCat struct{
	Category int `json:"category"`
}
func (route *Routes) InsertIntoFinalProd(w http.ResponseWriter, r *http.Request){
	ReadCatR := ReadCat{}
	err := helpers.ReadJSON(w,r,&ReadCatR)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("InsertIntoCategory ReadCatR",ReadCatR)
	FinalProd := "INSERT INTO tblCatFinalProd(CatFinalID, ProductID) VALUES(?,?)"
	route.DB.Exec(FinalProd, 1,ReadCatR.Category)
}

type CategoryReturn struct{
	CategoryName string `json:"CategoryName"`
	CategoryDescription string `json:"CategoryDescription"`
}

type CategoriesList struct{
	collection []CategoryReturn
}

func (route *Routes) ReturnAllPrimeCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesPrime"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)

}

func (route *Routes) ReturnAllSubCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesSub"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)
}

func (route *Routes) ReturnAllFinalCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesFinal"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)
}

// Product automatically creates Variation
type ProductCreate struct{
	Name string `json:"Product_Name"`
	Description string `json:"Product_Description"`
	Price string `json:"Product_Price"`
	VariationName string `json:"Variation_Name"`
	VariationDescription string `json:"Variation_Description"`
	VariationPrice float32 `json:"Variation_Price"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}

// Needs to get SKU, UPC, Primary Image to get created. PRimary Image needs to be a google/AWS bucket
func (route *Routes) CreateProduct(w http.ResponseWriter, r *http.Request){

	
}

type VariationCreate struct{
	Name string `json:"Variation_Name"`
	Description string `json:"Variation_Description"`
	Price string `json:"Variation_Price"`
	PrimaryImage string `json:"Primary_Image"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}

func (route *Routes) CreateVariation(w http.ResponseWriter, r *http.Request){

	
}