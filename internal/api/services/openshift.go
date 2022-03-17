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

const OPENSHIFT_LABORATORY_LABEL string = "ceai-laboratory"
const OPENSHIFT_DESCRIPTION_ANNOTATION string = "openshift.io/description"
const OPENSHIFT_DISPLAY_NAME_ANNOTATION string = "openshift.io/display-name"
const OPENSHIFT_REQUESTER_ANNOTATION string = "openshift.io/requester"

func MapOpenshiftProject(project *openshiftproject.Project) *models.OpenshiftProject {
	var openshiftProject models.OpenshiftProject

	openshiftProject.Id = project.Name
	openshiftProject.Description = project.Annotations[OPENSHIFT_DESCRIPTION_ANNOTATION]
	openshiftProject.Displayname = project.Annotations[OPENSHIFT_DISPLAY_NAME_ANNOTATION]

	return &openshiftProject
}

func MapOpenshiftProjectWithLab(project *openshiftproject.Project) *models.OpenshiftProjectWithLab {
	var openshiftProject models.OpenshiftProjectWithLab

	openshiftProject.OpenshiftProject = *MapOpenshiftProject(project)

	if idLab, ok := project.Labels[OPENSHIFT_LABORATORY_LABEL]; ok {
		openshiftProject.IdLab = idLab
	} else {
		openshiftProject.IdLab = "none"
	}

	return &openshiftProject
}

func MapOpenshiftProjectWithMeta(project *openshiftproject.Project) *models.OpenshiftProjectWithMeta {
	var openshiftProject models.OpenshiftProjectWithMeta

	openshiftProject.OpenshiftProjectWithLab = *MapOpenshiftProjectWithLab(project)

	if requester, ok := project.Annotations[OPENSHIFT_REQUESTER_ANNOTATION]; ok {
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
			if _, ok := project.Annotations[OPENSHIFT_REQUESTER_ANNOTATION]; ok {
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

func UpdateOpenshiftProject(projectId string, updateParam *models.OpenshiftProjectUpdate) error {
	project, err := openshift.GetProject(projectId)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_OPENSHIFT)
	}

	//If lab is changed detach project and attach to new lab
	if updateParam.IdLab != nil {
		if currentIdLab, ok := project.Labels[OPENSHIFT_LABORATORY_LABEL]; ok {
			if currentIdLab != *updateParam.IdLab {

				err = DetachOpenshiftProjectFromLaboratory(currentIdLab, projectId)

				if err != nil {
					return err
				}

				err = AttachOpenshiftProjectToLaboratory(*updateParam.IdLab, projectId)

				if err != nil {
					return err
				}

				//Obtain updated project
				project, err = openshift.GetProject(projectId)

				if err != nil {
					return err
				}
			}
		}
	}

	if updateParam.Description != nil {
		project.Annotations[OPENSHIFT_DESCRIPTION_ANNOTATION] = *updateParam.Description
	}

	if updateParam.Displayname != nil {
		project.Annotations[OPENSHIFT_DISPLAY_NAME_ANNOTATION] = *updateParam.Displayname
	}

	_, err = openshift.UpdateProject(project)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_OPENSHIFT)
	}

	return nil
}
