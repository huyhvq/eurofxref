FROM golang:1.16-alpine as build

WORKDIR /ws/eurofxref
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o /out/eurofxref .

FROM alpine:latest
WORKDIR /ws/eurofxref
COPY --from=build /out/eurofxref /ws/eurofxref/eurofxref
COPY --from=build /ws/eurofxref/migrations /ws/eurofxref/migrations
EXPOSE 8080
ENTRYPOINT ["/ws/eurofxref/eurofxref"]