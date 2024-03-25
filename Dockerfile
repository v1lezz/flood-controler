FROM golang:latest
RUN mkdir /checker
ADD . /checker/
WORKDIR /checker
RUN go build -o cmd/main ./cmd
CMD ["./cmd/main"]