package utils

import (
	"encoding/json"
	"fmt"
	"runtime"
	"io/ioutil"
	"os"
)


func FrontMix(path string) map[string]string  {
	var frontEndMix map[string]string
	//AbsolutePath := config.AbsolutePath()

	mixFile := "public/"+path

	if runtime.GOOS == "windows" {
		dir, _ := os.Getwd()

		fmt.Println("windows dir front:")
		fmt.Println(dir)

		mixFile := dir + "/public/" + path
		fmt.Println(mixFile)
	}

	data, err := ioutil.ReadFile(mixFile)

	if err != nil{
		fmt.Println("MIX FILE NOT FOUND")
	}
	err2 := json.Unmarshal(data, &frontEndMix)
	if err2 != nil{
		fmt.Println("File parse error")
	}

	return  frontEndMix
}


func CallMix(index string, mixData map[string]string) string{
	return mixData[index]
}