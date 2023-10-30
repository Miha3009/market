kubectl create namespace market
kubectl create configmap -n market products-init --from-file=products-init.sql
kubectl create configmap -n market products-test-values --from-file=products-test-values.sql
kubectl create configmap -n market inventory-init --from-file=inventory-init.sql
kubectl create configmap -n market inventory-test-values --from-file=inventory-test-values.sql
kubectl create configmap -n market graphql-config --from-file=graphql-config.yaml
kubectl create configmap -n market products-config --from-file=products-config.yaml
kubectl create configmap -n market orders-config --from-file=orders-config.yaml
kubectl create configmap -n market inventory-config --from-file=inventory-config.yaml
kubectl create configmap -n market notifier-config --from-file=notifier-config.yaml
kubectl apply -f init.yaml
