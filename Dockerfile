# builder image
FROM golang:1.17-alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build/cmd/web/
RUN CGO_ENABLED=0 GOOS=linux go build -a -o web-app .


# generate clean, final image for end users
FROM alpine:3.15
COPY --from=builder /build/cmd/web/web-app .

# executable
ENTRYPOINT [ "./web-app" ]
