package main

import (
	"fmt"

	"os"

	"github.com/clzhan/GoFlvParser/parser"
	"github.com/stackimpact/stackimpact-go"
)

const FileNmae = "flv_test.flv"

func main() {

	agent := stackimpact.Start(stackimpact.Options{
		AgentKey: "100ef4628a42754ee6f32391ec1d784225cecf5a",
		AppName: "MyGoApp",
	})

	fmt.Println("This is a flv parser ",agent)

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
