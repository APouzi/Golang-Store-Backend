package helpers

// Used for testing newly created endpoints
type TestResponse struct{
	Message string `json:"message"`
	StoreOpen bool `json:"storeOpen"`
}

// Error message
type ErrorJSONResponse struct{
	Error bool
	Message string
}

type UserLoginResponse struct{
	UserID string `json:"UserID"`
	Email string `json:"Email"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	IsAdmin bool `json:"IsAdmin"`
	IsSuperUser bool `json:"IsSuperUser"`
	CSRF string `json:"Token"`
}

type UserLoginRequest struct{
	Email string `json:"Email"`
	CSRF string `json:"Token"`
	Pass string `json:"Pass"`
}