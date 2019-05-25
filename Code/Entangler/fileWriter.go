package Entangler

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func WriteFile(data []byte, back, forward int) {
	filename := "p" + strconv.Itoa(back) + "_" + strconv.Itoa(forward)
	if _, err := os.Create(filename); err == nil {

	} else {
		fmt.Println("Fatal error ... " + err.Error())
		os.Exit(1)
	}

	ioutil.WriteFile(filename, data, os.ModeAppend)
}
