package apifirst

import (
	"encoding/json"
	"errors"
	"net/http"
)

const DefaultContentType = "application/json"

type Response struct {
	http.ResponseWriter
	Status      int
	Body        interface{}
	ContentType string
}

func NewResponse(w *http.ResponseWriter) *Response {
	response := &Response{
		ResponseWriter: *w,
	}

	response.SetStatus(http.StatusOK)
	response.SetContentType(DefaultContentType)

	return response
}

func (r *Response) SetStatus(status int) {
	r.Status = status
}

func (r *Response) SetContentType(contentType string) {
	r.ContentType = contentType
	r.Header().Set("Content-Type", r.ContentType)
}

func (r *Response) SetBody(body interface{}) {
	r.Body = body
}

func (r *Response) GetMarshaledBody() ([]byte, error) {
	if r.Body != nil {

		switch r.ContentType {
		case "application/json":
			return json.Marshal(r.Body)
		case "":
		}

	}

	err := errors.New("Unsupported content type")

	return nil, err
}

func (r *Response) WriteResponse() {
	r.WriteHeader(r.Status)

	mBody, err := r.GetMarshaledBody()

	if err != nil {
		//error
	}

	r.Write(mBody)
}
