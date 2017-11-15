# This file decribes the standard way to build stadygolang, using docker
#
# # Usage
#
# # # download the src and enter the dir first
# docker build -f Dockerfile.web -t studygolang .
#
# docker run --name mysqlDB -e MYSQL_ROOT_PASSWORD=123456 -d mysql
# docker run -d --name studygolang-web -v `pwd`:/studyglang -p 8090:8088 --link mysqlDB:db.localhost studygolang ./docker-entrypoint.sh
#
# # inside the container
# bin/studygolang
#
# # just compile
# docker run --rm -v `pwd`:/studyglang ./install.sh
# # and in production environment just put this binary file in jockerxu/ubuntu-golang and run it


FROM jockerxu/ubuntu-golang
MAINTAINER jockerxu <156082052@qq.com>

# download dep
RUN go get github.com/polaris1119/gvt
WORKDIR /studygolang
COPY . /studygolang
RUN cd src/ && gvt restore
RUN mkdir -p /vendor/src/ && mv src/vendor/* /vendor/src/
ENV GOPATH $GOPATH:/vendor

# run
CMD ["docker-entrypoint.sh"]
