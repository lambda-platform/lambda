package utils

import (
	"fmt"
	"go/format"
	"os"
)

func WriteFileFormat(fileContent string, path string) error {

	formatted, err := format.Source([]byte(fileContent))
	if err != nil {
		fmt.Println(path)
		fmt.Println(fileContent)
		fmt.Println(err.Error())
		return err
	} else {
		errWrite := WriteFile(string(formatted), path)
		if errWrite != nil {
			fmt.Println(errWrite.Error())
			return errWrite
		}
		return nil
	}
}

func WriteFile(fileContent string, path string) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	l2, err := f.WriteString(string(fileContent))
	if err != nil {
		//panic(err)
		panic(l2)
		f.Close()
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
