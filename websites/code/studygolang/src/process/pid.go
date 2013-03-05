package process

import (
	"io/ioutil"
	"os"
	"strconv"
)

func SavePidTo(pidFile string) error {
	return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(os.Getpid())), 0777)
}
