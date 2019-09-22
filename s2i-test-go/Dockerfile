FROM docker.io/golang:1.12

LABEL maintainer="Hongkai Liu <hongkailiu@users.noreply.github.com>"

ENV BUILDER_VERSION 1.0.0

LABEL io.k8s.description="Platform for building test-go" \
      io.k8s.display-name="builder 1.0.0" \
      io.openshift.expose-services="8080:http" \
      io.openshift.s2i.destination="/tmp" \
      io.openshift.tags="builder,html,test-go"

LABEL io.openshift.s2i.scripts-url=image:///usr/local/s2i
COPY ./s2i/bin/ /usr/local/s2i

USER 1001

WORKDIR /go/src/github.com/hongkailiu/test-go/

EXPOSE 8080

CMD ["/usr/local/s2i/usage"]
