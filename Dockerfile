FROM golang:1.21

WORKDIR /usr/src/app

COPY [".", "."]

RUN go build -o test_service ./main.go

FROM alpine

WORKDIR /usr/src/app

COPY --from=0 /usr/src/app/test_service /bin/test_service

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

CMD ["/bin/test_service"]


