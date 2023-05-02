package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var onceMix sync.Once

var Mix map[string]string

func init() {

	onceMix.Do(func() {

		mixFile := "public/mix-manifest.json"

		//if runtime.GOOS == "windows" {
		//	dir, _ := os.Getwd()
		//
		//	fmt.Println("dir windows:")
		//	fmt.Println(dir)
		//
		//	mixFile = dir + "/public/mix-manifest.json"
		//	fmt.Println(mixFile)
		//}

		data, err := os.ReadFile(mixFile)

		if err == nil {
			err2 := json.Unmarshal(data, &Mix)
			if err2 != nil {
				fmt.Println("File parse error")
			}
		}

	})

}
