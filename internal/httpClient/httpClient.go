package httpclient // Package 'httpclient' provides utilities for handling HTTP requests and responses

import (
	"encoding/json" // For encoding and decoding JSON data
	"fmt"           // For formatted I/O operations

	"github.com/valyala/fasthttp" // Importing the third-party package 'fasthttp' for handling HTTP client operations
)

// 'HttpRequestDetails' is a struct that encapsulates all information required for an HTTP request.
type HttpRequestDetails struct {
	URL         string            // URL of the HTTP request
	Method      string            // HTTP request method (GET, POST etc.)
	Headers     map[string]string // HTTP headers
	RequestBody string            // Body of the HTTP request; used in POST requests
}

// 'HttpResponseDetails' is a struct that holds the details of an HTTP response.
type HttpResponseDetails struct {
	Body     []byte      // Raw response body
	JsonData interface{} // Decoded JSON response data
	Error    error       // Error (if any) while making the HTTP request or parsing the response
}

// Function 'SendHttpRequest' takes in an object of HttpRequestDetails,
// makes the HTTP request using fasthttp library and returns the response as HttpResponseDetails object.
func SendHttpRequest(details HttpRequestDetails) HttpResponseDetails {
	req := fasthttp.AcquireRequest()     // Acquires an HTTP request instance
	resp := fasthttp.AcquireResponse()   // Acquires an HTTP response instance
	defer fasthttp.ReleaseRequest(req)   // Make sure to release request instance after it's no longer needed
	defer fasthttp.ReleaseResponse(resp) // Same goes for the response instance

	req.SetRequestURI(details.URL)         // Sets the request URL
	req.Header.SetMethod(details.Method)   // Sets the request method (GET, POST, etc.)
	req.SetBodyString(details.RequestBody) // Sets the request body

	for key, value := range details.Headers { // Iterating through each header and setting it on the request
		req.Header.Set(key, value)
	}

	err := fasthttp.Do(req, resp) // Executes the request and stores the response
	if err != nil {
		return HttpResponseDetails{Error: fmt.Errorf(" Error making request: %v", err)} // If there was an error, return it
	}

	body := resp.Body() // Get the response body

	var jsonData interface{}
	err = json.Unmarshal(body, &jsonData) // Try to unmarshal the response body into JSON format
	if err != nil {
		return HttpResponseDetails{Body: body, Error: nil} // If unable to unmarshal, return the raw body
	}

	return HttpResponseDetails{Body: body, JsonData: jsonData, Error: nil} // Return the response body, unmarshalled JSON data and no error
}
