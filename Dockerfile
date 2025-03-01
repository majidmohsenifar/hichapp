FROM golang:1.22.4-bullseye AS BUILD
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o main .
FROM base as builder

FROM debian:bullseye AS FINAL
RUN apt update && apt install -y ca-certificates 
WORKDIR /app
RUN groupadd -g 1001 -r hich && \
        useradd -u 1001 -r -s /bin/false -d /app -g hich hich && \
        chown -R hich:hich /app
USER hich:hich
COPY --from=BUILD --chown=hich:hich /app/main /app
CMD ["/app/main"]
