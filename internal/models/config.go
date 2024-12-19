package models

type Config struct {
	Name          *string          `json:"name,omitempty"`
	Description   *string          `json:"description,omitempty"`
	OpenAPIPath   *string          `json:"openApiPath,omitempty"`
	Port          *string          `json:"port,omitempty"`
	SwaggerUIPath *string          `json:"swaggerUIPath,omitempty"`
	GatewaySecret *string          `json:"gatewaySecret,omitempty"`
	Realm         *ConfigRealm     `json:"realm,omitempty"`
	Providers     *ConfigProviders `json:"Providers,omitempty"`
}

type ConfigRealm struct {
	Url                  *string `json:"url,omitempty"`
	Name                 *string `json:"name,omitempty"`
	ServiceAccountId     *string `json:"serviceAccountId,omitempty"`
	ServiceAccountSecret *string `json:"serviceAccountSecret,omitempty"`
}

type ConfigProviders struct {
	AwsLandingZone *ConfigAwsLandingZoneProvider `json:"awsLandingZone,omitempty"`
	Openshift      *ConfigOpenshiftProvider      `json:"openshift,omitempty"`
	DeployKF       *ConfigDeployKFProvider       `json:"deployKF,omitempty"`
	Github         *ConfigGithubProvider         `json:"github,omitempty"`
}

type ConfigProvider struct {
	Name        *string `json:"name,omitempty"`
	Type        *string `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
	ClientUIID  *string `json:"clientUIID,omitempty"`
}

type ConfigAwsLandingZoneProvider struct {
	ConfigProvider `yaml:",inline"`
	Config         *ConfigAwsLandingZoneConfig        `json:"config,omitempty"`
	AccountRoles   *[]ConfigAwsLandingZoneAccountRole `json:"accountRoles,omitempty"`
}

type ConfigAwsLandingZoneConfig struct {
	ScimEndpoint      *string `json:"scimEndpoint,omitempty"`
	ScimToken         *string `json:"scimToken,omitempty"`
	AwsAccessKey      *string `json:"awsAccessKey,omitempty"`
	AwsSecret         *string `json:"awsSecret,omitempty"`
	AwsSsoInstanceArn *string `json:"awsSsoInstanceARN,omitempty"`
}

type ConfigAwsLandingZoneAccountRole struct {
	Name             *string `json:"name,omitempty"`
	PermissionSetArn *string `json:"permissionSetArn,omitempty"`
}

type ConfigOpenshiftProvider struct {
	ConfigProvider `yaml:",inline"`
	Config         *ConfigOpenshiftConfig        `json:"config,omitempty"`
	ProjectRoles   *[]ConfigOpenshiftProjectRole `json:"projectRoles,omitempty"`
}

type ConfigOpenshiftConfig struct {
	KubeConfigPath *string `json:"kubeConfigPath,omitempty"`
}

type ConfigOpenshiftProjectRole struct {
	Name        *string `json:"name,omitempty"`
	ClusterRole *string `json:"clusterRole,omitempty"`
}

type ConfigDeployKFProvider struct {
	ConfigProvider `yaml:",inline"`
	Config         *ConfigDeployKFConfig    `json:"config,omitempty"`
	Profiles       *[]ConfigDeployKFProfile `json:"profiles,omitempty"`
}

type ConfigDeployKFConfig struct {
	GitRepo *string `json:"gitRepo,omitempty"`
}

type ConfigDeployKFProfile struct {
	Name            *string `json:"name,omitempty"`
	Role            *string `json:"role,omitempty"`
	NotebooksAccess *string `json:"notebookAccess,omitempty"`
}

type ConfigGithubProvider struct {
	ConfigProvider  `yaml:",inline"`
	Config          *ConfigGithubConfig           `json:"config,omitempty"`
	RepositoryRoles *[]ConfigGithubRepositoryRole `json:"repositoryRoles,omitempty"`
}

type ConfigGithubConfig struct {
	OrgName    *string `json:"orgName,omitempty"`
	Token      *string `json:"token,omitempty"`
	PrivateKey *string `json:"privateKey,omitempty"`
	TeamID     *string `json:"teamID,omitempty"`
}

type ConfigGithubRepositoryRole struct {
	Name *string `json:"name,omitempty"`
	Role *string `json:"role,omitempty"`
}
