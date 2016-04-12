FROM nginx:alpine

RUN mkdir -p /var/log/nginx && \
    mkdir -p /var/spool/nginx && \
    chown nobody /var/spool/nginx
COPY mime.types /mime.types
VOLUME ["/etc/nginx", "/etc/ss/certs"]