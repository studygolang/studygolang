package service_test

import (
	"testing"

	"service"
)

func TestParseProjectList(t *testing.T) {
	service.ParseProjectList("http://www.oschina.net/project/lang/358/go?tag=0&os=0&sort=view")
}

func TestParseOneProject(t *testing.T) {
	service.ParseOneProject("http://www.oschina.net/p/docker")
}
