package util

import (
	"os"
	"fmt"
)

//WriteFile -
func WriteFile(name,content string)(err error){
	fileObj,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_TRUNC,0644)
	if err != nil {
		fmt.Println("Failed to open the file",err.Error())
		return
	}
	defer fileObj.Close()
	if _,err := fileObj.WriteString(content);err != nil {
		fmt.Println("Failed to write the file",err.Error())
	}
	return
}

//ReadFile-
func ReadFile(name string)(context string,err error){
	fileObj,err := os.Open(name)
	if err != nil {
		fmt.Println("Failed to open the file",err.Error())
		return
	}
	defer fileObj.Close()
	buf := make([]byte,1024)
	if _,err := fileObj.Read(buf);err != nil {
		fmt.Println("Failed to read the file",err.Error())
	}
	return string(buf),err
}
