package parser

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type FLVScriptData struct {
	ScriptData []byte
}

var (
	AmfList = []string{
		"Number",
		"Boolean",
		"String",
		"Object",
		"MovieClip", // reserved, not supported
		"Null",
		"Undefined",
		"Reference",
		"ECMA Array",
		"Object and marker",
		"Strict array",
		"Date",
		"Long string",
	}
)

type ScriptDataValue struct {
	tp     string
	length uint32
	value  string
}

//TODO http://blog.csdn.net/cabbage2008/article/details/50500021
func (s *FLVScriptData) ParserTagBody(inData []byte) (err error, avc *AVCDecoderConfigurationRecord, aac *AudioSpecificConfig, outData []byte) {
	fmt.Println("script Tag Parser..........")
	var offset = len(inData)

	var j uint32 = 0
	for i := 0; i < offset; {
		tp := uint8(inData[i])
		i += 1
		var sz uint32

		switch tp {

		case 2:
			sz = uint32(inData[i])<<8 | uint32(inData[i+1])
			i += 2

			var value = string(inData[i : i+int(sz)])

			fmt.Printf("tp: %v  sz:%v value:%v\n", AmfList[tp], sz, value)
			i += int(sz)
		case 8:
			//读取数组个数
			if err := binary.Read(bytes.NewReader(inData[i:i+4]), binary.BigEndian, &sz); err != nil {
				return err, nil, nil, nil
			}
			i += 4

			Variables := inData[i:]

			var k uint32 = 0
			for k = 0; k < sz; {
				StringLength := uint32(Variables[j])<<8 | uint32(Variables[j+1])
				j += 2
				StringData := string(Variables[j : j+StringLength])
				fmt.Printf("StringData :%v\n", StringData)
				j += StringLength
				tmp8 := Variables[j]
				j += 1
				switch tmp8 {
				case 0: //double 8个字节
					floatbyte := Variables[j : j+8]
					fmt.Printf("floatbyte :%v\n", hex.Dump(floatbyte))
					j += 8
				case 1: //Boolean 8个字节
					Booleanbyte := Variables[j : j+1]
					fmt.Printf("Boolean :%v\n", hex.Dump(Booleanbyte))
					j += 1
				case 2: //string
					StringLength2 := uint32(Variables[j])<<8 | uint32(Variables[j+1])
					j += 2
					stringbyte := Variables[j : j+StringLength2]
					fmt.Printf("stringbyte :%v\n", hex.Dump(stringbyte))
					j += StringLength2
				default:
					defaultbyte := Variables[j : j+8]
					fmt.Printf("defaultbyte :%v\n", defaultbyte)
					j += 8

				}
				k += 1
			}
			i +=int(j)

		}

		//fmt.Println(amf)
	}

	return nil, nil, nil, nil
}
