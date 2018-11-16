FROM nginx:1.13.12-alpine

RUN apk --no-cache add nginx-mod-http-set-misc nginx-mod-http-lua
RUN mkdir -p /run/nginx

COPY conf.d/edge.conf /etc/nginx/nginx.conf
COPY conf.d/mime.types /etc/nginx/mime.types
COPY favicon.ico /var/lib/nginx/html/favicon.ico

CMD ["nginx"]