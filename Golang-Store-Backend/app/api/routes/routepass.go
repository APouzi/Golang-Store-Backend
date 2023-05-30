package routes

import (
	"database/sql"

	"github.com/Apouzi/golang-shop/app/api/authorization"
	"github.com/Apouzi/golang-shop/app/api/database"
	"github.com/go-chi/chi"
)

type Routes struct{
	DB *sql.DB
	ProductQuery *database.PrepareStatmentsProducts
	UserQuery *database.UserStatments
}

func RouteDigest(digest *chi.Mux, db *sql.DB) *chi.Mux{
	
	r := Routes{
		DB: db,
		ProductQuery: database.InitPrepare(db),
		UserQuery: database.InitUserStatments(db),
	}

	
	digest.Group(func(digest chi.Router){
		digest.Use(authorization.ValidateToken)
		digest.Post("/users/verify",r.VerifyTest)
		digest.Get("/users/profile",r.UserProfile)
	})



	
	//Index and Product
	digest.Get("/", r.Index)
	
	digest.Get("/products/{ProductID}",r.GetOneProductsEndPoint)
	digest.Get("/products/",r.GetAllProductsEndPoint)
	// digest.Get("/products/{CategoryName}",r.GetProductCategoryEndPointFinal)
	digest.Post("/users/",r.Register)
	digest.Post("/users/login",r.Login)
	
	// digest.Get("/categories/",r.GetAllCategories)

	// These are testing for categories
	digest.Post("/products/test-categories", r.CreateTestCategory)
	digest.Get("/products/test-categories/pullTest", r.PullTestCategory)
	digest.Post("/products/test-categories/InsertTest", r.InsertIntoFinalProd)

	return digest
}