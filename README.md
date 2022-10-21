# API d'administration de l'infrastructure du centre d'expertise appliquée en innovation du CQEN

Le centre d'expertise appliquée en innovation du CQEN offre des services de laboratoires d'expérimentation sur ses plateformes AWS SEA et Openshift. Il fournit une authentification unique vers ses plateformes via un serveur d'authentification Keycloak.

Le CEAI a conçu une API d'administration unifiant la gestion des usagers et des ressources allouées aux laboratoires sur les trois produits mentionnés ci-haut : AWS SEA, Openshift et Keycloak.

La création de cet API vise les buts suivants :

 * Propager les droits d'accès aux ressources de laboratoires via une ressource unique.
 * Permettre la création de consoles d'administration unifiées.
 * Automatiser la gestion de l'infrastructure via des scripts et des appels HTTP.
 * Produire une preuve de concept d'un API défini par sa documentation (API first) dans le langage Go et la spécification OpenAPI 3.0.

Le dépôt contient le serveur d'API ainsi qu'une console en ligne de commande (CLI) pour exploiter celui-ci.

# Serveur API

L'API s'appuis sur un Realm Keycloak lequel il exploite via l'API Keycloak. L'API nécessite des accès API au point d'accès du service SSO SCIM ainsi qu'un accès API IAM au compte AWS Master de l'environnement SEA. Finalement, un accès API à Openshift est aussi nécessaire.

## Keycloak

TODO

## AWS SSO SCIM

TODO

## AWS Master

TODO

## Openshift

TODO


## Sécurité

L'API ne contient pas de mécanisme d'authentification des usagers. Cela doit être pris en charge par un API Gateway ou un Reverse Proxy. Celui-ci doit prendre en charge l'authentification et l'identification de l'usager. Il doit passer le nom d'usager et son rôle à l'API via les entêtes _X-CEAI-Username_ et _X-CEAI-UserRoles_ .

Pour s'assurer que l'API ne soit pas accédée directement un "Gateway Secret" peut être défini dans les variables d'environnement. L'API validera que celui-ci est passé dans l'entête _X-CEAI-Gateway-Secret_ avant de traiter la requête.

Un exemple de déploiement derrière le produit AWS API Gateway est disponible dans le dépôt suivant : https://github.com/CQEN-QDCE/ceai-cqen-deployments

## Variables d'environnement

| Nom                           | Description                                                   |
| ----------------------------  | ------------------------------------------------------------- |
| `PORT`                        | Port sur lequel exposer l'API                                 |
| `OPENAPI_PATH`                | Chemin où est déposé la spécification OpenApi
| `SWAGGER_UI_PATH`             | Route où exposer la documentation SwaggerUI, laisser vide pour ne pas exposer.
| `GATEWAY_SECRET`              | Secret qu'un API Gateway/Reverse Proxy doit fournir dans l'entête X-CEAI-Gateway-Secret pour s'authentifier à l'API comme client valide.
| `SCIM_ENDPOINT`               | Url d'accès à l'API SCIM de AWS SSO.
| `SCIM_TOKEN`                  | Jeton d'accès à l'API SCIM de AWS SSO.
| `AWS_ACCESS_KEY`              | Clé d'accès au compte IAM de Service 
| `AWS_SECRET`                  | Secret du compte IAM de Service
| `AWS_SSO_INSTANCE_ARN`        | ARN de l'instance AWS SSO
| `AWS_ADMIN_PERMISSION_SET_ARN`| ARN du Permission Set lié au profil Administrateur dans AWS
| `AWS_DEV_PERMISSION_SET_ARN`  | ARN du Permission Set lié au profil Developpeur dans AWS
| `KEYCLOAK_URL`                | Url du serveur Keycloak
| `KEYCLOAK_REALM`              | Nom du realm Keycloak
| `KEYCLOAK_CLIENT_ID`          | Identifiant du client Keycloak dédié à l'API
| `KEYCLOAK_CLIENT_SECRET`      | Secret du client Keycloak dédié à l'API
| `KUBECONFIG_PATH`             | Chemin vers le fichier kubeconfig lié au Cluster Openshift

## Développement

### Architecture

TODO

### Tester

```
cp sample.env .env

# Renseigner les variables d'environnement dans le fichier .env

go run ./cmd/server
```

### Swagger UI

La documentation SwaggerUI peut être générée et rendue disponible par l'API. La variable d'environnement _SWAGGER_UI_PATH_ doit être définie pour rendre accessible SwaggerUI à l'url voulue.
 
La documentation SwaggerUI doit être générée et compilée à partir du fichier api/openapi-v1.yaml. En cas de modification à la définition, SwaggerUi doit être régénéré et recompilé :
 
```
go get github.com/rakyll/statik
cd third_party
/compile_spec.sh
```

### test.http

Des tests unitaires sont pour la plupart des routes sont disponibles dans le fichier test.http. Pour utiliser celui-ci, vous avez besoin de l'extension VS Code [Rest-Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)

# Console CLI

## Compiler et installer l'exécutable du CLI

Linux - MacOS - WSL
```
go build -o ./ceai ./cmd/cli

#Voir les répertoires pris en charge par votre système où copier l'exécutable
echo $PATH

cp ceai [/votre/repertoire/de/bin/prefere/] 

ceai --help
```