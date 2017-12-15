package parser

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	FLVAudio  = 0x08
	FLVVideo  = 0x09
	FLVScript = 0x12
)

//flv tag
type FlvTagHeader struct {
	PreviousTagSize   uint32 //前一个tag的长度
	Filter            uint8  //1bit 表示是否经过滤波  一般为0
	Reserved          uint8  // 2bit  Reserved for FMS, should be 0
	TagType           uint8  //5bit	- 8 = audio , 9 = video , 18 = script data
	DataSize          uint32 //24bit 表示当前tag的后续长度等于当前整个tag长度减去11（tag头信息）
	Timestamp         uint32 //24bit
	TimestampExtended uint8  //8bit
	StreamID          uint32 //24bit 一直为0

}

type FlvTag struct {
	Header FlvTagHeader
	Body   *MediaBodyParser

	Buffer []byte
}

func Bytes3ToUint32(b []byte) uint32 {
	nb := []byte{}
	nb = append(nb, 0)
	nb = append(nb, b...)
	return binary.BigEndian.Uint32(nb)
}

//TODO
func GetTagTimestamp(ts []byte) uint32 {
	nb := []byte{}
	//nb = append(nb, ts[3])
	nb = append(nb, 0)
	nb = append(nb, ts[0:3]...)
	return binary.BigEndian.Uint32(nb)
}

func (tag *FlvTag) ParseFlvTag(r io.Reader) (err error, avc *AVCDecoderConfigurationRecord, aac *AudioSpecificConfig) {

	fmt.Println("This a tag......")

	//
	if err := binary.Read(r, binary.BigEndian, &tag.Header.PreviousTagSize); err != nil {
		fmt.Println("err: " + err.Error())
		return err, nil, nil
	}

	fmt.Println("before tag.. len:%d ", tag.Header.PreviousTagSize)

	//tag Header
	var tmp uint8
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return err, nil, nil
	}
	tag.Header.Reserved = tmp >> 6 & 0x03
	tag.Header.Filter = tmp >> 5 & 0x01
	tag.Header.TagType = tmp & 0x1f

	dataSize := make([]byte, 3, 3)
	if _, err = io.ReadFull(r, dataSize); err != nil {
		return err, nil, nil
	}
	tag.Header.DataSize = Bytes3ToUint32(dataSize)
	fmt.Println("This a tag. DataSize %u.....", tag.Header.DataSize)

	timeStamp := make([]byte, 4, 4)
	if _, err = io.ReadFull(r, timeStamp); err != nil {
		return err, nil, nil
	}
	tag.Header.Timestamp = GetTagTimestamp(timeStamp)

	streamId := make([]byte, 3)
	if _, err = io.ReadFull(r, streamId); err != nil {
		return err, nil, nil
	}

	//tag data
	data := make([]byte, tag.Header.DataSize)
	if _, err = io.ReadFull(r, data); err != nil {
		return err, nil, nil
	}

	var body MediaBodyParser = nil
	fmt.Println("...............")
	switch tag.Header.TagType {
	case 8:
		fmt.Println("audio..........")
		body = &AudioTagData{}
	case 9:
		fmt.Println("video..........")
		body = &VideoTagData{}

	}

	if body != nil {
		_, avc, aac, buffer := body.ParserTagBody(data)
		tag.Buffer = buffer
		if tag.Header.TagType == FLVAudio {
			//fmt.Println("....buffer ",hex.Dump(tag.Buffer))
		}

		return err, avc, aac
	}

	return nil, nil, nil
}
