# Start from golang v1.12 base image
FROM golang:1.12

WORKDIR /app/studygolang

COPY . /app/studygolang

RUN make build

CMD ["bin/studygolang"]