# build binary
FROM golang:1.11-alpine3.7 AS build

RUN apk add --no-cache git

WORKDIR /src

COPY . /src/

RUN go mod download

RUN CGO_ENABLED=0 go build -o /out/resizer cmd/resizer/main.go

# copy to alpine image
FROM alpine:3.7
WORKDIR /app
COPY --from=build /out/resizer /app
CMD ["/app/resizer"]

