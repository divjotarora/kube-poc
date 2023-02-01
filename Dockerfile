FROM golang:1.18-alpine as base
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main ./main.go
ENTRYPOINT ["/app/main"]

FROM base as build_reader
CMD ["--mode", "reader"]

FROM base as build_writer
CMD ["--mode", "writer"]