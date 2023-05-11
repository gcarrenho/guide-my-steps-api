FROM golang:alpine AS build
WORKDIR $GOPATH/src
COPY . .
RUN go build -o guide-my-steps src/cmd/main.go


# Building image with the binary
FROM scratch
COPY --from=build /go/src/guide-my-steps .
ENTRYPOINT ["./guide-my-steps"]