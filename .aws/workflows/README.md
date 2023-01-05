# CI-CD AWS

# AWS CodeBuild

## Projets Build

### 1. Pull Request: "ceai-admin-api-pull-request-build-prj"
  Ce projet a été créé pour exécuter le build lors d'une création ou un mis à jour d'un Pull Request
  Les tâches réalisées par ce projet sont détaillées à continuation:

- Prend le code du dépôt GitHub, branch du Pull Request.
    - Il travaille avec le fichier de buildspec dans le code de source: [ci-pr.yml](ci-pr.yml)
- Build de l'application (go command)
- Scan de l'application avec l'outil de scan pour les projets go: govulncheck
- Dépose le résultat du scan dans un fichier nommé "**govulncheck-scan-results**".json**      
    - Il utilise le AWS S3 bucket: "**admin-api-pipeline-bucket**"
    - folder: "**pull-request-build**"

### 2. Dev: "ceai-admin-api-dev-build-prj"
  Ce projet a été créé pour exécuter le pipeline de déploiement: "**ceai-admin-api-dev-ci-cd-pipeline**"
  Les tâches réalisées par ce projet sont détaillées à continuation:

  - Prend le code du dépôt GitHub, branch **dev**
    - Il travaille avec le fichier de buildspec dans la source de code: [ci-cd-dev.yml](ci-cd-dev.yml)
  - Build de l'application (go command)
  - Scan de l'application avec l'outil de scan pour les projets go: govulncheck
  - Dépose le résultat du scan dans un fichier nommé "**govulncheck-scan-results.json**"
    - Il utilise le AWS S3 bucket: **admin-api-pipeline-bucket**
      - folder: **dev-build**
  - Build l'image docker de l'application
  - Tag l'image avec le build number du AWS CodeBuild
  - Push l'image docker dans AWS ECR: "**ceai-admin-api-codebuild-ecr**"
  - Prend le résultat de l'scan de l'image docker de AWS ECR et l'écrit dans un fichier: "**scan-results-$IMAGE_TAG.out**"
  - Produit le fichier "**imagedefinitions.json**" avec les informations de l'image à déployer dans AWS ECS

# AWS CodePipeline

## Pipeline: "ceai-admin-api-dev-ci-cd-pipeline"
  Ce pipeline déploie une nouvelle version de l'application dans AWS ECS. 
  Les tâches réalisées par ce pipeline sont détaillées à continuation:

  - Prends le code du dépôt GitHub, branch **dev**
  - Appel le projet de build: "**ceai-admin-api-dev-build-prj**" pour obtenir l'image Docker de l'application
  - Déploie l'image spécifiée dans le fichier "**imagedefinitions.json**" (dans le bucket S3), dans AWS ECS
    - AWS ECS: 
      - cluster: **ceai-admin-api-dev-ecs**
      - service: **ceai-admin-api-dev-app-service**
