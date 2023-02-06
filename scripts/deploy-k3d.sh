#!/usr/bin/env bash

EXAMPLESDIR="examples"

createCluster() {
    k3d registry create customer-registry --port 5050
    k3d cluster create customer --port '8081:80@loadbalancer' --registry-use k3d-customer-registry:5050 --registry-config ${EXAMPLESDIR}/registries.yaml
    echo "Before pushing images to customer-registry"
    echo "Add to /etc/hosts"
    echo "127.0.0.1 k3d-customer-registry"
}

deployCustomersNS() {
    kubectl create ns customer-a
    kubectl create ns customer-b
}

deployCustomers() {
    kubectl apply -f ${EXAMPLESDIR}/k8s-customer-a.yaml
    kubectl apply -f ${EXAMPLESDIR}/k8s-customer-b.yaml
}

deleteCustomers() { 
    kubectl delete -f ${EXAMPLESDIR}/k8s-customer-a.yaml
    kubectl delete -f ${EXAMPLESDIR}/k8s-customer-b.yaml
    kubectl delete ns customer-a
    kubectl delete ns customer-b
}

deleteCluster() {
    k3d cluster delete customer
    k3d registry delete customer-registry
}

testConnection() {
    echo "Connecting to Customer A"
    curl -H "Host: customer-a.localhost" http://127.0.0.1:8081/v1/greeting
    echo "Connecting to Customer B"
    curl -H "Host: customer-b.localhost" http://127.0.0.1:8081/v1/greeting
}

listCustomers() {
    kubectl get pods -n customer-a
    kubectl get pods -n customer-b
}

while getopts "cdlrt" option; do
case ${option} in
c ) 
echo "Creating K3D cluster customer"
createCluster
echo "Deploy Customer applications A and B"
deployCustomersNS
;;
d ) 
echo "Deleting deployments and cluster"
deleteCustomers
deleteCluster
;;
l ) 
listCustomers
;;
r ) 
echo "Deployment of Customer applications A and B"
deployCustomers
;;

t )
echo "Testing endpoints via ingress"
testConnection
;;
* ) 
echo "You have to use: [-c] to create or [-d] to delete everything or [-r] to re-deploy all manifests again"
;;
esac
done