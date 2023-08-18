FROM golang:1.18 as builder

WORKDIR /app
COPY . .
RUN go get -t ./...
RUN go build -o ./bin/app ./cmd/app
WORKDIR /app/bin
RUN chmod +x app

FROM debian:bookworm-slim
COPY --from=builder /app/bin/app ./app
# RUN apt update
# RUN apt install -y --no-install-recommends ca-certificates
# RUN update-ca-certificates
CMD ./app
