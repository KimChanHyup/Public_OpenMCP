#/bin/bash
NS=openmcp
CLUSTER=cluster2
NAME=$(kubectl get pod -n $NS --context $CLUSTER | grep -E 'cluster-metric-collector' | awk '{print $1}')

#echo "Exec Into '"$NAME"'"

#kubectl exec -it $NAME -n $NS /bin/sh
kubectl logs -f -n $NS $NAME --context $CLUSTER

