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
	Product_ID          int
	Product_Name        string
	Product_Description string
	Product_Price       float32
	SKU                 string
	UPC                 string
	PRIMARY_IMAGE       string
	ProductDateAdded    string
	ModifiedDate        string
}

type Category struct {
	Name string
}

type Inventory struct {
	Quantity int
}
