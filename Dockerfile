FROM golang:1.15-alpine AS build
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o graphql-schema-registry .

FROM alpine:3.12
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /go/src/app/graphql-schema-registry /app/graphql-schema-registry
EXPOSE 8080
CMD ["/app/graphql-schema-registry"]
