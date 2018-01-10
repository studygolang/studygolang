.PHONY: getpkg install reload start stop migrate

v=""

getpkg:
	./getpkg.sh

install:
	./install.sh

reload:
	./reload.sh	

start:
	./start.sh	

stop:
	./stop.sh	

migrate:
	./bin/migrator --changeVersion=${v}
