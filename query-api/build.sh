#!/bin/bash

gradle clean build -x test
docker build -f Dockerfile ./ -t 192.168.1.21:8082/docker-local/query-api:0.1
docker login 192.168.1.21:8082 -u admin -p Wzjzzl2015
docker push 192.168.1.21:8082/docker-local/query-api:0.1