#!/bin/bash

res=$(kubectl get pod --namespace=$1 -l svc_id=$2 --no-headers | awk '{print $1}')
echo -n $res
