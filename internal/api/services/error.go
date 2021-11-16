package services

const (
	ERROR_SERVER_KEYCLOAK  string = "Keycloak"
	ERROR_SERVER_OPENSHIFT string = "Openshift"
	ERROR_SERVER_AWS       string = "AWS"
)

type ErrorExternal struct {
	Err    error
	Server string
}

func (e ErrorExternal) Error() string {
	return e.Server + " error: " + e.Err.Error()
}

type ErrorExternalRessourceNotFound struct {
	*ErrorExternal
}

func NewErrorExternalRessourceNotFound(err error, server string) ErrorExternalRessourceNotFound {
	return ErrorExternalRessourceNotFound{
		ErrorExternal: &ErrorExternal{
			Err:    err,
			Server: server,
		},
	}
}

type ErrorExternalRessourceExist struct {
	*ErrorExternal
}

func NewErrorExternalRessourceExist(err error, server string) ErrorExternalRessourceExist {
	return ErrorExternalRessourceExist{
		ErrorExternal: &ErrorExternal{
			Err:    err,
			Server: server,
		},
	}
}

type ErrorExternalServerDown struct {
	*ErrorExternal
}

func NewErrorExternalServerDown(err error, server string) ErrorExternalServerDown {
	return ErrorExternalServerDown{
		ErrorExternal: &ErrorExternal{
			Err:    err,
			Server: server,
		},
	}
}

type ErrorExternalServerError struct {
	*ErrorExternal
}

func NewErrorExternalServerError(err error, server string) ErrorExternalServerError {
	return ErrorExternalServerError{
		ErrorExternal: &ErrorExternal{
			Err:    err,
			Server: server,
		},
	}
}

/*
if errors.As(err, ErrorExternalRessourceNotFound) {
    // err is a *QueryError, and e is set to the error's value
}
*/
