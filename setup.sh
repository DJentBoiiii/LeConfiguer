#!/bin/bash


kubectl apply -f kubernetes/secrets.yaml
kubectl apply -f kubernetes/deployment.yaml
kubectl apply -f kubernetes/services.yaml