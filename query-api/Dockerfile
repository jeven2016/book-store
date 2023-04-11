FROM 192.168.1.21:8082/docker-group/openjdk:20-ea-17-jdk-oracle

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8


ADD build/libs/query-api-1.0.jar /usr/local/bin/query-api.jar

RUN chmod a+x /usr/local/bin/query-api.jar

ENTRYPOINT ["java","-jar","/usr/local/bin/query-api.jar"]