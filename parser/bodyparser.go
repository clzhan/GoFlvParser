package parser

type MediaBodyParser interface {
	//audio
	//Video
	//matedata
	ParserTagBody(inData []byte) (err error,avc *AVCDecoderConfigurationRecord, aac* AudioSpecificConfig, outData []byte)
}
