package parser

import "fmt"

type FLVScriptData struct {
	ScriptData []byte
}

func (s *FLVScriptData) ParserTagBody(inData []byte) (err error, avc *AVCDecoderConfigurationRecord, aac *AudioSpecificConfig, outData []byte) {
	fmt.Println("script Tag Parser..........")






	return nil, nil, nil, nil
}
