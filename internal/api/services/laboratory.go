package services

import (
	"errors"

	scim "github.com/CQEN-QDCE/aws-sso-scim-goclient"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/api/globalvar"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/aws"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/keycloak"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/api/openshift"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/models"
	"github.com/Nerzal/gocloak/v11"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
	openshiftuser "github.com/openshift/api/user/v1"
	openshiftmeta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const LAB_TOP_GROUP = "/Laboratories/"

const KEYCLOAK_LAB_TOP_GROUP = "Laboratories"
const OPENSHIFT_LAB_GROUP_PREFIX = "Lab_"

const KEYCLOAK_OPENSHIFT_PROJECT_ATTRIBUTE = "openshift_projects"
const KEYCLOAK_AWS_ACCOUNT_ATTRIBUTE = "aws_accounts"

type LaboratoryState struct {
	Keycloak  *gocloak.Group
	Aws       *scim.Group
	Openshift *openshiftuser.Group
}

func MapLaboratory(kgroup gocloak.Group) (*models.Laboratory, error) {
	var lab models.Laboratory

	lab.Id = *kgroup.Name

	//TODO test required attributes
	if kgroup.Attributes != nil {
		attributes := *kgroup.Attributes

		if attributes["displayname"] != nil {
			lab.Displayname = attributes["displayname"][0]
		}

		if attributes["type"] != nil {
			lab.Type = attributes["type"][0]
		}

		if attributes["description"] != nil {
			lab.Description = attributes["description"][0]
		}

		if attributes["gitrepo"] != nil {
			lab.Gitrepo = &attributes["gitrepo"][0]
		}
	}

	return &lab, nil
}

func MapLaboratoryWithUsers(kgroup gocloak.Group) (*models.LaboratoryWithUsers, error) {
	var lab models.LaboratoryWithUsers

	laboratory, _ := MapLaboratory(kgroup)

	lab.Laboratory = *laboratory

	members, err := keycloak.GetGroupMembers(&kgroup)

	if err != nil {
		return nil, err
	}

	var users []string

	for _, member := range members {
		users = append(users, *member.Email)
	}

	lab.Users = &users

	return &lab, nil
}

func MapLaboratoryWithResources(kgroup gocloak.Group) (*models.LaboratoryWithResources, error) {
	var lab models.LaboratoryWithResources

	laboratory, _ := MapLaboratoryWithUsers(kgroup)

	lab.LaboratoryWithUsers = *laboratory

	attributes := *kgroup.Attributes

	var openshiftprojects []models.OpenshiftProject

	if attributes["openshift_projects"] != nil {
		for _, projectName := range attributes["openshift_projects"] {
			var project models.OpenshiftProject

			//TODO this is getting not mapping....
			oProject, err := openshift.GetProject(projectName)

			if err == nil {
				project.Id = oProject.Name
				project.Description = oProject.Annotations["openshift.io/description"]
				project.Displayname = oProject.Annotations["openshift.io/display-name"]

				openshiftprojects = append(openshiftprojects, project)
			}
			//else ignore project TODO Log?
		}

		lab.Openshiftprojects = &openshiftprojects
	}

	var awsAccounts []models.AWSAccount

	if attributes["aws_accounts"] != nil {
		for _, accountId := range attributes["aws_accounts"] {
			var account models.AWSAccount

			//TODO service function
			accountInfo, err := aws.DescribeAccount(accountId)

			if err == nil {
				account.Email = accountInfo.Email
				account.Id = accountInfo.Id
				account.Name = accountInfo.Name

				awsAccounts = append(awsAccounts, account)
			}
			//else ignore account TODO log?
		}

		lab.AWSAccounts = &awsAccounts
	}

	return &lab, nil
}

func CreateGroupKeycloak(plab *models.Laboratory) error {
	labTopGroup, err := keycloak.GetGroup(KEYCLOAK_LAB_TOP_GROUP)

	if err != nil {
		return err
	}

	attributes := map[string][]string{
		"description": {plab.Description},
		"type":        {plab.Type},
		"displayname": {plab.Displayname},
	}

	if plab.Gitrepo != nil {
		attributes["gitrepo"] = []string{*plab.Gitrepo}
	}

	kgroup := gocloak.Group{
		Name:       &plab.Id,
		Attributes: &attributes,
	}

	return keycloak.CreateChildGroup(labTopGroup, &kgroup)
}

func CreateGroupAws(plab *models.Laboratory) error {
	group := scim.NewGroup(AWS_LAB_GROUP_PREFIX + plab.Id)

	_, err := aws.CreateGroup(group)

	return err
}

func CreateGroupOpenshift(plab *models.Laboratory) error {
	group := openshiftuser.Group{
		ObjectMeta: openshiftmeta.ObjectMeta{
			Name: OPENSHIFT_LAB_GROUP_PREFIX + plab.Id,
		},
	}

	_, err := openshift.CreateGroup(&group)

	return err
}

//Gets current Lab states across all products: Keycloak|AWS|Openshift
func GetLaboratoryState(idLab string) (*LaboratoryState, error) {
	var state LaboratoryState
	var kerr, aerr, oerr error

	fKeycloak := func() {
		state.Keycloak, kerr = keycloak.GetGroup(idLab)
	}

	fAws := func() {
		state.Aws, aerr = aws.GetGroup(AWS_LAB_GROUP_PREFIX + idLab)
	}

	fOpenshift := func() {
		state.Openshift, oerr = openshift.GetGroup(OPENSHIFT_LAB_GROUP_PREFIX + idLab)
	}

	Parallelize(fKeycloak, fAws, fOpenshift)

	if kerr != nil || oerr != nil || aerr != nil {
		var err error

		if kerr != nil {
			err = NewErrorExternalRessourceNotFound(kerr, ERROR_SERVER_KEYCLOAK)
		} else if oerr != nil {
			err = NewErrorExternalRessourceNotFound(oerr, ERROR_SERVER_OPENSHIFT)
		} else if aerr != nil {
			err = NewErrorExternalRessourceNotFound(aerr, ERROR_SERVER_AWS)
		}

		return nil, err
	}

	return &state, nil
}

func AddUsersToGroupKeycloak(labState *LaboratoryState, newUsersList []*UserState) error {
	group := labState.Keycloak

	for _, userState := range newUsersList {

		kuser := userState.Keycloak

		err := keycloak.AddUserToGroup(kuser, group)

		if err != nil {
			return err
		}
	}

	return nil
}

func AddUsersToGroupAws(labState *LaboratoryState, newUsersList []*UserState) error {
	group := labState.Aws

	for _, userState := range newUsersList {

		user := userState.Aws

		err := aws.AddUserToGroup(user, group)

		if err != nil {
			return err
		}
	}

	return nil
}

func AddUsersToGroupOpenshift(labState *LaboratoryState, newUsersList []*UserState) error {
	group := labState.Openshift

	for _, userState := range newUsersList {

		user := userState.Openshift

		err := openshift.AddUserInGroup(user.Name, group.Name)

		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveUsersFromGroupKeycloak(labState *LaboratoryState, newUsersList []*UserState) error {
	group := labState.Keycloak

	for _, userState := range newUsersList {

		kuser := userState.Keycloak

		err := keycloak.DeleteUserFromGroup(kuser, group)

		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveUsersFromGroupAws(labState *LaboratoryState, newUsersList []*UserState) error {
	group := labState.Aws

	for _, userState := range newUsersList {

		user := userState.Aws

		err := aws.RemoveUserFromGroup(user, group)

		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveUsersFromGroupOpenshift(labState *LaboratoryState, newUsersList []*UserState) error {
	group := labState.Openshift

	for _, userState := range newUsersList {

		user := userState.Openshift

		err := openshift.RemoveUserFromGroup(user.Name, group.Name)

		if err != nil {
			return err
		}
	}

	return nil
}

func GetKeycloakUserList(userList *[]string) (*[]gocloak.User, error) {
	var users []gocloak.User

	for _, username := range *userList {
		user, err := keycloak.GetUser(username)

		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return &users, nil
}

func AddElementToKeycloakGroupArrayAttribute(kgroup *gocloak.Group, attribute string, element string) error {
	var AttrElements []string

	if _, ok := (*kgroup.Attributes)[attribute]; ok {
		AttrElements = (*kgroup.Attributes)[attribute]

		//Verify that element is not already there
		for _, attrElem := range AttrElements {
			if attrElem == element {
				return nil
			}
		}

		AttrElements = append(AttrElements, element)
	} else {
		AttrElements = []string{element}
	}

	(*kgroup.Attributes)[attribute] = AttrElements

	return keycloak.UpdateGroup(kgroup)
}

func RemoveElementFromKeycloakGroupArrayAttribute(kgroup *gocloak.Group, attribute string, element string) error {
	AttrElements := (*kgroup.Attributes)[attribute]

	if len(AttrElements) > 0 {
		AttrElements = RemoveStringElementFromArray(AttrElements, element)

		(*kgroup.Attributes)[attribute] = AttrElements

		return keycloak.UpdateGroup(kgroup)
	}

	return nil
}

func GetLaboratories() ([]*models.LaboratoryWithResources, error) {
	labGroups, err := keycloak.GetGroups(gocloak.StringP(KEYCLOAK_LAB_TOP_GROUP))

	if err != nil {
		return nil, NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	labsList := make([]*models.LaboratoryWithResources, 0, len(*labGroups))

	if labGroups != nil {

		for _, group := range *labGroups {
			lab, err := MapLaboratoryWithResources(group)

			if err == nil {
				labsList = append(labsList, lab)
			}
			//TODO Log error?
		}
	}

	return labsList, nil
}

func GetLaboratoryFromId(laboratoryid string) (*models.LaboratoryWithResources, error) {
	labGroup, err := keycloak.GetGroup(laboratoryid)

	if err != nil {
		return nil, NewErrorExternalRessourceNotFound(err, ERROR_SERVER_KEYCLOAK)
	}

	lab, err := MapLaboratoryWithResources(*labGroup)

	if err != nil {
		return nil, NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	return lab, nil
}

func CreateLaboratory(pLab *models.Laboratory) error {
	var kerr, oerr, aerr error

	kfunc := func() {
		kerr = CreateGroupKeycloak(pLab)
	}

	ofunc := func() {
		if globalvar.IsProdEnv() {
			oerr = CreateGroupOpenshift(pLab)
		}
	}

	afunc := func() {
		aerr = CreateGroupAws(pLab)
	}

	Parallelize(kfunc, ofunc, afunc)

	if kerr != nil || oerr != nil || aerr != nil {
		var err error

		if kerr != nil {
			err = NewErrorExternalRessourceExist(kerr, ERROR_SERVER_KEYCLOAK)
		} else if oerr != nil {
			err = NewErrorExternalRessourceExist(oerr, ERROR_SERVER_OPENSHIFT)
		} else if aerr != nil {
			err = NewErrorExternalRessourceExist(aerr, ERROR_SERVER_AWS)
		}

		return err
	}

	return nil
}

func UpdateLaboratory(laboratoryId string, pLab *models.LaboratoryUpdate) error {
	//change info in keycloak group only as groups in AWS and Openshift only has id in them
	group, err := keycloak.GetGroup(laboratoryId)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_KEYCLOAK)
	}

	if pLab.Description != nil {
		(*group.Attributes)["description"] = []string{*pLab.Description}
	}

	if pLab.Displayname != nil {
		(*group.Attributes)["displayname"] = []string{*pLab.Displayname}
	}

	if pLab.Gitrepo != nil {
		(*group.Attributes)["gitrepo"] = []string{*pLab.Gitrepo}
	}

	if pLab.Type != nil {
		(*group.Attributes)["type"] = []string{*pLab.Type}
	}

	err = keycloak.UpdateGroup(group)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	return nil
}

func AddLaboratoryUsers(laboratoryId string, usernameList []string) error {
	labState, err := GetLaboratoryState(laboratoryId)

	if err != nil {
		return err
	}

	var userStates []*UserState

	//Validate and gather user states
	for _, username := range usernameList {

		userState, err := GetUserState(username)

		if err != nil {
			return err
		}

		userStates = append(userStates, userState)
	}

	//Add users to group in each resources
	var kerr, oerr, aerr error

	kfunc := func() {
		kerr = AddUsersToGroupKeycloak(labState, userStates)
	}

	ofunc := func() {
		if globalvar.IsProdEnv() {
			oerr = AddUsersToGroupOpenshift(labState, userStates)
		}
	}

	afunc := func() {
		aerr = AddUsersToGroupAws(labState, userStates)
	}

	Parallelize(kfunc, ofunc, afunc)

	if kerr != nil || oerr != nil || aerr != nil {
		var err error

		if kerr != nil {
			err = NewErrorExternalServerError(kerr, ERROR_SERVER_KEYCLOAK)
		} else if oerr != nil {
			err = NewErrorExternalServerError(oerr, ERROR_SERVER_OPENSHIFT)
		} else if aerr != nil {
			err = NewErrorExternalServerError(aerr, ERROR_SERVER_AWS)
		}

		return err
	}

	return nil
}

func RemoveLaboratoryUsers(laboratoryId string, usernameList []string) error {
	labState, err := GetLaboratoryState(laboratoryId)

	if err != nil {
		return err
	}

	var userStates []*UserState

	//Validate and gather user states
	for _, username := range usernameList {

		userState, err := GetUserState(username)

		if err != nil {
			return err
		}

		userStates = append(userStates, userState)
	}

	//Remove users from group in each resources
	var kerr, oerr, aerr error

	kfunc := func() {
		kerr = RemoveUsersFromGroupKeycloak(labState, userStates)
	}

	ofunc := func() {
		if globalvar.IsProdEnv() {
			oerr = RemoveUsersFromGroupOpenshift(labState, userStates)
		}
	}

	afunc := func() {
		aerr = RemoveUsersFromGroupAws(labState, userStates)
	}

	Parallelize(kfunc, ofunc, afunc)

	if kerr != nil || oerr != nil || aerr != nil {
		var err error

		if kerr != nil {
			err = NewErrorExternalServerError(kerr, ERROR_SERVER_KEYCLOAK)
		} else if oerr != nil {
			err = NewErrorExternalServerError(oerr, ERROR_SERVER_OPENSHIFT)
		} else if aerr != nil {
			err = NewErrorExternalServerError(aerr, ERROR_SERVER_AWS)
		}

		return err
	}

	return nil
}

func AttachOpenshiftProjectToLaboratory(laboratoryId string, projectId string) error {
	kgroup, err := keycloak.GetGroup(laboratoryId)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_KEYCLOAK)
	}

	openshiftGroupName := OPENSHIFT_LAB_GROUP_PREFIX + laboratoryId

	//Add lab label to namespace
	namespace, err := openshift.GetNamespace(projectId)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_OPENSHIFT)
	}

	if namespace.Labels == nil {
		namespace.Labels = map[string]string{}
	} else if idLab, ok := namespace.Labels["ceai-laboratory"]; ok {
		//project already assigned to a lab
		err = errors.New("Project is already attached to laboratory " + idLab)
		return NewErrorExternalRessourceExist(err, ERROR_SERVER_OPENSHIFT)
	}

	namespace.Labels["ceai-laboratory"] = laboratoryId

	_, err = openshift.UpdateNamespace(namespace)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_OPENSHIFT)
	}

	//Add Rolebinding between group and project namespace
	_, err = CreateGroupRoleBinding(openshiftGroupName, projectId)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_OPENSHIFT)
	}

	//Add project to keycloak openshift_projects attribute
	err = AddElementToKeycloakGroupArrayAttribute(kgroup, KEYCLOAK_OPENSHIFT_PROJECT_ATTRIBUTE, projectId)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	return nil
}

func DetachOpenshiftProjectFromLaboratory(laboratoryId string, projectId string) error {
	kgroup, err := keycloak.GetGroup(laboratoryId)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_KEYCLOAK)
	}

	if globalvar.IsProdEnv() {
		openshiftGroupName := OPENSHIFT_LAB_GROUP_PREFIX + laboratoryId

		//Remove lab label to namespace
		namespace, err := openshift.GetNamespace(projectId)

		if err != nil {
			return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_OPENSHIFT)
		}

		if _, ok := namespace.Labels["ceai-laboratory"]; ok {
			delete(namespace.Labels, "ceai-laboratory")

			_, err = openshift.UpdateNamespace(namespace)

			if err != nil {
				return NewErrorExternalServerError(err, ERROR_SERVER_OPENSHIFT)
			}
		}

		//Remove Rolebinding between group and project namespace
		roleBindings, err := openshift.GetNamespaceRoleBindings(projectId)

		if err != nil {
			return NewErrorExternalServerError(err, ERROR_SERVER_OPENSHIFT)
		}

		for _, roleBinding := range *roleBindings {
			for _, subject := range roleBinding.Subjects {
				if subject.Kind == "Group" && subject.Name == openshiftGroupName {
					err = openshift.DeleteRoleBinding(projectId, &roleBinding)

					if err != nil {
						return NewErrorExternalServerError(err, ERROR_SERVER_OPENSHIFT)
					}
				}
			}
		}
	}

	//Remove project to keycloak openshift_projects attribute
	err = RemoveElementFromKeycloakGroupArrayAttribute(kgroup, KEYCLOAK_OPENSHIFT_PROJECT_ATTRIBUTE, projectId)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	return nil
}

func AttachAwsAccountToLaboratory(laboratoryId string, accountId string) error {
	awsGroup, err := aws.GetGroup(GetAwsLabGroupName(laboratoryId))

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_AWS)
	}

	ssoInstanceArn := aws.GetSsoInstanceArn()
	devPermissionSetArn := aws.GetDevPermissionSetArn()

	accountAssigmnent := ssoadmin.CreateAccountAssignmentInput{
		InstanceArn:      &ssoInstanceArn,
		PermissionSetArn: &devPermissionSetArn,
		TargetId:         &accountId,
		TargetType:       types.TargetTypeAwsAccount,
		PrincipalId:      &awsGroup.ID,
		PrincipalType:    types.PrincipalTypeGroup,
	}

	err = aws.CreateAccountAssigment(&accountAssigmnent)

	if err != nil {
		return NewErrorExternalRessourceExist(err, ERROR_SERVER_AWS)
	}

	//Add account to keycloak aws_projects attribute
	kgroup, err := keycloak.GetGroup(laboratoryId)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_KEYCLOAK)
	}

	err = AddElementToKeycloakGroupArrayAttribute(kgroup, KEYCLOAK_AWS_ACCOUNT_ATTRIBUTE, accountId)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	return nil
}

func DetachAwsAccountFromLaboratory(laboratoryId string, accountId string) error {
	awsGroup, err := aws.GetGroup(GetAwsLabGroupName(laboratoryId))

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_AWS)
	}

	ssoInstanceArn := aws.GetSsoInstanceArn()
	devPermissionSetArn := aws.GetDevPermissionSetArn()

	accountAssigmnent := ssoadmin.DeleteAccountAssignmentInput{
		InstanceArn:      &ssoInstanceArn,
		PermissionSetArn: &devPermissionSetArn,
		TargetId:         &accountId,
		TargetType:       types.TargetTypeAwsAccount,
		PrincipalId:      &awsGroup.ID,
		PrincipalType:    types.PrincipalTypeGroup,
	}

	err = aws.DeleteAccountAssigment(&accountAssigmnent)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_AWS)
	}

	//Remove account from keycloak aws_projects attribute
	kgroup, err := keycloak.GetGroup(laboratoryId)

	if err != nil {
		return NewErrorExternalRessourceNotFound(err, ERROR_SERVER_KEYCLOAK)
	}

	err = RemoveElementFromKeycloakGroupArrayAttribute(kgroup, KEYCLOAK_AWS_ACCOUNT_ATTRIBUTE, accountId)

	if err != nil {
		return NewErrorExternalServerError(err, ERROR_SERVER_KEYCLOAK)
	}

	return nil
}
