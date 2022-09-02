# Builder
ARG GITHUB_PATH=gitlab.com/g6834/team17/task-service

FROM golang:1.18-alpine as builder

WORKDIR /app/${GITHUB_PATH}

RUN apk add --update make git curl
COPY Makefile Makefile
COPY . .
RUN make build

# Mail server
FROM alpine:latest as server
LABEL org.opencontainers.image.source = https://${GITHUB_PATH}
WORKDIR /root/

COPY --from=builder /app/${GITHUB_PATH}/bin/task-service .
COPY --from=builder /app/${GITHUB_PATH}/migrations/ ./migrations
COPY --from=builder /app/${GITHUB_PATH}/config.yml .

RUN chown root:root app-service

EXPOSE 3000
EXPOSE 8080

CMD ["./task-service"]
