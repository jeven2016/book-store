FROM docker-public.home.com/library/nginx:latest

ENV LANG C.UTF-8
ENV LC_ALL C.UTF-8

COPY build /static
COPY config/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]