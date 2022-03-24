package apifirst

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

type Response struct {
	*http.Response
	OapiResponse *openapi3.Response
}

func NewResponse(httpResponse *http.Response, oapiResponse *openapi3.Response) (*Response, error) {

	response := Response{
		Response:     httpResponse,
		OapiResponse: oapiResponse,
	}

	if httpResponse.StatusCode > 300 {
		//Send error with specification message
		if oapiResponse.Description != nil {
			return nil, fmt.Errorf(*oapiResponse.Description)
		} else {
			return nil, fmt.Errorf("Server Error: %s", httpResponse.Status)
		}
	}

	return &response, nil
}

func (r *Response) UnmarshalBody(v interface{}) error {
	if r.Body != nil {
		switch r.Header.Get("Content-Type") {
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&v)
			return err
		case "": //TODO Support more contentTypes?
		}
	} else {
		return nil
	}

	err := errors.New("Unsupported content type")

	return err
}
