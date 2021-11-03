package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	scim "github.com/CQEN-QDCE/aws-sso-scim-goclient"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/aws"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/keycloak"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/internal/openshift"
	"github.com/CQEN-QDCE/ceai-cqen-admin-api/pkg/apifirst"
	"github.com/Nerzal/gocloak/v8"
	"github.com/gorilla/mux"
	userv1 "github.com/openshift/api/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type LaboratoryHandlersInterface interface {

	// (GET /laboratory)
	GetLaboratories(response *apifirst.Response, request *http.Request) error

	// (GET /laboratory/{laboratoryid})
	GetLaboratoryFromId(response *apifirst.Response, r *http.Request) error

	// (POST /laboratory)
	CreateLaboratory(response *apifirst.Response, r *http.Request) error
}

// AWSAccount defines model for AWSAccount.
type AWSAccount struct {
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
	Id    *string `json:"number,omitempty"`
}

// Laboratory defines model for Laboratory.
type Laboratory struct {
	Id          string  `json:"id"`
	Description string  `json:"description"`
	Displayname string  `json:"displayname"`
	Type        string  `json:"type"`
	Gitrepo     *string `json:"gitrepo,omitempty"`
}

// LaboratoryWithResources defines model for LaboratoryWithResources.
type LaboratoryWithResources struct {
	// Embedded struct due to allOf(#/components/schemas/LaboratoryWithUsers)
	*LaboratoryWithUsers `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	AWSAccounts       *[]AWSAccount       `json:"AWSAccounts,omitempty"`
	Openshiftprojects *[]OpenshiftProject `json:"openshiftprojects,omitempty"`
}

// LaboratoryWithUsers defines model for LaboratoryWithUsers.
type LaboratoryWithUsers struct {
	// Embedded struct due to allOf(#/components/schemas/Laboratory)
	*Laboratory `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Users *[]string `json:"users,omitempty"`
}

// OpenshiftProject defines model for OpenshiftProject.
type OpenshiftProject struct {
	Description string `json:"description"`
	Displayname string `json:"displayname"`
	Id          string `json:"id"`
}

// OpenshiftProjectWithLab defines model for OpenshiftProjectWithLab.
type OpenshiftProjectWithLab struct {
	// Embedded struct due to allOf(#/components/schemas/OpenshiftProject)
	OpenshiftProject `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	IdLab *string `json:"idLab,omitempty"`
}

func MapLaboratory(kgroup gocloak.Group) (*Laboratory, error) {
	var lab Laboratory

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

func MapLaboratoryWithUsers(kgroup gocloak.Group) (*LaboratoryWithUsers, error) {
	var lab LaboratoryWithUsers

	lab.Laboratory, _ = MapLaboratory(kgroup)

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

func MapLaboratoryWithResources(kgroup gocloak.Group) (*LaboratoryWithResources, error) {
	var lab LaboratoryWithResources

	lab.LaboratoryWithUsers, _ = MapLaboratoryWithUsers(kgroup)

	attributes := *kgroup.Attributes

	var openshiftprojects []OpenshiftProject

	if attributes["openshift_projects"] != nil {
		for _, projectName := range attributes["openshift_projects"] {
			var project OpenshiftProject

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

	var awsAccounts []AWSAccount

	if attributes["aws_accounts"] != nil {
		for _, accountId := range attributes["aws_accounts"] {
			var account AWSAccount

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

func CreateGroupKeycloak(plab *Laboratory) error {
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

func CreateGroupAws(plab *Laboratory) error {
	group := scim.NewGroup(AWS_LAB_GROUP_PREFIX + plab.Id)

	_, err := aws.CreateGroup(group)

	return err
}

func CreateGroupOpenshift(plab *Laboratory) error {
	group := userv1.Group{
		ObjectMeta: metav1.ObjectMeta{
			Name: OPENSHIFT_LAB_GROUP_PREFIX + plab.Id,
		},
	}

	_, err := openshift.CreateGroup(&group)

	return err
}

func (s ServerHandlers) GetLaboratories(response *apifirst.Response, request *http.Request) error {
	labGroups, err := keycloak.GetGroups(gocloak.StringP(KEYCLOAK_LAB_TOP_GROUP))

	if err != nil {
		response.SetStatus(http.StatusInternalServerError)
		log.Println(err)
		return err
	}

	labsList := make([]*LaboratoryWithResources, 0, len(*labGroups))

	if labGroups != nil {

		for _, group := range *labGroups {
			lab, err := MapLaboratoryWithResources(group)

			if err == nil {
				labsList = append(labsList, lab)
			}
			//TODO Log error?
		}
	}

	response.SetStatus(http.StatusOK)
	response.SetBody(labsList)

	return nil
}

func (s ServerHandlers) GetLaboratoryFromId(response *apifirst.Response, request *http.Request) error {
	params := mux.Vars(request)
	laboratoryid := params["laboratoryid"]

	labGroup, err := keycloak.GetGroup(laboratoryid)

	if err != nil {
		response.SetStatus(http.StatusNotFound)
		log.Println(err)
		return err
	}

	lab, err := MapLaboratoryWithResources(*labGroup)

	if err == nil {
		response.SetStatus(http.StatusOK)
		response.SetBody(lab)
	}

	return err
}

func (s ServerHandlers) CreateLaboratory(response *apifirst.Response, request *http.Request) error {
	plab := Laboratory{}
	if err := json.NewDecoder(request.Body).Decode(&plab); err != nil {
		response.SetStatus(http.StatusBadRequest)
		log.Println(err)
		return err
	}

	var kerr, oerr, aerr error

	kfunc := func() {
		kerr = CreateGroupKeycloak(&plab)
	}

	ofunc := func() {
		oerr = CreateGroupOpenshift(&plab)
	}

	afunc := func() {
		aerr = CreateGroupAws(&plab)
	}

	Parallelize(kfunc, ofunc, afunc)

	//TODO Error map
	if kerr != nil {
		log.Println("Keycloak error: " + kerr.Error())
		response.SetStatus(http.StatusConflict)
		return kerr
	}

	if oerr != nil {
		log.Println("Openshift error: " + oerr.Error())
		response.SetStatus(http.StatusConflict)
		return oerr
	}

	if aerr != nil {
		log.Println("AWS error: " + aerr.Error())
		response.SetStatus(http.StatusConflict)
		return aerr
	}

	response.SetStatus(http.StatusCreated)

	return nil
}
