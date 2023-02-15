# Exemples d'utilisation du client en ligne de commandes (CLI)

## Se connecter au serveur d'API

```bash
# Authentification en spécifiant l'identifiant
ceai login https://admin-api.ceai.cqen.ca -u [votre_identifiant_ceai]

# Véfifier la connexion
ceai whoami
```

## Consulter une ressource
```bash
# Consulter tous les usagers
ceai get users

# Consulter un usager
ceai get user [courriel]

# Consulter des laboratoires
ceai get labs
ceai get lab [identifiant_du_lab]

# Consulter des projets Openshift
ceai get projects
ceai get project [Identifiant_du_project]

# Consulter des comptes AWS
ceai get accounts
ceai get account [Numéro_du_compte]

# Toutes les commandes get peuvent être obtenues en format json ou yaml
ceai get user [courriel] -o json
ceai get user [courriel] -o yaml

```

## Créer un usager

```bash
# Création interactive 
ceai create user -i

# Création inline format yaml (Format recommandé pour gérer les caractères spéciaux)
ceai create user --yaml "$(cat <<'EOF'
- email: [courriel]
  firstname: [prenom]
  lastname: [nom]
  organisation: [organisation]
  infrarole: [Developer|Admin]
  disabled: false
EOF
)"

# Création inline format json
ceai create user --json "$(cat <<'EOF'
[
  {
    "email": "[courriel]",
    "firstname": "[courriel]",
    "lastname": "[prenom]",
    "organisation": "[organisation]",
    "infrarole": "[Developer|Admin]",
    "disabled": false
  }
]
EOF
)"

# Le paramètre d'entrée étant un tableau, plusieurs usagers peuvent être créés à la fois
ceai create user --yaml "$(cat <<'EOF'
- email: [courriel]
  firstname: [prenom]
  lastname: [nom]
  organisation: [organisation]
  infrarole: [Developer|Admin]
  disabled: false
- email: [courriel]
  firstname: [prenom]
  lastname: [nom]
  organisation: [organisation]
  infrarole: [Developer|Admin]
  disabled: false
EOF
)"

# À partir d'un fichier (Fichier exemple dans le répertoire)
ceai create user --yamlfile ./usager.yaml
ceai create user --jsonfile ./usager.json
```

## Modifier un usager
```bash
# On modifie une ou plusieurs informations via les *flags* de modification
ceai update user [courriel] --Firstname [NouveauPrenom] --Lastname [NouveauNom] --Organisation [NouvelleOrganisation] --Infrarole [Developer|Admin] --Disabled [true|false]
```


## Réinitialiser le mot de passe ou le jeton OTP
```bash
# Le mot de passe, le jeton OTP ou les deux secrets peuvent être réinitialisés
ceai reset [password|otp|all] [courriel]
```

## Renvoyer le courriel de création de compte
```bash
# Le courriel de création de compte expire après 24h, pour en renvoyer un nouveau
ceai send required-actions [courriel]
```

## Créer un laboratoire
```bash
# Création interactive 
ceai create laboratory -i

# Création format yaml (Format recommandé pour gérer les caractères spéciaux)
ceai create laboratory --yaml "$(cat <<'EOF'
- id: [Identifiant]
  displayname: [Nom complet du laboratoire]
  description: [Description]
  type: [projet|experimentation]
  gitrepo: [Url d'un dépot Git]
EOF
)"

# Création format json
ceai create laboratory --json "$(cat <<'EOF'
[
  {
    "id": "[Identifiant]",
    "displayname": "[Nom complet du laboratoire]",
    "description": "[Description]",
    "type": "[projet|experimentation]",
    "gitrepo": "[Url d'un dépot Git]"
  }
]
EOF
)"

# À partir d'un fichier (Fichier exemple dans le répertoire)
ceai create user --yamlfile ./lab.yaml
ceai create user --jsonfile ./lab.json
```

## Créer un *project* Openshift
```bash
# Création interactive 
ceai create project -i

# Création format yaml (Format recommandé pour gérer les caractères spéciaux)
ceai create project --yaml "$(cat <<'EOF'
- id: [Identifiant]
  displayname: [Nom complet du project]
  description: [Description]
  idLab: [Identifiant du laboratoire]
EOF
)"

# Création format json
ceai create project --json "$(cat <<'EOF'
[
  {
    "id": "[Identifiant]",
    "displayname": "[Nom complet du project]",
    "description": "[Description]",
    "idLab": "[Identifiant du laboratoire]"
  }
]
EOF
)"

# À partir d'un fichier (Fichier exemple dans le répertoire)
ceai create user --yamlfile ./project.yaml
ceai create user --jsonfile ./project.json
```

## Associer des ressources à un laboratoire
```bash
# Associer/ Désassocier un usager d'un laboratoire
ceai add lab-user [Id laboratoire] [courriel1 courriel2 courrielX...]
ceai remove lab-user [Id laboratoire] [courriel1 courriel2 courrielX...]

# Associer/ Désassocier un project Openshift d'un laboratoire
ceai add lab-project [Id laboratoire] [Id project]
ceai remove lab-project [Id laboratoire] [Id project]

# Associer/ Désassocier un compte AWS d'un laboratoire
ceai add lab-account [Id laboratoire] [Numéro de compte AWS]
ceai remove lab-account [Id laboratoire] [Numéro de compte AWS]
```
