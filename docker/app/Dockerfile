FROM golang:1.23 AS build

WORKDIR /src

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/ .
RUN CGO_ENABLED=0 go build -o /out/bookingapp ./cmd/api

FROM alpine:3.20

COPY --from=build /out/bookingapp /bookingapp

EXPOSE 8080

CMD ["/bookingapp"]
