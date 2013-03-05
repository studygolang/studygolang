package process

import (
	"fmt"
	"os"
)

func Die(s string) {
	fmt.Println(s)
	os.Exit(1)
}
