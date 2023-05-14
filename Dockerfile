FROM golang:alpine AS build
WORKDIR $GOPATH/src
COPY . .
RUN go build -o guide-my-steps src/cmd/main.go


# Building image with the binary
FROM alpine:latest AS production
COPY --from=build /go/src/guide-my-steps .
COPY src/internal/locales ./internal/locales 
EXPOSE 8080
CMD ["./guide-my-steps"]