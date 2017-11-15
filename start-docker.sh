#!/usr/bin/env bash

# ***************************************************************************
# *
# * @author:jockerxu
# * @date:2017-11-14 22:20
# * @version 1.0
# * @description: Shell script
#*
#**************************************************************************/

#---------tool function---------------
echo_COLOR_GREEN=$(  echo -e "\e[32;49m")
echo_COLOR_RESET=$(  echo -e "\e[0m")
function echo-info()
{
    echo -e "${echo_COLOR_GREEN}[$(date "+%F %T")]\t$*${echo_COLOR_RESET}";
}
#---------end tool function-----------
if [[ $USER != "root" ]]; then
    echo "you must be root!!!!!"
    exit 1
fi

if [[ $1 == "" ]]; then
    echo "Usage start-docker.sh [local | remote]"
    exit 1
fi

STUDYGOLANG_IMG=

if [[ $1 == "local" ]]; then
    STUDYGOLANG_IMG=studygolang
    docker images ${STUDYGOLANG_IMG} | grep -q ${STUDYGOLANG_IMG} || {
        docker build -f Dockerfile.web -t $STUDYGOLANG_IMG .
    }
elif [[ $1 == "remote" ]]; then
    STUDYGOLANG_IMG="jockerxu/studygolang"
else
    exit 1
fi

docker ps -a | grep -q mysqlDB || {
    docker run --name mysqlDB -e MYSQL_ROOT_PASSWORD=123456 -d mysql
}
docker ps -a | grep -q studygolang-web && {
    docker rm -f studygolang-web
}
docker run -d --name studygolang-web -v `pwd`:/studygolang -p 8090:8088 --link mysqlDB:db.localhost $STUDYGOLANG_IMG ./docker-entrypoint.sh

if [[ $? == 0 ]]; then
    echo-info "studygolang-web start, waiting several seconds to install..."
    sleep 5
    echo-info "open browser: http://localhost:8090"
    echo-info "mysql-host is: db.localhost "
    echo-info "mysql-password is: 123456"
fi
