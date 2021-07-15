# Ressources Openshift requises par l'API d'administration

## Installation

Avec l'outil cli oc (https://docs.openshift.com/container-platform/4.7/cli_reference/openshift_cli/getting-started-cli.html) préalablement installé:

```bash
oc login [avec un usager ayant droit cluster-admin]
./install.sh
```

Le script va créer les ressources suivantes sur le serveur:

## Groupes

### Admin

Le groupe des administrateurs associé au role cluster-admin. 

### Developer

Les développeurs n'ont pas le droit de créer des projects. Ceux-ci doivent être créé par un administrateur. 

## Service Account

L'API nécéssite l'accès vers Openshift via un ServiceAccount. le script crée un compte nommé ceai-admin-api dans le project Openshift. Le fichier kubeconfig peut alors être complété avec les informations d'installation du cluster ainsi que le token généré pour le ServiceAccount.