package apifirst

import (
	"encoding/json"
	"errors"
	"net/http"
)

const DefaultContentType = "application/json"

type ResponseWriter struct {
	http.ResponseWriter
	Status      int
	Body        interface{}
	ContentType string
}

func NewResponseWriter(w *http.ResponseWriter) *ResponseWriter {
	response := &ResponseWriter{
		ResponseWriter: *w,
	}

	response.SetStatus(http.StatusOK)
	response.SetContentType(DefaultContentType)

	return response
}

func (r *ResponseWriter) SetStatus(status int) {
	r.Status = status
}

func (r *ResponseWriter) SetContentType(contentType string) {
	r.ContentType = contentType
	r.Header().Set("Content-Type", r.ContentType)
}

func (r *ResponseWriter) SetBody(body interface{}) {
	r.Body = body
}

func (r *ResponseWriter) GetMarshaledBody() ([]byte, error) {
	if r.Body != nil {

		switch r.ContentType {
		case "application/json":
			return json.Marshal(r.Body)
		case "": //TODO Support more contentTypes?
		}
	} else {
		return nil, nil
	}

	err := errors.New("unsupported content type")

	return nil, err
}

func (r *ResponseWriter) WriteResponse() error {
	r.WriteHeader(r.Status)

	mBody, err := r.GetMarshaledBody()

	if err != nil {
		return err
	}

	r.Write(mBody)

	return nil
}
