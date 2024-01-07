FROM golang:alpine AS builder
WORKDIR /app
ADD go.mod .
COPY . . 
RUN go build -o app ./main.go

FROM alpine
RUN apk update && \
    apk add --no-cache curl
WORKDIR /app
COPY --from=builder /app/app /app/app
HEALTHCHECK --interval=10s --timeout=5s --start-period=3s \ 
   CMD curl --fail localhost || exit 1
CMD ["/app/app"]