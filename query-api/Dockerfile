FROM docker-public.home.com/library/openjdk:21-jdk-slim
#FROM docker-public.home.com/library/openjdk:21-jdk-oracle

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8


ADD build/libs/query-api-1.0.jar /usr/local/bin/query-api.jar

RUN chmod a+x /usr/local/bin/query-api.jar

ENTRYPOINT ["java","-jar","/usr/local/bin/query-api.jar"]