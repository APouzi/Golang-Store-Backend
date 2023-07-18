package adminendpoints

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Apouzi/golang-shop/app/api/helpers"
	"github.com/go-chi/chi"
)

type AdminRoutes struct{
	DB *sql.DB
}

func InstanceAdminRoutes(db *sql.DB,  ) *AdminRoutes {
	r := &AdminRoutes{
		DB: db,
	}
	return r
}

type CategoryInsert struct{
	CategoryName string `json:"CategoryName"`
	CategoryDescription string `json:"CategoryDescription"`
}

func(route *AdminRoutes)  AdminTableScopeCheck(adminTable string, tableName string ,adminID any, w http.ResponseWriter) bool{
	// strQ := "SELECT AdminID FROM" + adminTable + "WHERE Tablename = " + string(adminID) + " AND AdminID = " + adminID
	var exists bool
	var strBuild strings.Builder
	strBuild.WriteString("SELECT AdminID FROM ")
	strBuild.WriteString(adminTable)
	strBuild.WriteString(" WHERE TableName = ? AND AdminID = ?")
	route.DB.QueryRow(strBuild.String(), tableName, adminID).Scan(&exists)
	
	if exists == false{
		fmt.Println("Failed Query AdminTableScopeCheck endpoint")
		return false
	}

	return true
}

// Product automatically creates Variation
type ProductCreate struct{
	Name string `json:"Product_Name"`
	Description string `json:"Product_Description"`
	Price float64 `json:"Product_Price"`
	VariationName string `json:"Variation_Name"`
	VariationDescription string `json:"Variation_Description"`
	VariationPrice float32 `json:"Variation_Price"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}
type ProductCreateRetrieve struct{
	ProductID int64 `json:"Product_ID"`
	VarID int64 `json:"Variation_ID"`
	ProdInvLoc int64 `json:"Inv_ID,omitempty"`

}

// Needs to get SKU, UPC, Primary Image to get created. Primary Image needs to be a google/AWS bucket
func(route *AdminRoutes) CreateProduct(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("userid")
	if !route.AdminTableScopeCheck("tblCreateTables","tblProducts",userID, w){
		err := errors.New("Failed Query")
		helpers.ErrorJSON(w, err, 400)
		return
	}
	transaction, err := route.DB.Begin()
	if err != nil{
		log.Println("Error creating a transation in CreateProduct")
		log.Println(err)
	}

	productRetrieve := &ProductCreate{}

	helpers.ReadJSON(w, r, &productRetrieve)

	tRes, err := transaction.Exec("INSERT INTO tblProducts(Product_Name, Product_Description) VALUES(?,?)", productRetrieve.Name,productRetrieve.Description)
	if err != nil{
		fmt.Println("transaction at tblProduct has failed")
		fmt.Println(err)
		transaction.Rollback()
	}
	prodID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProduct has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}
	tRes, err = transaction.Exec("INSERT INTO tblProductVariation(Product_ID,Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)",prodID, productRetrieve.VariationName, productRetrieve.VariationDescription, productRetrieve.VariationPrice)
	if err != nil{
		fmt.Println("transaction at tblProductVariation has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}
	
	ProdVarID, err :=  tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProductVariation has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}
	PCR := ProductCreateRetrieve{
		ProductID: prodID,
		VarID: ProdVarID,
	}
	if productRetrieve.LocationAt == ""{
		
		err = transaction.Commit()
		if err != nil{
			fmt.Println(err)
			transaction.Rollback()
			return
		}
		helpers.WriteJSON(w,http.StatusAccepted,&PCR)
		return
	}

	tRes, err = transaction.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_AT) VALUES(?,?,?)",  ProdVarID,productRetrieve.VariationQuantity, productRetrieve.LocationAt)
	if err != nil {
		fmt.Println("transaction at tblProductInventory has failed")
		fmt.Println(err)
	}
	invID, err := tRes.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	PCR.ProdInvLoc = invID
	err = transaction.Commit()
	if err != nil{
		fmt.Println(err)
	}
	helpers.WriteJSON(w,http.StatusAccepted,&PCR)
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

func (route *AdminRoutes) CreateVariation(w http.ResponseWriter, r *http.Request){
	ProductID := chi.URLParam(r, "ProductID")
	variation := VariationCreate{}
	helpers.ReadJSON(w,r, &variation)
// Check if product exists, if not, then return false
	row := route.DB.QueryRow("SELECT Product_ID FROM tblProducts WHERE Product_ID = ?",ProductID)
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
		varitCrt := variCrtd{}
		varit, err = route.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", ProductID,variation.Name, variation.Description, variation.Price)
		if err != nil{
			fmt.Println("insert into tblProductVariation failed")
			fmt.Println(err)
		}
		varitCrt.VariationID, err = varit.LastInsertId()
		if err != nil{
			fmt.Println(err)
		}
		// helpers.WriteJSON(w, http.StatusCreated,varitCrt)
	}
	varit, err = route.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price, PRIMARY_IMAGE) VALUES(?,?,?,?,?)", ProductID,variation.Name, variation.Description, variation.Price, variation.PrimaryImage)
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

func(route *AdminRoutes) CreateInventoryLocation(w http.ResponseWriter, r *http.Request){
	// Test for Variantion existness
	pil := ProdInvLocCreation{}
	helpers.ReadJSON(w,r,&pil)
	row := route.DB.QueryRow("SELECT Variation_ID FROM tblProductVariation WHERE Variation_ID = ?",pil.VarID)
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
	res ,err:= route.DB.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_At) VALUES(?,?,?)", pil.VarID,pil.Quantity,pil.Location)
	
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



func (route *AdminRoutes) CreatePrimeCategory(w http.ResponseWriter, r *http.Request){
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

func (route *AdminRoutes) CreateSubCategory(w http.ResponseWriter, r *http.Request){
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

func (route *AdminRoutes) CreateFinalCategory(w http.ResponseWriter, r *http.Request){
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


func (route *AdminRoutes) ConnectPrimeToSubCategory(w http.ResponseWriter, r *http.Request){
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

func (route *AdminRoutes) ConnectSubToFinalCategory(w http.ResponseWriter, r *http.Request){
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
func (route *AdminRoutes) ConnectFinalToProdCategory(w http.ResponseWriter, r *http.Request){
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
func (route *AdminRoutes) InsertIntoFinalProd(w http.ResponseWriter, r *http.Request){
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

func (route *AdminRoutes) ReturnAllPrimeCategories(w http.ResponseWriter, r *http.Request){
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

func (route *AdminRoutes) ReturnAllSubCategories(w http.ResponseWriter, r *http.Request){
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

func (route *AdminRoutes) ReturnAllFinalCategories(w http.ResponseWriter, r *http.Request){
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


type ProductEdit struct{
	Name string `json:"Product_Name"`
	Description string `json:"Product_Description"`
}


func (route *AdminRoutes) EditProduct(w http.ResponseWriter, r *http.Request){
	ProdID := chi.URLParam(r, "ProductID")
	prodEdit := ProductEdit{}
	helpers.ReadJSON(w,r, &prodEdit)
	var buf strings.Builder
	buf.WriteString("UPDATE tblProducts SET")
	var count int = 0
	Varib := []any{}
	if prodEdit.Name != "" {
		if count == 0{
			buf.WriteString(" Product_Name = ?")
			Varib = append(Varib, prodEdit.Name)
			count++
		}
		buf.WriteString(", Product_Name = ?")
		Varib = append(Varib, prodEdit.Name)
	}
	if prodEdit.Description != "" {
		if count == 0{
			buf.WriteString(" Product_Description = ?")
			Varib = append(Varib, prodEdit.Description)
			count++
		}
		buf.WriteString(", Product_Description = ?")
		Varib = append(Varib, prodEdit.Description)
	}
	if count  == 0 {
		helpers.WriteJSON(w,http.StatusAccepted,"failed")
		return
	}

	buf.WriteString(", Modified_Date = ? WHERE Product_ID = ?")
	Varib = append(Varib, time.Now(),ProdID)
	_, err := route.DB.Exec(buf.String(), Varib...)
	if err != nil{
		fmt.Println("err with exec Edit Product Update")
		fmt.Println(err)
	}

	helpers.WriteJSON(w,http.StatusAccepted,&prodEdit)
	
}

type VariationEdit struct{
	VariationID int64 `json:"Variation_ID"`
	VariationProductID int64 `json:"Product_ID"`
	VariationName string `json:"Variation_Name"`
	VariationDescription string `json:"Variation_Description"`
	VariationPrice float32 `json:"Variation_Price"`
	SKU string `json:"SKU"`
	UPC string `json:"UPC"`
	PrimaryImage string `json:"Primary_Image,omitempty"`
	VariationQuantity int  `json:"Variation_Quantity"`
	LocationAt string `json:"Location_At"`
}

func (route *AdminRoutes) EditVariation(w http.ResponseWriter, r *http.Request){
	r.Header.Get("Authorization")
	VarID := chi.URLParam(r, "VariationID")
	VaritEdit := VariationEdit{}
	helpers.ReadJSON(w,r, &VaritEdit)
	var buf strings.Builder
	Varib := []any{}
	buf.WriteString("UPDATE tblProductVariation SET")
	var count int = 0
	if VaritEdit.VariationName != "" {
		if count == 0{
			buf.WriteString(" Variation_Name = ?")
			Varib = append(Varib, VaritEdit.VariationName)
			count++
		}
		buf.WriteString(", Variation_Name = ?")
		Varib = append(Varib, VaritEdit.VariationName)
	}
	if VaritEdit.VariationDescription != ""{
		if count == 0{
			buf.WriteString(" Variation_Description = ?")
			Varib = append(Varib, VaritEdit.VariationDescription)
			count++
		}
		buf.WriteString(", Variation_Description = ?")
		Varib = append(Varib, VaritEdit.VariationDescription)
	}
	if VaritEdit.SKU != ""{
		if count == 0 {
			buf.WriteString(" SKU = ?")
			Varib = append(Varib, VaritEdit.SKU)
			count++
		}
		buf.WriteString(", SKU = ?")
		Varib = append(Varib, VaritEdit.SKU)
	}
	if VaritEdit.UPC != ""{
		if count == 0{
			buf.WriteString(" UPC = ?")
			Varib = append(Varib, VaritEdit.UPC)
			count++
		}
		buf.WriteString(", UPC = ?")
		Varib = append(Varib, VaritEdit.UPC)
	}
	if VaritEdit.VariationPrice != 0 {
		if count == 0{
			buf.WriteString(" Variation_Price = ?")
			Varib = append(Varib, VaritEdit.VariationPrice)
			count++
		}
		buf.WriteString(", Variation_Price = ?")
		Varib = append(Varib, VaritEdit.VariationPrice)
	}
	buf.WriteString(" WHERE Variation_ID = ?")
	Varib = append(Varib, VarID)
	_,err := route.DB.Exec(buf.String(),Varib...)
	if err != nil{
		fmt.Println(err)
	}
	helpers.WriteJSON(w, http.StatusAccepted, VaritEdit)
}

type DeletedSendBack struct{
	SendBack bool `json:"Deleted"`
}
type AddedSendBack struct{
	IDSendBack int64 `json:"AddedID"`
}
func (route *AdminRoutes) DeletePrimeCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatPrimeName")
	if CatName == ""{
		fmt.Println("No CatPrimeName wasn't pulled")
		return
	}
	_, err := route.DB.Exec("DELETE FROM tblCategoriesPrime WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatPrimeName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}

	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}


func (route *AdminRoutes) DeleteSubCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatSubName")
	if CatName == ""{
		fmt.Println("No CatSubName wasn't pulled")
		return
	}
	
	_, err := route.DB.Exec("DELETE FROM tblCategoriesSub WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatSubName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}
	
	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}


func (route *AdminRoutes) DeleteFinalCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatFinalName")
	if CatName == ""{
		fmt.Println("No CatPrimeName wasn't pulled")
		return
	}
	
	_, err := route.DB.Exec("DELETE FROM tblCategoriesFinal WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatPrimeName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}

	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}
type Attribute struct{
	Attribute string `json:"attribute"`
}
func (route *AdminRoutes) AddAttribute(w http.ResponseWriter, r *http.Request){
	VarID := chi.URLParam(r,"VariationID")
	if VarID == ""{
		helpers.ErrorJSON(w, errors.New("please input VariationID"),400)
		return
	}
	att := Attribute{}
	
	err := helpers.ReadJSON(w,r,&att)
	if err != nil{
		helpers.ErrorJSON(w, err, 500)
		return
	}
	sql, err := route.DB.Exec("INSERT INTO tblProductAttribute (Variation_ID, AttributeName) VALUES(?,?)",VarID,att.Attribute)
	if err != nil{
		helpers.ErrorJSON(w,err, 400)
		return
	}
	var id int64
	id, err = sql.LastInsertId()
	if err != nil{
		helpers.ErrorJSON(w,errors.New("failed attribute LastInsertID"))
		return
	}
	sendBack := AddedSendBack{IDSendBack: id}
	helpers.WriteJSON(w, 200, sendBack)
}




func (route *AdminRoutes) GetAllTables(w http.ResponseWriter, r *http.Request){
	sql,err := route.DB.Query("show tables")
	if err != nil{
		fmt.Println("failed to get all tables")
		return
	}
	var table string
	list := []string{}
	for sql.Next(){
		sql.Scan(&table)
		list = append(list, table)
	}
	helpers.WriteJSON(w,200,list)
}


func(route *AdminRoutes) UserToAdmin(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r,"UserID")
	fmt.Println("UserToAdmin:",id)
	var exists bool
	route.DB.QueryRow("SELECT UserID FROM tblUser WHERE UserID = ?",id).Scan(&exists)
	if exists == false {
		helpers.ErrorJSON(w,errors.New("user doesn't exist") ,400)
		return
	}

	var UserID int64
	err := route.DB.QueryRow("SELECT UserID FROM tblUser WHERE UserID = ?", id).Scan(&UserID)
	if err != nil{
		helpers.ErrorJSON(w,errors.New("issue with scanning user into struct ") ,500)
		return
	}

	sql, err := route.DB.Exec("INSERT INTO tblAdminUsers (UserID, SuperUser) VALUES(?,?)",UserID,false)
	if err != nil{
		helpers.ErrorJSON(w,errors.New("failed insertinginto tblAdminUsers") ,500)
		return
	}
	type returnAdminID struct{
		UserID int64 `json:"AdminUserID"`
	}
	adminID, err := sql.LastInsertId()
	if err != nil{
		helpers.ErrorJSON(w,errors.New("couldn't retrieve id from LastInsertId") ,500)
		return
	}
	rAID := returnAdminID{UserID:adminID}
	helpers.WriteJSON(w,200,rAID)
}

