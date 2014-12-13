all: build install
build:
	go build 

install: build
	go build
	cp templates/* /opt/ameske/go_nfl/templates/
	cp go_nfl /opt/ameske/go_nfl/
