package database

import (
	"database/sql"
	"fmt"
	"log"
)

type PrepareStatmentsProducts struct{
	GetAllProdStmt *sql.Stmt
	GetOneProductStmt *sql.Stmt
	GetAllProductByCategoryStmt *sql.Stmt
	GetAllProductByCategoryPrimeStmt *sql.Stmt
	GetAllProductByCategorySubStmt *sql.Stmt
	GetAllProductByCategoryFinalStmt *sql.Stmt
	GetProductPrimeCategorybyID *sql.Stmt
}

func InitPrepare(db *sql.DB) *PrepareStatmentsProducts{
	prep := &PrepareStatmentsProducts{}
	var err error
	prep.GetAllProdStmt, err = db.Prepare("SELECT Product_ID, Product_Name, Product_Description, PRIMARY_IMAGE FROM tblProducts")
	if err != nil{
		log.Fatal(err)
	}

	prep.GetOneProductStmt, err = db.Prepare("SELECT Product_ID, Product_Name, Product_Description, PRIMARY_IMAGE FROM tblProducts where Product_ID = ?")
	if err != nil{
		log.Fatal(err)
	}

	// prep.GetProductPrimeCategoryByID, err = db.Prepare("SELECT tblProducts.Product_ID, tblProducts.Product_Name FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.Product_ID = tblProducts.Product_ID JOIN tblCategoriesFinal ON tblCategoriesFinal.Category_ID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.Category_ID JOIN tblCategoriesSub ON tblCategoriesSub.Category_ID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.Category_ID JOIN tblCategoriesPrime ON tblCategoriesPrime.Category_ID = tblCatPrimeSub.CatPrimeID where tblCategoriesPrime.Category_ID")

	// prep.GetAllProductByCategoryStmt, err = db.Prepare("SELECT * FROM tblProducts JOIN tblCategory ON tblProducts.Category_ID = tblCategories.id WHERE tblCategori.id = ?")
	// if err != nil{
	// 	log.Fatal(err)
	// }

	// prep.GetAllProductByCategoryPrimeStmt, err = db.Prepare("SELECT tblProducts.Product_ID, tblProducts.Product_Name, tblProducts.Product_Description, tblProducts.Product_Price FROM tblProducts JOIN tblProductsCategoriesPrime ON tblProducts.Product_ID = tblProductsCategoriesPrime.Product_ID JOIN tblCategoriesPrime ON tblProductsCategoriesPrime.Category_ID = tblCategoriesPrime.Category_ID WHERE tblCategoriesPrime.CategoryName = ?") 
	// if err != nil{
	// 	log.Fatal(err)
	// }

	// prep.GetAllProductByCategorySubStmt, err = db.Prepare("SELECT tblProducts.Product_ID, tblProducts.Product_Name, tblProducts.Product_Description, tblProducts.Product_Price FROM tblProducts JOIN tblProductsCategoriesSub ON tblProducts.Product_ID = tblProductsCategoriesSub.Product_ID JOIN tblCategoriesSub ON tblProductsCategoriesSub.Category_ID = tblCategoriesSub.Category_ID WHERE tblCategoriesSub.CategoryName = ?") 
	// if err != nil{
	// 	log.Fatal(err)
	// }

	// prep.GetAllProductByCategoryFinalStmt, err = db.Prepare("SELECT tblProducts.Product_ID, tblProducts.Product_Name, tblProducts.Product_Description, tblProducts.Product_Price FROM tblProducts JOIN tblProductsCategoriesFinal ON tblProducts.Product_ID = tblProductsCategoriesFinal.Product_ID JOIN tblCategoriesFinal ON tblProductsCategoriesFinal.Category_ID = tblCategoriesFinal.Category_ID WHERE tblCategoriesFinal.CategoryName = ?") 
	if err != nil{
		log.Fatal(err)
	}



	return prep
}

type ProductList struct {
	Products []Product `json:"products"`
}

func(prep *PrepareStatmentsProducts) GetAllProducts(db *sql.DB) []ProductJSON {
	rows,err:=prep.GetAllProdStmt.Query()
	if err != nil{
		fmt.Println(err)
	}
	
	products := []ProductJSON{}
	
	defer rows.Close()
	for rows.Next(){
		prodJSON := ProductJSON{}
		err := rows.Scan(
			&prodJSON.Product_ID, 
			&prodJSON.Product_Name, 
			&prodJSON.Product_Description,  
			&prodJSON.PRIMARY_IMAGE,
		)
		if err != nil{
			fmt.Println("scanning error:",err)
		}
		products = append(products, prodJSON)
	}

	return products
}

func(prep *PrepareStatmentsProducts) GetOneProduct(db *sql.DB, id int) ProductJSON {
	rows :=prep.GetOneProductStmt.QueryRow(id)
	
	prodJSON := ProductJSON{}
	
	err := rows.Scan(
		&prodJSON.Product_ID, 
		&prodJSON.Product_Name, 
		&prodJSON.Product_Description, 
		&prodJSON.Product_Price, 
		&prodJSON.SKU, 
		&prodJSON.UPC, 
		&prodJSON.PRIMARY_IMAGE,
	)
	if err != nil{
		fmt.Println("scanning error:",err)
	}
	

	return prodJSON
}

func(prep *PrepareStatmentsProducts) GetCategoryProduct(db *sql.DB, category int, catMap map[string]int) []ProductJSON{
	rows, err :=prep.GetOneProductStmt.Query(category)
	if err!=nil{
		fmt.Println("Error in GetCategoryProduct")
	}
	prodJSON := ProductJSON{}
	products := []ProductJSON{}
	for rows.Next(){
		err := rows.Scan(
			&prodJSON.Product_ID, 
			&prodJSON.Product_Name, 
			&prodJSON.Product_Description, 
			// &prodJSON.PRIMARY_IMAGE,
		)
		if err != nil{
			fmt.Println("scanning error:",err)
		}
		products = append(products, prodJSON)
	}
	
	return products
}

func(prep *PrepareStatmentsProducts) GetProductCategoryFinal(db *sql.DB, category string) []ProductJSON{
	rows, err :=prep.GetAllProductByCategoryStmt.Query(category)
	if err!=nil{
		fmt.Println("Error in GetCategoryProduct")
	}
	prodJSON := ProductJSON{}
	products := []ProductJSON{}
	for rows.Next(){
		err := rows.Scan(
			&prodJSON.Product_ID, 
			&prodJSON.Product_Name, 
			&prodJSON.Product_Description, 
			// &prodJSON.PRIMARY_IMAGE,
		)
		if err != nil{
			fmt.Println("scanning error:",err)
		}
		products = append(products, prodJSON)
	}
	
	return products
}