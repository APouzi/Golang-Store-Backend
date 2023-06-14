package routes

import (
	"database/sql"

	"github.com/Apouzi/golang-shop/app/api/authorization"
	adminendpoints "github.com/Apouzi/golang-shop/app/api/routes/admin"
	indexendpoints "github.com/Apouzi/golang-shop/app/api/routes/index"
	productendpoints "github.com/Apouzi/golang-shop/app/api/routes/product"
	userendpoints "github.com/Apouzi/golang-shop/app/api/routes/user"
	"github.com/go-chi/chi"
)



func RouteDigest(digest *chi.Mux, db *sql.DB) *chi.Mux{

	rIndex := indexendpoints.InstanceIndexRoutes(db)

	rProduct := productendpoints.InstanceProductsRoutes(db)

	rUser := userendpoints.InstanceUserRoutes(db)

	rAdmin := adminendpoints.InstanceAdminRoutes(db)

	
	digest.Group(func(digest chi.Router){
		digest.Use(authorization.ValidateToken)
		digest.Get("/users/profile",rUser.UserProfile)
	})

	//Index and Product
	digest.Get("/", rIndex.Index)
	digest.Post("/superusercreation",rUser.AdminSuperUserCreation)
	
	digest.Get("/product/{Product_ID}",rProduct.GetOneProductsEndPoint)
	digest.Get("/products/",rProduct.GetAllProductsEndPoint)
	// digest.Get("/products/{CategoryName}",r.GetProductCategoryEndPointFinal)
	digest.Post("/users/",rUser.Register)
	digest.Post("/users/login",rUser.Login)
	
	// digest.Get("/categories/",r.GetAllCategories)

	// These are testing for categories
	digest.Post("/products/test-categories", rProduct.CreateTestCategory)
	digest.Get("/products/test-categories/pullTest", rProduct.PullTestCategory)
	digest.Post("/products/test-categories/InsertTest", rAdmin.InsertIntoFinalProd)

	// Admin need to lockdown based on jwt payload and scope
	digest.Post("/products/", rAdmin.CreateProduct)
	digest.Post("/products/variation", rAdmin.CreateVariation)
	digest.Post("/products/inventory", rAdmin.CreateInventoryLocation)
	digest.Post("/category/prime", rAdmin.CreatePrimeCategory)
	digest.Post("/category/sub", rAdmin.CreateSubCategory)
	digest.Post("/category/final", rAdmin.CreateFinalCategory)
	digest.Post("/category/primetosub",rAdmin.ConnectPrimeToSubCategory)
	digest.Post("/category/subtofinal",rAdmin.ConnectSubToFinalCategory)
	digest.Post("/category/finaltoprod",rAdmin.ConnectFinalToProdCategory)
	digest.Get("/category/primes", rAdmin.ReturnAllPrimeCategories)
	digest.Get("/category/subs", rAdmin.ReturnAllSubCategories)
	digest.Get("/category/finals", rAdmin.ReturnAllFinalCategories)
	return digest
}