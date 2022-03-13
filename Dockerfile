# Start from golang v1.17 base image
FROM golang:1.17

WORKDIR /app/studygolang

COPY . /app/studygolang/

RUN make

ENTRYPOINT ["bin/studygolang", "-embed_crawler", "-embed_indexing"]
