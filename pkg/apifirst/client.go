package apifirst

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type Client struct {
	HttpClient    *http.Client
	OpenAPIDoc    *openapi3.T
	Operations    map[string]*Operation
	BaseURL       *url.URL
	UserAgent     string
	Authorization *string
}

type Operation struct {
	*openapi3.Operation
	Method string
	Path   string
}

const CLIENT_USER_AGENT = "CEAI CLI Version 0.1"

func NewClient(OpenAPIDoc *openapi3.T) *Client {
	var client Client

	client.OpenAPIDoc = OpenAPIDoc
	client.UserAgent = CLIENT_USER_AGENT

	client.HttpClient = http.DefaultClient //TODO Better option?

	//Create Operations dictionary
	operations := make(map[string]*Operation)

	for path, pathItem := range OpenAPIDoc.Paths {
		oapiOperations := pathItem.Operations()

		for method, oapiOperation := range oapiOperations {

			operations[oapiOperation.OperationID] = &Operation{
				Operation: oapiOperation,
				Method:    method,
				Path:      path,
			}
		}
	}

	client.Operations = operations

	//TODO Set ServerUrl if present in spec?

	return &client
}

func (c *Client) SetServerUrl(serverUrl string) error {
	oUrl, err := url.Parse(serverUrl)

	if err != nil {
		return err
	}

	c.BaseURL = oUrl

	return nil
}

func (c *Client) SetAuthorization(authorization string) {
	c.Authorization = &authorization
}

func (c *Client) GenerateUrl(operation *Operation, parameters *map[string]string) (*url.URL, error) {
	path := operation.Path

	//Query Parameters object
	var queryValues *url.Values

	for _, parm := range operation.Parameters {
		if parm.Value.In == "path" {
			value, exists := (*parameters)[parm.Value.Name]

			if parm.Value.Required && !exists {
				return nil, fmt.Errorf("missing required path parameter %s", parm.Value.Name)
			}

			path = strings.Replace(path, "{"+parm.Value.Name+"}", url.PathEscape(value), 1)
		}

		if parm.Value.In == "query" {
			value, exists := (*parameters)[parm.Value.Name]

			if parm.Value.Required && !exists {
				return nil, fmt.Errorf("missing required query parameter %s", parm.Value.Name)
			}

			if queryValues == nil {
				queryValues = &url.Values{}
			}

			queryValues.Set(parm.Value.Name, value)
		}
	}

	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)

	if queryValues != nil {
		u.RawQuery = queryValues.Encode()
	}

	return u, nil
}

func (c *Client) Request(operationId string, parameters *map[string]string, requestBody interface{}) (*Response, error) {
	operation, exists := c.Operations[operationId]

	if !exists {
		return nil, fmt.Errorf("Operation %s not found in OpenAPI specification", operationId)
	}

	u, err := c.GenerateUrl(operation, parameters)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter

	if requestBody != nil && operation.RequestBody != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(requestBody)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(operation.Method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	if c.Authorization != nil {
		req.Header.Set("Authorization", *c.Authorization)
	}

	//TODO Request validation? Optional?

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	//defer resp.Body.Close()

	//TODO Response validation?? Optional?

	oapiResponse := operation.Responses.Get(resp.StatusCode)

	response, err := NewResponse(resp, oapiResponse.Value)

	if err != nil {
		return nil, err
	}

	return response, nil
}
