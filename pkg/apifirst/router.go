// API First Approach Router
//
// Made from https://github.com/getkin/kin-openapi/blob/master/routers/gorillamux/router.go
//
// Bind routes to Handler functions named after the OperationId property of operations in a OpenAPI 3.0 schema
//
// Validation of request and response are performed before and after handlers.

package apifirst

import (
	"context"
	"log"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/gorilla/mux"

	_ "github.com/CQEN-QDCE/ceai-cqen-admin-api/api"
)

// Router helps link http.Request.s and an OpenAPIv3 spec
type Router struct {
	Muxes    []*mux.Route
	Routes   []*routers.Route
	Router   *mux.Router
	Handlers interface{}
	Options  *RouterOptions
}

type RouterOptions struct {
	AuthenticationFunc *openapi3filter.AuthenticationFunc
	//Add more options as needed
}

// NewRouter creates a gorilla/mux router with handlers attached via the CallRouteFunc function
// Assumes spec is .Validate()d
func NewRouter(doc *openapi3.T, serverWrapper interface{}, options *RouterOptions) *Router {
	r := &Router{}

	r.Handlers = serverWrapper

	muxRouter := mux.NewRouter() /*.UseEncodedPath()?*/

	log.Println("API First router initialization")

	for _, path := range orderedPaths(doc.Paths) {
		pathItem := doc.Paths[path]

		operations := pathItem.Operations()
		methods := make([]string, 0, len(operations))
		for method := range operations {
			methods = append(methods, method)

			//Closures for Http methods handlers
			op := operations[method]

			muxRoute := muxRouter.HandleFunc(path, func(w http.ResponseWriter, request *http.Request) {
				response := r.CallRouteFunc(op, w, request)
				response.WriteResponse()
				log.Printf("%v %v %v", request.Method, request.RequestURI, response.Status)
			}).Methods(method)

			r.Muxes = append(r.Muxes, muxRoute)

			r.Routes = append(r.Routes, &routers.Route{
				Spec:      doc,
				Server:    nil,
				Path:      path,
				PathItem:  pathItem,
				Method:    method,
				Operation: op,
			})

			log.Printf("Route %v %v attached to handler %v \n", method, path, operations[method].OperationID)
		}
	}

	r.Router = muxRouter

	r.Options = options

	return r
}

// FindRoute extracts the route and parameters of an http.Request
func (r *Router) FindRoute(req *http.Request) (*routers.Route, map[string]string, error) {
	for i, muxRoute := range r.Muxes {
		var match mux.RouteMatch

		if muxRoute.Match(req, &match) {

			//Ensure there is no error in the match
			if err := match.MatchErr; err == nil {
				route := r.Routes[i]
				route.Method = req.Method
				route.Operation = route.Spec.Paths[route.Path].GetOperation(route.Method)
				return route, match.Vars, nil
			}
		}
	}

	return nil, nil, routers.ErrPathNotFound
}

func (r *Router) Serve(port string) error {
	//TODO Better start logs, Fiber Style?
	log.Printf("listening incoming requests on port %v \n", port)

	return http.ListenAndServe(":"+port, r.Router)
}

// Call the handler method associated with request route
// Validate request and response against OpenAPI Spec
// Then return a apifirst.Response
//
// TODO: This method is way too huge. Need to split/use middlewares?
//
func (r *Router) CallRouteFunc(operation *openapi3.Operation, w http.ResponseWriter, request *http.Request) *Response {
	//Convert ResponseWriter to apifirst.response
	response := NewResponse(&w)

	//Find handler method using Reflect package
	handlerFunc := operation.OperationID

	inputs := make([]reflect.Value, 2)

	v := reflect.ValueOf(r.Handlers)
	m := v.MethodByName(handlerFunc)

	//test m, return unimplemented response if method is undefined
	if !m.IsValid() {
		response.SetStatus(http.StatusNotImplemented)
		return response
	}

	//Find route in spec and extract path params
	route, pathParams, err := r.FindRoute(request)
	if err != nil {
		//Could not match request with any route in the OpenAPI spec
		response.SetStatus(http.StatusNotFound)
		return response
	}

	//Prepare request validation
	filterOptions := openapi3filter.Options{}
	if r.Options != nil {
		if r.Options.AuthenticationFunc != nil {
			filterOptions.AuthenticationFunc = *r.Options.AuthenticationFunc
		}
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    request,
		PathParams: pathParams,
		Route:      route,
		Options:    &filterOptions,
	}

	if err := openapi3filter.ValidateRequest(context.Background(), requestValidationInput); err != nil {

		if e, ok := err.(*openapi3filter.SecurityRequirementsError); ok {
			log.Println(e.Error())
			log.Println(e.Errors)
			response.SetStatus(http.StatusUnauthorized)
		} else {
			log.Println(e.Error())
			response.SetStatus(http.StatusBadRequest)

			//TODO Add switch in .env to output validation error or not
			//TODO Seems I can't output text with this Content-Type...
			//response.SetContentType("text/plain")
			response.SetBody(err.Error())
		}

		return response
	}

	//Call method

	inputs[0] = reflect.ValueOf(response)
	inputs[1] = reflect.ValueOf(request)

	m.Call(inputs)

	//TODO Writing Header at this moment could be problematic if response validation fails?
	//response.WriteResponse()

	//Validate response
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 response.Status,
		Header:                 response.Header(),
	}

	if response.Body != nil {
		data, _ := response.GetMarshaledBody()
		responseValidationInput.SetBodyBytes(data)
	}

	if err := openapi3filter.ValidateResponse(context.Background(), responseValidationInput); err != nil {
		response.SetStatus(http.StatusBadGateway)

		//TODO Add switch in .env to output validation error or not
		//response.SetContentType("text/plain")
		response.SetBody(err.Error())

		return response
	}

	return response
}

func orderedPaths(paths map[string]*openapi3.PathItem) []string {
	// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.3.md#pathsObject
	// When matching URLs, concrete (non-templated) paths would be matched
	// before their templated counterparts.
	// NOTE: sorting by number of variables ASC then by lexicographical
	// order seems to be a good heuristic.
	vars := make(map[int][]string)
	max := 0
	for path := range paths {
		count := strings.Count(path, "}")
		vars[count] = append(vars[count], path)
		if count > max {
			max = count
		}
	}
	ordered := make([]string, 0, len(paths))
	for c := 0; c <= max; c++ {
		if ps, ok := vars[c]; ok {
			sort.Strings(ps)
			ordered = append(ordered, ps...)
		}
	}
	return ordered
}
