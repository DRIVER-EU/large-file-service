FROM golang:latest as builder

COPY . /build/
WORKDIR /build

RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o large-file-service

FROM scratch

COPY --from=builder /build/large-file-service /app/
COPY --from=builder /build/templates /app/templates
COPY --from=builder /build/swagger /app/swagger
WORKDIR /app
ENTRYPOINT ["./large-file-service"]