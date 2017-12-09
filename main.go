package main

import (
	"fmt"

	"os"

	"github.com/clzhan/GoFlvParser/parser"
)

const FileNmae = "flv_test.flv"

func main() {

	fmt.Println("This is a flv parser")

	f, err := os.Open(FileNmae)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	flvHeader := &parser.FlvHeader{}

	if err = flvHeader.ParseFlvHeader(f); err != nil{
		fmt.Println("main err : " + err.Error())
		return
	}

	for{
		tag := &parser.FlvTag{}
		if err= tag.ParseFlvTag(f);err != nil{
			fmt.Println("err : "+ err.Error())
			break;
		}
	}


	fmt.Println("decode tag finished")
}
