FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/k8s-dns-exposer \
    USER_UID=1001 \
    USER_NAME=k8s-dns-exposer

# install operator binary
COPY build/_output/bin/k8s-dns-exposer ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
