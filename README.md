# API d'administration de l'infrastructure du centre d'expertise appliquée en innovation du CQEN

## Installation

 * sample.env

## Swagger UI

 * go get github.com/rakyll/statik
 * cd api
 * ./genswagger.sh

## test.http

* ext install humao.rest-client

## Installation environnement Go

### Fedora Linux

Installer les dépendances Golang:
```
sudo dnf install golang
```

Ajouter dans $HOME/.bashrc:

```
export GOPATH="$HOME/go"
export PATH=$PATH:$GOPATH/bin
```

### MacOS

```
#TODO
```