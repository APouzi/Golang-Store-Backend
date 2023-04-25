package database

type Models struct {
	UserCustomer UserCustomer
	UserAdmin    UserAdmin
	Customer     Customer
	Admin        Admin
	Product      Product
	Category     Category
	Inventory    Inventory
}

type UserCustomer struct {
	Email      string
	First_Name string
	Last_Name  string
	Customer   *Customer
}

type UserAdmin struct {
	Email      string
	First_Name string
	Last_Name  string
	Admin      *Admin
}

type Customer struct {
	Street_Address string
	Phone_Number   string
	State          string
}

type Profile struct {
}

type Admin struct {
	Privlages []string
	SuperUser bool
}

// --------- Product ---------

type Product struct {
	ProductID          int
	ProductName        string
	ProductDescription string
	ProductPrice       float32
	SKU                string
	UPC                string
	PRIMARY_IMAGE      string
	ProductDateAdded   string
	ModifiedDate       string
}

type ProductJSON struct {
	ProductID          int     `json:"ProductID"`
	ProductName        string  `json:"ProductName"`
	ProductDescription string  `json:"ProductDescription"`
	ProductPrice       float32 `json:"ProductPrice"`
	SKU                string  `json:"SKU"`
	UPC                string  `json:"UPC"`
	PRIMARY_IMAGE      string  `json:"PRIMARY_IMAGE,omitempty"`
	// ProductDateAdded   string  `json:"DateAdded"`
	// ModifiedDate       string `json:"ModifiedDate"`
}

type Category struct {
	Name string
}

type Inventory struct {
	Quantity int
}
