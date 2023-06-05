package routes

import (
	"database/sql"
	"fmt"
	"log"
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


func (route *Routes) GetOneProductsEndPoint(w http.ResponseWriter, r *http.Request){
	query :=  chi.URLParam(r,"Product_ID")
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

func (route *Routes) PullTestCategory(w http.ResponseWriter, r *http.Request){
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
	result, err := route.DB.Exec("INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)", FinalProd.Cat, FinalProd.Prod)

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
	FinalProd := "INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)"
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
	Price float32 `json:"Product_Price"`
	VariationName string `json:"Variation_Name"`
	VariationDescription string `json:"Variation_Description"`
	VariationPrice float32 `json:"Variation_Price"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}

// Needs to get SKU, UPC, Primary Image to get created. PRimary Image needs to be a google/AWS bucket
func(route *Routes) CreateProduct(w http.ResponseWriter, r *http.Request){
	transaction, err := route.DB.Begin()
	if err != nil{
		log.Println("Error creating a transation in CreateProduct")
		log.Println(err)
	}

	productRetrieve := &ProductCreate{}

	helpers.ReadJSON(w, r, &productRetrieve)

	tRes, err := transaction.Exec("INSERT INTO tblProducts(Product_Name, Product_Description, Product_Price) VALUES(?,?,?)", productRetrieve.Name,productRetrieve.Description,productRetrieve.Price)
	if err != nil{
		fmt.Println("transaction at tblProduct has failed")
		fmt.Println(err)
	}
	prodID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProduct has failed")
		fmt.Println(err)
	}
	tRes, err = transaction.Exec("INSERT INTO tblProductVariation(Product_ID,Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)",prodID, productRetrieve.VariationName, productRetrieve.VariationDescription, productRetrieve.VariationPrice)
	if err != nil{
		fmt.Println("transaction at tblProductVariation has failed")
		fmt.Println(err)
	}
	ProdVarID, err :=  tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProductVariation has failed")
		fmt.Println(err)
	}
	tRes, err = transaction.Exec("INSERT INTO tblProductInventory(Variation_ID, Quantity) VALUES(?,?)",  ProdVarID,productRetrieve.VariationQuantity)
	if err != nil {
		fmt.Println("transaction at tblProductInventory has failed")
		fmt.Println(err)
	}
	ProdInvID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProductInventory has failed")
		fmt.Println(err)
	}

	_, err = transaction.Exec("INSERT INTO tblLocation(Inv_ID,Location_AT) VALUES(?,?)", ProdInvID,productRetrieve.LocationAt)
	if err != nil {
		fmt.Println("transaction at tblProductInventory has failed")
		fmt.Println(err)
	}	
	transaction.Commit()
}


type VariationCreate struct{
	ProductID int64 `json:"Product_ID"`
	Name string `json:"Variation_Name"`
	Description string `json:"Variation_Description"`
	Price float32 `json:"Variation_Price"`
	PrimaryImage string `json:"Primary_Image,omitempty"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}

type ProdExist struct{
	ProductExists bool `json:"Product_Exists"`
	Message string `json:"Message"`
}

type variCrtd struct{
	VariationID int64 `json:"Product_ID"`
	LocationExists bool `json:"Location_Exists"`
}

func (route *Routes) CreateVariation(w http.ResponseWriter, r *http.Request){

	variation := VariationCreate{}
	helpers.ReadJSON(w,r, &variation)
// Check if product exists, if not, then return false
	row := route.DB.QueryRow("SELECT Product_ID FROM tblProducts WHERE Product_ID = ?",variation.ProductID)
	if row.Err() != nil{
		fmt.Println(row.Err().Error())
		return
	}
	// DELETE

	var exists bool
	if row.Scan(&exists); exists == false {
		
		msg := ProdExist{}
		msg.ProductExists = false
		msg.Message = "Product provided does not exist"
		helpers.WriteJSON(w,http.StatusAccepted,msg)
		log.Println("Variation Creation completed")
		return
	}
	// Implement the returns for this to allow for proper exiting 
	var varit sql.Result
	var err error
	if variation.PrimaryImage != "" {
		varit, err = route.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", variation.ProductID,variation.Name, variation.Description, variation.Price)
		if err != nil{
			fmt.Println("insert into tblProductVariation failed")
			fmt.Println(err)
		}

		helpers.WriteJSON(w, http.StatusCreated,)
	}
	varit, err = route.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price, PRIMARY_IMAGE) VALUES(?,?,?,?,?)", variation.ProductID,variation.Name, variation.Description, variation.Price, variation.PrimaryImage)
	if err != nil{
		fmt.Println("insert into tblProductVariation failed")
		fmt.Println(err)
	}
//Check if location exists, if not, then we should prompt them to create one
	varitID, err := varit.LastInsertId()
	if err != nil{
		fmt.Println("issue with Variation_ID failed")
		fmt.Println(err)
	}
	if variation.LocationAt == ""{
		msg := variCrtd{}
		msg.LocationExists = false
		msg.VariationID = varitID
		helpers.WriteJSON(w, http.StatusAccepted, msg)
		return
	}

	 
	
}

type ProdInvLocCreation struct{
	VarID int64 `json:"Variation_ID"`
	Quantity int `json:"Quantity"`
	Location string `json:"Location"`
}
type PILCreated struct{
	InvID int64 `json:"Inv_ID"`
	Quantity int `json:"Quantity"`
	Location string `json:"Location"`
}

func(route *Routes) CreateInventoryLocation(w http.ResponseWriter, r *http.Request){
	// Test for Variantion existness
	pil := ProdInvLocCreation{}
	helpers.ReadJSON(w,r,&pil)
	row := route.DB.QueryRow("SELECT Product_ID FROM tblProductVariation WHERE Product_ID = ?",pil.VarID)
	if row.Err() != nil{
		fmt.Println(row.Err().Error())
		return
	}
	var exists bool
	if row.Scan(&exists); exists == false {
		
		msg := ProdExist{}
		msg.ProductExists = false
		msg.Message = "Variation record provided does not exist"
		helpers.WriteJSON(w,http.StatusAccepted,msg)
		log.Println("Location Creation failed")
		return
	}
	res ,err:= route.DB.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_At) VALUES(?,?,?)")
	
	if err != nil{
		fmt.Println("failed to create tblProductInventoryLocation")
		fmt.Println(err)
		helpers.ErrorJSON(w,err,http.StatusForbidden)
		return
	}

	pilID, err := res.LastInsertId()
	
	if err != nil{
		fmt.Println("result of tblProductInventoryLocation failed")
	}
	pilReturn := PILCreated{}
	pilReturn.InvID = pilID
	pilReturn.Quantity = pil.Quantity
	pilReturn.Location = pil.Location
	helpers.WriteJSON(w, http.StatusAccepted, pil)
}