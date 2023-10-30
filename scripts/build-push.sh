docker build -t miha3009/graphql:latest ../graphql
docker build -t miha3009/inventory:latest ../inventory
docker build -t miha3009/notifier:latest ../notifier
docker build -t miha3009/orders:latest ../orders
docker build -t miha3009/products:latest ../products

docker push miha3009/graphql:latest
docker push miha3009/inventory:latest
docker push miha3009/notifier:latest
docker push miha3009/orders:latest
docker push miha3009/products:latest
