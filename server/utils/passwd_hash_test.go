package utils

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestHashIt(t *testing.T) {
	result := HashIt("lx_admin123456")
	log.Println(result)
}
