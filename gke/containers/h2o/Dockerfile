FROM alpine

RUN mkdir -p /var/log/h2o && mkdir /var/spool/h2o && chown nobody /var/spool/h2o
COPY h2o /
COPY run-h2o.sh /
VOLUME ["/etc/h2o", "/etc/ssl/certs"]
ENTRYPOINT ["/run-h2o.sh"]