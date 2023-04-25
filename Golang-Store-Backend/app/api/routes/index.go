package routes

import (
	"fmt"
	"log"
	"net/http"

	// "./helpers"

	"github.com/Apouzi/golang-shop/app/api/authorization"
	"github.com/Apouzi/golang-shop/app/api/helpers"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

func(route *Routes) Index(w http.ResponseWriter, r *http.Request) {
	
	db := route.DB
	result, err := db.Query("SELECT from tblTEST where id = 1",1)
	if err != nil{
		log.Println(err)
	}
	fmt.Println("TEST:", result)
	payload := helpers.ErrorJSONResponse{Error:false, Message: "All good"}



	helpers.WriteJSON(w, http.StatusAccepted, payload)
}

func(route *Routes) Login(w http.ResponseWriter, r *http.Request){
	request := helpers.UserLoginRequest{}
	err := helpers.ReadJSON(w, r, &request)
	if err != nil{
		log.Panic("Sucka broke")
	}
	response := helpers.UserLoginResponse{UserID: request.Email}
	
	
	tokens := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	claims := tokens.CustomClaims.(*authorization.CustomClaims)
	if !claims.HasScope("something"){
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("message","Insufficent scope")
		helpers.WriteJSON(w,http.StatusBadRequest, "Oh nooo")
		return
	}
	
	err = helpers.WriteJSON(w,http.StatusAccepted, &response)
	if err != nil{
		log.Panic(err)
	}
	
	// helpers.WriteJSON(w, http.StatusAccepted,helpers.ErrorJSONResponse{Error: false, Message: "We read it"})
}