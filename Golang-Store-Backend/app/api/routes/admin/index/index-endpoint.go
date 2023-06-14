package indexendpoints

import (
	"database/sql"
	"log"
	"net/http"

	// "./helpers"

	"github.com/Apouzi/golang-shop/app/api/helpers"
)

type ProductRoutes struct{
	DB *sql.DB
}

func InstanceIndexRoutes(db *sql.DB ) *ProductRoutes {
	r := &ProductRoutes{
		DB: db,
	}
	return r
}


func(route *ProductRoutes) Index(w http.ResponseWriter, r *http.Request) {
	
	db := route.DB
	_, err := db.Query("SELECT from tblTEST where id = 1",1)
	if err != nil{
		log.Println(err)
	}
	payload := helpers.ErrorJSONResponse{Error:false, Message: "All good"}



	helpers.WriteJSON(w, http.StatusAccepted, payload)
}