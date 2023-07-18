package routes

import (
	"database/sql"

	"github.com/Apouzi/golang-shop/app/api/authorization"
	adminendpoints "github.com/Apouzi/golang-shop/app/api/routes/admin"
	indexendpoints "github.com/Apouzi/golang-shop/app/api/routes/index"
	productendpoints "github.com/Apouzi/golang-shop/app/api/routes/product"
	testroutes "github.com/Apouzi/golang-shop/app/api/routes/test-routes"
	userendpoints "github.com/Apouzi/golang-shop/app/api/routes/user"
	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
)



func RouteDigest(digest *chi.Mux, db *sql.DB, redis *redis.Client) *chi.Mux{

	rIndex := indexendpoints.InstanceIndexRoutes(db)

	rProduct := productendpoints.InstanceProductsRoutes(db, redis)

	rUser := userendpoints.InstanceUserRoutes(db)

	rAdmin := adminendpoints.InstanceAdminRoutes(db)

	AuthMiddleWare := authorization.InjectDBRef(db)

	rTestRoutes := testroutes.InjectDBRef(db, redis)

	digest.Group(func(digest chi.Router){
		digest.Use(AuthMiddleWare.ValidateToken)
		digest.Get("/users/profile",rUser.UserProfile)
	})

	//Index and Product
	digest.Get("/", rIndex.Index)
	digest.Post("/superusercreation",rUser.AdminSuperUserCreation)
	
	digest.Get("/products/{ProductID}",rProduct.GetOneProductsEndPoint)
	digest.Get("/products/",rProduct.GetAllProductsEndPoint)
	// digest.Get("/products/{CategoryName}",r.GetProductCategoryEndPointFinal)
	digest.Post("/users/",rUser.Register)
	digest.Post("/users/login",rUser.Login)

	// Testing Routes
	digest.Get("/products-test-redis",rTestRoutes.GetOneProductRedis)
	digest.Get("/products-test-sql",rTestRoutes.GetOneProductSQL)
	digest.Get("/products/test-categories/pullTest", rTestRoutes.PullTestCategory)
	digest.Post("/products/test-categories", rTestRoutes.CreateTestCategory)
	
	// digest.Get("/categories/",r.GetAllCategories)
	
	digest.Post("/products/test-categories/InsertTest", rAdmin.InsertIntoFinalProd)

	// Admin need to lockdown based on jwt payload and scope
	digest.Group(func(digest chi.Router){
		digest.Use(AuthMiddleWare.ValidateToken)
		digest.Use(AuthMiddleWare.HasAdminScope)
		digest.Post("/products/", rAdmin.CreateProduct)
	})
	
	digest.Post("/products/{ProductID}/variation", rAdmin.CreateVariation)
	digest.Post("/products/inventory", rAdmin.CreateInventoryLocation)
	digest.Post("/category/prime", rAdmin.CreatePrimeCategory)
	digest.Post("/category/sub", rAdmin.CreateSubCategory)
	digest.Post("/category/final", rAdmin.CreateFinalCategory)
	digest.Delete("/category/prime/{CatPrimeName}",rAdmin.DeletePrimeCategory)
	digest.Delete("/category/sub/{CatSubName}",rAdmin.DeleteSubCategory)
	digest.Delete("/category/final/{CatFinalName}",rAdmin.DeleteFinalCategory)
	digest.Post("/category/primetosub",rAdmin.ConnectPrimeToSubCategory)
	digest.Post("/category/subtofinal",rAdmin.ConnectSubToFinalCategory)
	digest.Post("/category/finaltoprod",rAdmin.ConnectFinalToProdCategory)
	digest.Get("/category/primes", rAdmin.ReturnAllPrimeCategories)
	digest.Get("/category/subs", rAdmin.ReturnAllSubCategories)
	digest.Get("/category/finals", rAdmin.ReturnAllFinalCategories)
	digest.Patch("/products/{ProductID}",rAdmin.EditProduct)
	digest.Patch("/products/{ProductID}/variation/{VariationID}",rAdmin.EditVariation)
	digest.Post("/admin/{UserID}", rAdmin.UserToAdmin)

	digest.Get("/tables",rAdmin.GetAllTables)
	return digest
}