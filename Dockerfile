FROM golang:alpine as build
MAINTAINER Tiberiu Craciun <tibi@happysoft.ro>

WORKDIR /go/src/project
COPY ./backup_agent.go /go/src/project
RUN go build -o /usr/local/bin/backup_agent


FROM alpine
COPY --from=build /usr/local/bin/backup_agent /usr/local/bin/backup_agent

EXPOSE 9191

ENTRYPOINT ["/usr/local/bin/backup_agent"]
CMD ["-h"]
