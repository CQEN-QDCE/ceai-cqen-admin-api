package handlers

import (
	"sync"

	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/keycloak"
	"github.com/Nerzal/gocloak/v8"
)

type ServerHandlers struct {
	UserHandlersInterface
	LaboratoryHandlersInterface
}

//Common handlers assets

const LAB_TOP_GROUP = "/Laboratories/"

const KEYCLOAK_LAB_TOP_GROUP = "Laboratories"
const AWS_LAB_GROUP_PREFIX = "Lab_"
const OPENSHIFT_LAB_GROUP_PREFIX = "Lab_"

const ADMIN_ROLE_NAME = "Admin"
const DEV_ROLE_NAME = "Developer"

func GetKeycloakAdminGroup() (*gocloak.Group, error) {
	return keycloak.GetGroup(ADMIN_ROLE_NAME)
}

// Parallelize parallelizes function calls
func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

func RemoveStringElementFromArray(array []string, element string) []string {
	index := -1

	for i, e := range array {
		if e == element {
			index = i
			break
		}
	}

	if index != -1 {
		ret := make([]string, 0)
		ret = append(ret, array[:index]...)
		return append(ret, array[index+1:]...)
	} else {
		return array
	}
}
