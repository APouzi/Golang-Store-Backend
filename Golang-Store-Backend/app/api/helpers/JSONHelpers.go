package helpers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)



func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)//We are decoding into the data here.
	if err != nil{
		return err
	}
	// The purpose of this code is to decode a JSON-encoded input stream and discard the decoded value. It is commonly used when you only want to check if the input stream is a valid JSON data, without actually processing the decoded value. &struct{}{} creates a pointer to an empty struct literal. This is used as a placeholder to store the decoded JSON value, which is discarded afterwards.The purpose of this code is to decode a JSON-encoded input stream and discard the decoded value. It is commonly used when you only want to check if the input stream is a valid JSON data, without actually processing the decoded value.
	err = dec.Decode(&struct{}{})
	// When asking "err != io.EOF", we ask this code to see if there isn't a single value of JSON being passed in.
	if err != io.EOF{
		return errors.New("body must have only a single json objects")
	}

	// Because we aree returning an error, here we are saying that we have went through this function and it successfully completed.
	return nil

}

func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil{
		return err
	}

	if len(headers) >0 {
		for key, value := range headers[0]{
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil{
		return err
	}

	return nil
}

//Lets write the json error if there was an error in the json.

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error{
	statusCode := http.StatusBadRequest

	if len(status) > 0{
		statusCode = status[0]
	}

	var payLoad ErrorJSONResponse
	payLoad.Error = true
	payLoad.Message = err.Error()

	WriteJSON(w, statusCode, payLoad)

	return nil
} 