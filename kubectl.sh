#!/bin/bash
namespace=$1
name=$2
shift
shift

#set -x
#kubectl exec  --namespace=$namespace $name $@

echo $@ |  kubectl exec -i --namespace=$namespace $name -- /bin/sh
