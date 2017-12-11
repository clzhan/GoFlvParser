package parser

import (
	"encoding/binary"
	"fmt"
	"io"
	"container/list"
	"errors"
	"encoding/hex"
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

type AVCDecoderConfigurationRecord struct {
	ConfigurationVersion uint8 // 1 bytes
	AVCProfileIndication uint8 // 1 bytes
	ProfileCompatibility uint8 // 1 bytes
	AVCLevelIndication uint8   // 1 bytes
	LengthSizeMinusOne           uint8
	//number_of_sequence_parameter_sets 1 bytes
	NumOfSequenceParameterSets   int
	SPS                          *list.List
	NumOfPictureParameterSets    int
	PPS                          *list.List

	//SPS []byte //sequence parameter set length(2 bytes)

	// number of picture parameter sets(1 bytes)
	//PPS []byte // picture parameter set length (2 bytes)
}
//AVCDecoderConfigurationRecord
/*
aligned(8) class AVCDecoderConfigurationRecord {
    ||0		unsigned int(8) configurationVersion = 1;
    ||1		unsigned int(8) AVCProfileIndication;
    ||2		unsigned int(8) profile_compatibility;
    ||3		unsigned int(8) AVCLevelIndication;
    ||4		bit(6) reserved = ‘111111’b;
            unsigned int(2) lengthSizeMinusOne; // offset 4
    ||5		bit(3) reserved = ‘111’b;
            unsigned int(5) numOfSequenceParameterSets;
    ||6		for (i = 0; i< numOfSequenceParameterSets; i++) {
                ||0	    unsigned int(16) sequenceParameterSetLength;
                ||2	    bit(8 * sequenceParameterSetLength) sequenceParameterSetNALUnit;
    }
    ||6+X	unsigned int(8) numOfPictureParameterSets;
            for (i = 0; i< numOfPictureParameterSets; i++) {
                ||0		unsigned int(16) pictureParameterSetLength;
                ||2		bit(8 * pictureParameterSetLength) pictureParameterSetNALUnit;
            }
}
*/

//HEVCDecoderConfigurationRecord
/*
aligned(8) class HEVCDecoderConfigurationRecord
{
    ||0		unsigned int(8) configurationVersion = 1;
    //vps[4]
    ||1		unsigned int(2) general_profile_space;
            unsigned int(1) general_tier_flag;
            unsigned int(5) general_profile_idc;
    //vps[5..8]
    ||2		unsigned int(32) general_profile_compatibility_flags;
    //
    ||6		unsigned int(48) general_constraint_indicator_flags;
    //vps[14]
    ||12	unsigned int(8) general_level_idc;
    ||13	bit(4) reserved = ‘1111’b;
            unsigned int(12) min_spatial_segmentation_idc;
    ||15	bit(6) reserved = ‘111111’b;
            unsigned int(2) parallelismType;
    ||16	bit(6) reserved = ‘111111’b;
            unsigned int(2) chroma_format_idc;
    ||17	bit(5) reserved = ‘11111’b;
            unsigned int(3) bit_depth_luma_minus8; //0
    ||18	bit(5) reserved = ‘11111’b;
            unsigned int(3) bit_depth_chroma_minus8; //0
    ||19	bit(16) avgFrameRate;
    ||21	bit(2) constantFrameRate;
            bit(3) numTemporalLayers;
            bit(1) temporalIdNested;
            unsigned int(2) lengthSizeMinusOne;
    ||22	unsigned int(8) numOfArrays;
    ||23	for (j=0; j < numOfArrays; j++)
            {
        ||0		bit(1) array_completeness;
                unsigned int(1) reserved = 0;
                unsigned int(6) NAL_unit_type;
        ||1		unsigned int(16) numNalus;
        ||3		for (i=0; i< numNalus; i++)
                {
                    unsigned int(16) nalUnitLength;
                    bit(8*nalUnitLength) nalUnit;
                }
            }
}
*/
func (v *VideoTagData) ParserTagBody(data []byte) (err error) {
	fmt.Println("Video Tag Parser..........")

	tmp := data[0]

	v.FrameType = tmp >> 4 & 0x0f
	v.CodecID = tmp  & 0x0f

	if v.CodecID == 7{  //avc 才有的字段
		v.AVCPacketType = data[1]
		v.CompositionTime = Bytes3ToUint32(data[2:5])

		if v.AVCPacketType == 0 {
			//AVC sequence header
			tmp = data[5]
			var info AVCDecoderConfigurationRecord

			info.ConfigurationVersion = data[5]
			info.AVCProfileIndication = data[6]
			info.ProfileCompatibility = data[7]
			info.AVCLevelIndication = data[8]

			info.LengthSizeMinusOne = data[9] & 0x03

			info.NumOfSequenceParameterSets = int(data[10] & 0x1f)

			fmt.Printf("lengthSizeMinusOne : %v number_of_sequence_parameter_sets:%v\n",info.LengthSizeMinusOne,info.NumOfSequenceParameterSets)

			//var i uint8 = 0
			//for i = 0;i < number_of_sequence_parameter_sets;i++{
			//	fmt.Println(i)
			//
			//}
			info.SPS = list.New()
			var size int
			for i := 0; i < info.NumOfSequenceParameterSets; i++ {

				size = (int(data[11]) << 8) | (int(data[12]))
				if size == 0 {
					err = errors.New("invalid sps size")
					return
				}

				dataSps := data[13:13 + size]

				fmt.Printf("dataSps : %v \n",hex.Dump(dataSps))
				info.SPS.PushBack(dataSps)
			}

			datapps := data[13+size:]
			info.PPS = list.New()
			info.NumOfPictureParameterSets =int(datapps[0])

			for i := 0; i < info.NumOfPictureParameterSets; i++ {

				size = (int(datapps[1]) << 8) | (int(datapps[2]))
				if size == 0 {
					err = errors.New("invalid pps size")
					return
				}

				dataPps := datapps[3:]

				fmt.Printf("dataPps : %v \n",hex.Dump(dataPps))
				info.SPS.PushBack(dataPps)
			}


		}else{
			//data


		}
	}else{
		fmt.Println("Not AVC Encode")
	}







	return nil
}
/*
aligned(8) class AVCDecoderConfigurationRecord {
    ||0		unsigned int(8) configurationVersion = 1;
    ||1		unsigned int(8) AVCProfileIndication;
    ||2		unsigned int(8) profile_compatibility;
    ||3		unsigned int(8) AVCLevelIndication;
    ||4		bit(6) reserved = ‘111111’b;
            unsigned int(2) lengthSizeMinusOne; // offset 4
    ||5		bit(3) reserved = ‘111’b;
            unsigned int(5) numOfSequenceParameterSets;
    ||6		for (i = 0; i< numOfSequenceParameterSets; i++) {
                ||0	    unsigned int(16) sequenceParameterSetLength;
                ||2	    bit(8 * sequenceParameterSetLength) sequenceParameterSetNALUnit;
    }
    ||6+X	unsigned int(8) numOfPictureParameterSets;
            for (i = 0; i< numOfPictureParameterSets; i++) {
                ||0		unsigned int(16) pictureParameterSetLength;
                ||2		bit(8 * pictureParameterSetLength) pictureParameterSetNALUnit;
            }
}
*/