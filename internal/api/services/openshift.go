package services

import (
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/keycloak"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/openshift"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	openshiftauthorization "github.com/openshift/api/authorization/v1"
	openshiftproject "github.com/openshift/api/project/v1"
	openshiftcore "k8s.io/api/core/v1"
	openshiftmeta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func MapOpenshiftProject(project *openshiftproject.Project) *models.OpenshiftProject {
	var openshiftProject models.OpenshiftProject

	openshiftProject.Id = project.Name
	openshiftProject.Description = project.Annotations["openshift.io/description"]
	openshiftProject.Displayname = project.Annotations["openshift.io/display-name"]

	return &openshiftProject
}

func MapOpenshiftProjectWithLab(project *openshiftproject.Project) *models.OpenshiftProjectWithLab {
	var openshiftProject models.OpenshiftProjectWithLab

	openshiftProject.OpenshiftProject = MapOpenshiftProject(project)

	if idLab, ok := project.Labels["ceai-laboratory"]; ok {
		openshiftProject.IdLab = idLab
	} else {
		openshiftProject.IdLab = "none"
	}

	return &openshiftProject
}

func MapOpenshiftProjectWithMeta(project *openshiftproject.Project) *models.OpenshiftProjectWithMeta {
	var openshiftProject models.OpenshiftProjectWithMeta

	openshiftProject.OpenshiftProjectWithLab = MapOpenshiftProjectWithLab(project)

	if requester, ok := project.Annotations["openshift.io/requester"]; ok {
		openshiftProject.Requester = &requester
	}

	openshiftProject.CreationDate = &project.CreationTimestamp.Time

	return &openshiftProject
}

func CreateGroupRoleBinding(groupName string, projectNamespace string) (*openshiftauthorization.RoleBinding, error) {
	rolebinding := openshiftauthorization.RoleBinding{
		ObjectMeta: openshiftmeta.ObjectMeta{
			Name: groupName + "_" + projectNamespace,
		},
		Subjects: []openshiftcore.ObjectReference{
			{
				Kind: "Group",
				Name: groupName,
			},
		},
		RoleRef: openshiftcore.ObjectReference{
			Kind: "ClusterRole",
			Name: "edit",
		},
	}

	return openshift.CreateRoleBinding(projectNamespace, &rolebinding)
}

func GetOpenshiftProjects() ([]*models.OpenshiftProjectWithMeta, error) {
	projects, err := openshift.GetProjects()

	if err != nil {
		return nil, NewErrorExternalServerError(err, ERROR_SERVER_OPENSHIFT)
	}

	projectList := make([]*models.OpenshiftProjectWithMeta, 0)

	if projects != nil {

		for _, project := range *projects {
			if _, ok := project.Annotations["openshift.io/requester"]; ok {
				projectList = append(projectList, MapOpenshiftProjectWithMeta(&project))
			}
		}
	}

	return projectList, nil
}

func GetOpenshiftProjectFromId(projectId string) (*models.OpenshiftProjectWithMeta, error) {
	project, err := openshift.GetProject(projectId)

	if err != nil {
		return nil, NewErrorExternalRessourceNotFound(err, ERROR_SERVER_OPENSHIFT)
	}

	return MapOpenshiftProjectWithMeta(project), nil
}

func CreateOpenshiftProject(createParam *models.OpenshiftProjectWithLab) error {
	//Validate that lab exist
	_, err := keycloak.GetGroup(createParam.IdLab)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_KEYCLOAK)
	}

	//Create new project
	projectRequest := openshiftproject.ProjectRequest{
		DisplayName: createParam.Displayname,
		Description: createParam.Description,
		ObjectMeta: openshiftmeta.ObjectMeta{
			Name: createParam.Id,
		},
	}

	_, err = openshift.CreateProject(&projectRequest)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_OPENSHIFT)
	}

	return AttachOpenshiftProjectToLaboratory(createParam.IdLab, createParam.Id)
}
