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
		digest.Use(authorization.EnsureValidToken())
		digest.Post("/login", r.Login)
	})


	
	//Index and Product
	digest.Get("/", r.Index)
	digest.Get("/products/",r.GetAllProductsEndPoint)
	digest.Get("/products/{ProductID}",r.GetOneProductsEndPoint)
	digest.Get("/products/{CategoryName}",r.GetProductCategoryEndPointFinal)
	digest.Post("/users/",r.Register)
	// digest.Get("/categories/",r.GetAllCategories)

	return digest
}