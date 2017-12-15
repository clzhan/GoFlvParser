package parser

import (
	"github.com/clzhan/GoFlvParser/log"
	"io"
)

type FlvReader struct {
	ior    io.Reader
	log    *log.Log
	buffer []byte
}

func InitReader(r io.Reader, l *log.Log) *FlvReader {

	reader := new(FlvReader)
	reader.ior = r
	reader.log = l

	reader.log.Info("InitReader success!!")

	return reader
}

func (flv *FlvReader) FlvReadHeader() (header *FlvHeader, err error) {

	header = &FlvHeader{}

	err = header.ParseFlvHeader(flv.ior)

	return header, err
}


func (flv *FlvReader) FlvReadTag() (tag *FlvTag,avc *AVCDecoderConfigurationRecord, aac* AudioSpecificConfig, err error) {

	tag = &FlvTag{}

	err,avc,aac = tag.ParseFlvTag(flv.ior)

	return tag,avc,aac, err
}



