FROM alpine:3.4

RUN apk add --update ca-certificates tzdata && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*
ADD apiserver /apiserver

CMD ["/apiserver"]