#!/bin/bash

# This script will create or patch all required assets to use the Admin API on a Openshift cluster

#Creating the API service account
oc create sa ceai-admin-api -n openshift
oc apply -f ceai-admin-api-admin.yaml

#Creating groups and ClusterRoleBindings
oc apply -f Admin.yaml
oc apply -f Developer.yaml
oc apply -f Admin-cluster-admin.yaml

# Removing Self Provisioning Projects rights
oc patch clusterrolebinding.rbac self-provisioners -p '{"subjects": null}'
oc patch clusterrolebinding.rbac self-provisioners -p '{ "metadata": { "annotations": { "rbac.authorization.kubernetes.io/autoupdate": "false" } } }'

# Change message on create project denial

oc patch --type=merge project.config.openshift.io cluster -p '{"spec":{"projectRequestMessage":"Contactez un administrateur du CEAI pour créer un projet."}}'