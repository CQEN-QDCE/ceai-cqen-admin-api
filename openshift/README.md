# Ressources Openshift requises par l'API d'administration

## Installation

 **Exécuter seulement sur un cluster neuf, sous peine de supprimer tous les usagers de vos groupes!**

Avec l'outil cli oc (https://docs.openshift.com/container-platform/4.7/cli_reference/openshift_cli/getting-started-cli.html) préalablement installé:

```bash
oc login [avec un usager ayant droit cluster-admin]
./install.sh
```

Le script va créer les ressources suivantes sur le cluster:

## Groupes

### Admin

Le groupe des administrateurs associé au role cluster-admin. 

### Developer

Les développeurs n'ont pas le droit de créer des projects. Ceux-ci doivent être créés par un administrateur. 

## Service Account

L'API nécéssite l'accès vers Openshift via un ServiceAccount. Le script crée un compte nommé ceai-admin-api dans le project Openshift. Le fichier kubeconfig peut alors être complété avec les informations d'installation du cluster ainsi que le token généré pour le ServiceAccount.

## Création d'un laboratoire

En attendant la fonctionalité de création de laboratoire, voici la manière de créer un laboratoire et d'y associer des projects et des utilisateurs.

Créer un groupe contenant les utilisateur du laboratoire

```
oc adm groups new [Lab_NomDuLab]
oc adm groups add-users [Lab_NomDuLab] [usager1 usager2]
```

Créer les projects nécéssaires au laboratoire et les associer au groupe

```
oc new project [NomDuLab_Project1]
oc adm policy add-role-to-group [admin|edit|view] [Lab_NomDuLab] --rolebinding-name="[nom_role_binding]" -n [NomDuLab_Project1]
```
