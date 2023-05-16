package authorization

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Apouzi/golang-shop/app/api/helpers"
	"github.com/golang-jwt/jwt/v5"
)


type JWTtest struct{
	Token string `json:"JWT"`
}

func ValidateToken(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		jwttest := &JWTtest{}
		helpers.ReadJSON(w, r, &jwttest)
		token, err := jwt.Parse(jwttest.Token, func(token *jwt.Token) (interface{}, error) {
			return []byte("Testing key"), nil
		})
		if err != nil{
			fmt.Println("middleware test error")
			fmt.Println(err)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println("email claim",claims["email"])
		ctx := context.WithValue(r.Context(), "email",claims["email"])
		if token.Valid{
			fmt.Println("token validated in middleware")
		}
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}


// HasScope checks whether our claims have a specific scope.
// func (c CustomClaims) HasScope(expectedScope string) bool {
// 	fmt.Println("Validate hit -  scope")
//     result := strings.Split(c.Scope, " ")
//     for i := range result {
//         if result[i] == expectedScope {
//             return true
//         }
//     }

//     return false
// }


















































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