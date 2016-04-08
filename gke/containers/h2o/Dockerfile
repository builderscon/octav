FROM alpine

RUN mkdir -p /var/log/h2o && mkdir /var/spool/h2o && chown nobody /var/spool/h2o
COPY h2o /
VOLUME ["/etc/h2o", "/etc/ssl/certs"]
ENTRYPOINT ["/h2o"]