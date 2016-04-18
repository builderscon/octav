FROM busybox:latest
RUN mkdir -p /etc/ssl/certs 
COPY ca-certificates.crt /etc/ssl/certs/
COPY acmebot /
ENTRYPOINT ["/acmebot"]

