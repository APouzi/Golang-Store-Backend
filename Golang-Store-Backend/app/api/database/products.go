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
}

func InitPrepare(db *sql.DB) *PrepareStatmentsProducts{
	prep := &PrepareStatmentsProducts{}
	var err error
	prep.GetAllProdStmt, err = db.Prepare("SELECT ProductID, ProductName, ProductDescription, ProductPrice, SKU, UPC, PRIMARY_IMAGE FROM tblProducts")
	if err != nil{
		log.Fatal(err)
	}

	prep.GetOneProductStmt, err = db.Prepare("SELECT ProductID, ProductName, ProductDescription, ProductPrice, SKU, UPC, PRIMARY_IMAGE FROM tblProducts where ProductID = ?")
	if err != nil{
		log.Fatal(err)
	}

	// prep.GetAllProductByCategoryStmt, err = db.Prepare("SELECT * FROM tblProducts JOIN tblCategory ON tblProducts.CategoryID = tblCategories.id WHERE tblCategori.id = ?")
	// if err != nil{
	// 	log.Fatal(err)
	// }

	prep.GetAllProductByCategoryPrimeStmt, err = db.Prepare("SELECT tblProducts.ProductID, tblProducts.ProductName, tblProducts.ProductDescription, tblProducts.ProductPrice FROM tblProducts JOIN tblProductsCategoriesPrime ON tblProducts.ProductID = tblProductsCategoriesPrime.ProductID JOIN tblCategoriesPrime ON tblProductsCategoriesPrime.CategoryID = tblCategoriesPrime.CategoryID WHERE tblCategoriesPrime.CategoryName = ?") 
	if err != nil{
		log.Fatal(err)
	}

	prep.GetAllProductByCategorySubStmt, err = db.Prepare("SELECT tblProducts.ProductID, tblProducts.ProductName, tblProducts.ProductDescription, tblProducts.ProductPrice FROM tblProducts JOIN tblProductsCategoriesSub ON tblProducts.ProductID = tblProductsCategoriesSub.ProductID JOIN tblCategoriesSub ON tblProductsCategoriesSub.CategoryID = tblCategoriesSub.CategoryID WHERE tblCategoriesSub.CategoryName = ?") 
	if err != nil{
		log.Fatal(err)
	}

	prep.GetAllProductByCategoryFinalStmt, err = db.Prepare("SELECT tblProducts.ProductID, tblProducts.ProductName, tblProducts.ProductDescription, tblProducts.ProductPrice FROM tblProducts JOIN tblProductsCategoriesFinal ON tblProducts.ProductID = tblProductsCategoriesFinal.ProductID JOIN tblCategoriesFinal ON tblProductsCategoriesFinal.CategoryID = tblCategoriesFinal.CategoryID WHERE tblCategoriesFinal.CategoryName = ?") 
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
			&prodJSON.ProductID, 
			&prodJSON.ProductName, 
			&prodJSON.ProductDescription, 
			&prodJSON.ProductPrice, 
			&prodJSON.SKU, 
			&prodJSON.UPC, 
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
		&prodJSON.ProductID, 
		&prodJSON.ProductName, 
		&prodJSON.ProductDescription, 
		&prodJSON.ProductPrice, 
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
			&prodJSON.ProductID, 
			&prodJSON.ProductName, 
			&prodJSON.ProductDescription, 
			&prodJSON.ProductPrice, 
			&prodJSON.SKU, 
			&prodJSON.UPC, 
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
			&prodJSON.ProductID, 
			&prodJSON.ProductName, 
			&prodJSON.ProductDescription, 
			&prodJSON.ProductPrice, 
			// &prodJSON.SKU, 
			// &prodJSON.UPC, 
			// &prodJSON.PRIMARY_IMAGE,
		)
		if err != nil{
			fmt.Println("scanning error:",err)
		}
		products = append(products, prodJSON)
	}
	
	return products
}