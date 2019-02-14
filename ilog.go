package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

func writeLog(preFix string, contents string) bool {

	filePath, _ := filepath.Abs("log")
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(filePath, os.ModePerm)
			if err != nil {
				fmt.Println("create dir error:", err.Error())
				return false
			}
		}
	}

	filePath = path.Join(filePath, time.Now().Format("200601"))
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(filePath, os.ModePerm)
			if err != nil {
				fmt.Println("create dir error:", err.Error())
				return false
			}
		}
	}

	filePath = path.Join(filePath, time.Now().Format("20060102")+".log")
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(filePath)
			if err != nil {
				fmt.Println("create file error:", err.Error())
				return false
			}
		}
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("open file error:", err.Error())
		return false
	}
	defer file.Close()

	loger := log.New(file, "["+preFix+"] ", log.Ldate|log.Ltime|log.Lshortfile)
	loger.Println(contents + "\r\n")
	fmt.Println(contents)

	return true
}
