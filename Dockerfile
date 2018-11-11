FROM debian

RUN  apt-get update && apt-get install -y ca-certificates

COPY kube-sqs-autoscaler /

CMD ["/kube-sqs-autoscaler"]
