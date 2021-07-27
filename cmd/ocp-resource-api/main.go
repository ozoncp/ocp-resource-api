package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("This is an OCP Resource API")

	OpenFile := func(filename string) (err error) {
		file, err := os.Open(filename)
		if err != nil {
			return
		}
		defer func() {
			fileCloseErr := file.Close()
			if fileCloseErr != nil {
				fmt.Printf("Error during file closing %v: %v\n", filename, err)
				err = fileCloseErr
				return
			}
			fmt.Printf("Closed file %v\n", filename)
		}()
		return
	}
	configFiles := []string{
		"config.json", "config.prod.json", "config.aws.json",
	}

	for _, configFileName := range configFiles {
		if configFileErr := OpenFile(configFileName); configFileErr != nil {
			fmt.Printf("Error during opening config file %v: %v\n", configFileName, configFileErr)
			continue
		}
		fmt.Printf("Success opening of file %v\n", configFileName)
	}
}
