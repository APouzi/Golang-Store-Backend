package authorization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Apouzi/golang-shop/app/api/helpers"
	"github.com/golang-jwt/jwt/v5"
)


type JWTtest struct{
	Token string `json:"JWT"`
}

type AuthMiddleWareStruct struct{
	db *sql.DB
}

func InjectDBRef(db *sql.DB) *AuthMiddleWareStruct{
	AMWS := AuthMiddleWareStruct{}
	AMWS.db = db
	return &AMWS
}

func(db *AuthMiddleWareStruct) ValidateToken(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		jwttoken := r.Header.Get("Authorization")
		// When the jwttoken comes in, it will input "bearer" into the token and we have to remove this from the token so we can parse it. 
		jwttoken = strings.Split(jwttoken, "Bearer ")[1]
		token, err := jwt.Parse(jwttoken, func(token *jwt.Token) (interface{}, error) {
			return []byte("Testing key"), nil
		})
		if err != nil{
			fmt.Println("ValidateToken Failed")
			helpers.ErrorJSON(w,err)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		ctx := context.WithValue(r.Context(), "userid", claims["userId"])
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}


// Start of checking if given user is a SuperUser
func(db *AuthMiddleWareStruct) HasSuperUserScope(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		jwttoken := r.Header.Get("Authorization")
		if jwttoken == ""{
			fmt.Println("No Authorization")
			return
		}
		jwttoken = strings.Split(jwttoken, "Bearer ")[1]
		token, err := jwt.Parse(jwttoken, func(token *jwt.Token) (interface{}, error) {
			return []byte("Testing key"), nil
		})
		if err != nil{
			fmt.Println("HasSuperUserScope failed")
			fmt.Println(err)
			helpers.ErrorJSON(w,err, 400)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["admin"] != "True"{
			err := errors.New("failed superUser check")
			helpers.ErrorJSON(w,err,400)
			return
		}

		ctx := context.WithValue(r.Context(), "userid", claims["userId"])
		var exists bool
		db.db.QueryRow("SELECT UserID FROM tblAdminUsers WHERE UserID = ? AND SuperUser = 1", claims["userId"]).Scan(&exists)
		if exists == false{
			fmt.Println("User not in Admin, HasAdminScope has failed")
			err := errors.New("failed admin check")
			helpers.ErrorJSON(w,err, 400)
			return
		}
		next.ServeHTTP(w,r.WithContext(ctx))

	})
}




func(db *AuthMiddleWareStruct) HasAdminScope(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwttoken := r.Header.Get("Authorization")
		if jwttoken == ""{
			fmt.Println("No Authorization")
			return
		}
		jwttoken = strings.Split(jwttoken, "Bearer ")[1]
		token, err := jwt.Parse(jwttoken, func(token *jwt.Token) (interface{}, error) {
			return []byte("Testing key"), nil
		})
		if err != nil{
			fmt.Println("HasSuperUserScope failed")
			fmt.Println(err)
			helpers.ErrorJSON(w,err, 400)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["admin"] != "True"{
			err := errors.New("Failed Admin Check")
			helpers.ErrorJSON(w,err,400)
			return
		}
		ctx := context.WithValue(r.Context(), "userid", claims["userId"])
		var exists bool
		db.db.QueryRow("SELECT UserID FROM tblAdminUsers WHERE UserID = ?", claims["userId"]).Scan(&exists)
		if exists == false{
			fmt.Println("User not in Admin, HasAdminScope has failed")
			err := errors.New("failed admin check")
			helpers.ErrorJSON(w,err, 400)
			return
		}
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}


















































// CustomClaims contains custom data we want from the token.
// type CustomClaims struct {
// 	Scope string `json:"scope"`
// }

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
// func (c CustomClaims) Validate(ctx context.Context) error {
// 	return nil
// }

// EnsureValidToken is a middleware that will check the validity of our JWT.
// The return is a function that also returns an http.Handler.
// func EnsureValidToken() func(next http.Handler) http.Handler {
// 	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
// 	if err != nil {
// 		log.Fatalf("Failed to parse the issuer url: %v", err)
// 	}

// 	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

// 	jwtValidator, err := validator.New(
// 		provider.KeyFunc,
// 		validator.RS256,
// 		issuerURL.String(),
// 		[]string{os.Getenv("AUTH0_API_AUDIENCE")},

// 		validator.WithCustomClaims(
// 			func() validator.CustomClaims {
// 				return &CustomClaims{}
// 			},
// 		),
// 		validator.WithAllowedClockSkew(time.Minute),
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to set up the jwt validator")
// 	}

// 	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
// 		log.Printf("Encountered error while validating JWT: %v", err)

// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusUnauthorized)
// 		w.Write([]byte(`{"message":"Failed to validate JWT."}`))
// 	}
// 	// After creating the errorHandler, we are going to be handling any possible issue that could arise from checking jwt. This is then passed on line 71.
// 	middleware := jwtmiddleware.New(
// 		jwtValidator.ValidateToken,
// 		jwtmiddleware.WithErrorHandler(errorHandler),
// 	)

// 	return func(next http.Handler) http.Handler {
// 		return middleware.CheckJWT(next)
// 	}
// }