FROM golang:alpine AS build
WORKDIR /go/src/guide-my-steps
COPY . .
RUN go build -o /go/bin/guide-my-steps cmd/main.go


# Building image with the binary
FROM scratch
COPY --from=build /go/bin/guide-my-steps /go/bin/guide-my-steps
ENTRYPOINT ["/go/bin/guide-my-steps"]