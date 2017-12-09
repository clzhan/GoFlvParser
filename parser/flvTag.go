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

	//
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
	AACPacketType uint8 //8bit if SoundFormat== 10  defined 如果不是AAC编码 没有这个字节
}

type AudioSpecificConfig struct{
	AacProfile  uint8   //5bit
	SampleRateIndex uint8 //4bit
	ChannelConfig  uint8  //4bit
	OtherConfig   uint8  //3bit      这个为0
}

func (a *AudioTagData) ParserTagBody(data []byte) (err error) {
	fmt.Println("Audio Tag Parser..........")
	tmp := data[0]

	a.SoundFormat = tmp >> 4 & 0x0f
	a.SoundRate = tmp >> 2 & 0x03
	a.SoundSize = tmp >> 1 & 0x01
	a.SoundType = tmp  & 0x01

	if a.SoundFormat ==10{
		a.AACPacketType = data[1]

		if a.AACPacketType == 0{
			//AudioSpecificConfig 2个字节，这个值就是faacEncGetDecoderSpecificInfo出来的
			//前5位，表示编码结构类型，AAC main编码为1，LOW低复杂度编码为2，SSR为3
			//4位，表示采样率。 按理说，应该是：0 ~ 96000， 1~88200， 2~64000， 3~48000， 4~44100， 5~32000， 6~24000， 7~ 22050， 8~16000...)
			// 通常aac固定选中44100，即应该对应为4，但是试验结果表明，当音频采样率小于等于44100时，应该选择3，而当音频采样率为48000时，应该选择2
			//接着4位，表示声道数。
			//最后3位，固定为0吧

			var config AudioSpecificConfig

			config.AacProfile  = data[2] >> 3 & 0x1f
			config.SampleRateIndex  = ((data[2]&0x07)<<1) | (data[3]>>7)
			config.ChannelConfig = (data[3]>>3) & 0x0f;
			config.OtherConfig = data[3] & 0x03;
			fmt.Printf("aac AacProfile:%v, SampleRateIndex:%v, ChannelConfig:%v, OtherConfig:%v\n",
				GetAACProfile(config.AacProfile),GetSampleRate(config.SampleRateIndex),GetChannel(config.ChannelConfig) ,config.OtherConfig)
		}else{
			//aac data


		}


		fmt.Printf("audio format:%v, rate:%v, size:%v, type:%v aactype:%v \n",
			GetSoundFormat(a.SoundFormat),GetSoundRate(a.SoundRate),GetSoundSize(a.SoundSize),GetSoundType(a.SoundType), GetAACPacketType(a.AACPacketType))

	} else {

	}



	return nil
}

type VideoTagData struct {
	FrameType       uint8  //4bit 帧类型
	CodecID         uint8  //4bit 视频编码类型
	AVCPacketType   uint8  //8bit
	CompositionTime uint32 //24
}

func (v *VideoTagData) ParserTagBody(data []byte) (err error) {
	fmt.Println("Video Tag Parser..........")

	tmp := data[0]

	v.FrameType = tmp >> 4 & 0x0f
	v.CodecID = tmp  & 0x0f

	v.AVCPacketType = data[1]

	if v.AVCPacketType == 0 {
		//AVC sequence header


	}else{

	}




	return nil
}
