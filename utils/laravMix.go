package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
)

var onceMix sync.Once

var Mix map[string]string

func init() {

	onceMix.Do(func() {

		mixFile := "public/mix-manifest.json"

		if runtime.GOOS == "windows" {
			dir, _ := os.Getwd()

			fmt.Println("dir windows:")
			fmt.Println(dir)

			mixFile = dir + "/public/mix-manifest.json"
			fmt.Println(mixFile)
		}

		data, err := ioutil.ReadFile(mixFile)

		if err != nil {
			fmt.Println("MIX FILE NOT FOUND")
		}
		err2 := json.Unmarshal(data, &Mix)
		if err2 != nil {
			fmt.Println("File parse error")
		}
	})

}
