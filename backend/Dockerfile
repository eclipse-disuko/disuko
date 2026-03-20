FROM golang:1.26.1-alpine AS builder
WORKDIR /go/server/
ENV CGO_ENABLED=0
COPY ./go.mod .
RUN go mod download

COPY . .
RUN go build -o /go/server/app


FROM alpine:3.22
ENV USERNAME=disukoUser
WORKDIR /go/server/

COPY --from=builder  /go/server/app /go/server/
COPY --from=builder  /go/server/conf/dbseeds/ /go/server/conf/dbseeds/
COPY --from=builder  /go/server/conf/config.yml /go/server/conf/config.yml
COPY --from=builder  /go/server/resources/ /go/server/resources/
COPY ./certs/root/ /usr/local/share/ca-certificates/
RUN apk update
RUN apk add terminus-font ttf-inconsolata ttf-dejavu font-noto font-noto-cjk ttf-font-awesome font-noto-extra

RUN apk add --no-cache \
    --repository http://dl-cdn.alpinelinux.org/alpine/v3.22/main \
    ca-certificates
RUN update-ca-certificates
RUN apk add chromium
RUN apk del ffmpeg
RUN apk del libvpx
RUN apk add poppler-utils
RUN addgroup -S 1
RUN adduser -S $USERNAME -G 1 -u 1000
RUN echo "$USERNAME:$USERNAME" | chpasswd
RUN chown -R $USERNAME:1 /tmp
RUN chown -R $USERNAME:1 /go/server/
USER $USERNAME

CMD ["/go/server/app"]

