package parser

import (
	"encoding/binary"
	"fmt"
	"io"
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
	Body   *FlvTagBody
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

func (tag *FlvTag) ParseFlvTag(r io.Reader) (err error) {

	fmt.Println("This a tag......")

	if err := binary.Read(r, binary.BigEndian, &tag.Header.PreviousTagSize); err != nil {
		fmt.Println("err: " + err.Error())
		return err
	}

	fmt.Println("before tag.. len:%d ",tag.Header.PreviousTagSize)

	//tag Header
	var tmp uint8
	if err = binary.Read(r, binary.BigEndian, &tmp); err != nil {
		return err
	}
	tag.Header.Reserved = tmp >> 6 & 0x03
	tag.Header.Filter = tmp >> 5 & 0x01
	tag.Header.TagType = tmp & 0x1f

	dataSize := make([]byte, 3,3)
	if _, err = io.ReadFull(r, dataSize); err != nil {
		return err
	}
	tag.Header.DataSize = Bytes3ToUint32(dataSize)
	fmt.Println("This a tag. DataSize %u.....",tag.Header.DataSize)


	timeStamp := make([]byte, 4,4)
	if _, err = io.ReadFull(r, timeStamp); err != nil {
		return err
	}
	tag.Header.Timestamp = GetTagTimestamp(timeStamp)

	streamId := make([]byte, 3)
	if _, err = io.ReadFull(r, streamId); err != nil {
		return err
	}

	//tag data
	data := make([]byte, tag.Header.DataSize)
	if _, err = io.ReadFull(r, data); err != nil {
		return err
	}

	var body FlvTagBody=nil
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
		body.ParserTagBody(data)
	}

	return nil
}

type FlvTagBody interface {
	//audio
	//Video
	//matedata
	ParserTagBody(data []byte) (err error)
}

type AudioTagData struct {
	SoundFormat   uint8 //4bit
	SoundRate     uint8 //2bit
	SoundSize     uint8 //1bit
	SoundType     uint8 //1bit
	AACPacketType uint8 //8bit if SoundFormat== 10  defined
}

func (a *AudioTagData) ParserTagBody(data []byte) (err error) {
	fmt.Println("Audio Tag Parser..........")
	return nil
}

type VideoTagData struct {
	FrameType       uint8  //4bit
	CodecID         uint8  //4bit
	AVCPacketType   uint8  //8bit
	CompositionTime uint32 //24
}

func (v *VideoTagData) ParserTagBody(data []byte) (err error) {
	fmt.Println("Video Tag Parser..........")
	return nil
}
