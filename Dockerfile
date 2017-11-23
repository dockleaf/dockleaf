FROM golang:1.9 as builder
WORKDIR /go/src/github.com/edkellena/dockleaf
ADD dockleaf.go .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dockleaf .


FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/src/github.com/edkellena/dockleaf /app/
ENTRYPOINT ./dockleaf