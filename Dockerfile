FROM alpine:3.5

RUN apk add --no-cache --update ca-certificates

COPY bin/k8s-job-cleaner /k8s-job-cleaner

ENTRYPOINT ["/k8s-job-cleaner"]
