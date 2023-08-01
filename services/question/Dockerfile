ARG SERVICE_NAME=question
ARG GO_VERSION=1.19
ARG APPCONFIG

FROM golang:${GO_VERSION}-alpine AS build_base
ARG SERVICE_NAME
LABEL stage=build_base
RUN apk update && apk add gcc libc-dev make git --no-cache ca-certificates  && \
    mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

ENV GO111MODULE=on
RUN mkdir -p /go/src/${SERVICE_NAME}
WORKDIR /go/src/${SERVICE_NAME}

COPY Makefile Makefile
COPY services/base.Makefile services/base.Makefile
COPY services/Makefile services/Makefile
COPY services/common services/common
COPY services/${SERVICE_NAME} services/${SERVICE_NAME}
COPY go.work go.work
COPY go.work.sum go.work.sum
RUN go work sync

RUN make ${SERVICE_NAME}

# ================ copy from stage build base ===========
FROM alpine:3.18
ARG SERVICE_NAME
ARG APPCONFIG

RUN apk update &&  apk add --no-cache ca-certificates git && \
    mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group


RUN mkdir -p /smartkid

COPY --from=build_base /go/src/${SERVICE_NAME}/services/${SERVICE_NAME}/${SERVICE_NAME} /smartkid/${SERVICE_NAME}
COPY --from=build_base /go/src/${SERVICE_NAME}/services/${SERVICE_NAME}/${APPCONFIG} /smartkid/app.yml

RUN chmod -R 777 /smartkid
RUN chown -R nobody:nobody /smartkid
RUN chmod -R 755 /smartkid

USER nobody:nobody

EXPOSE 3000
ENV API_DIR="/smartkid"
ENTRYPOINT ["/smartkid/question"]