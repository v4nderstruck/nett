#!/bin/bash

docker exec -it clab-s2s-PC1_VN1001 ip link set eth1 up
docker exec -it clab-s2s-PC1_VN1001 ip addr add 10.10.1.100/24 dev eth1
docker exec -it clab-s2s-PC1_VN1001 ip route del default
docker exec -it clab-s2s-PC1_VN1001 ip route add default via 10.10.1.1 dev eth1

docker exec -it clab-s2s-PC2_VN1002 ip link set eth1 up
docker exec -it clab-s2s-PC2_VN1002 ip addr add 10.10.2.100/24 dev eth1
docker exec -it clab-s2s-PC2_VN1002 ip route del default
docker exec -it clab-s2s-PC2_VN1002 ip route add default via 10.10.2.1 dev eth1

docker exec -it clab-s2s-PC3_VN2001 ip link set eth1 up
docker exec -it clab-s2s-PC3_VN2001 ip addr add 10.20.1.100/24 dev eth1
docker exec -it clab-s2s-PC3_VN2001 ip route del default
docker exec -it clab-s2s-PC3_VN2001 ip route add default via 10.20.1.1 dev eth1

docker exec -it clab-s2s-PC4_VN2002 ip link set eth1 up
docker exec -it clab-s2s-PC4_VN2002 ip addr add 10.20.2.100/24 dev eth1
docker exec -it clab-s2s-PC4_VN2002 ip route del default
docker exec -it clab-s2s-PC4_VN2002 ip route add default via 10.20.2.1 dev eth1
