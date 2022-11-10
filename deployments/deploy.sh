#!/bin/sh

docker build -t remi-app .
docker tag remi-app quangngoc430/remi:v1.0.2
docker push quangngoc430/remi:v1.0.2
cd deployments/
kubectl apply -f .
kubectl rollout restart deployment remi