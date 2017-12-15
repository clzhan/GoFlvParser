package parser

import "fmt"

type AudioTagData struct {
	SoundFormat   uint8 //4bit
	SoundRate     uint8 //2bit
	SoundSize     uint8 //1bit
	SoundType     uint8 //1bit
	AACPacketType uint8 //8bit if SoundFormat== 10  defined 如果不是AAC编码 没有这个字节

	SoundData []byte
}

type AudioSpecificConfig struct {
	AacProfile      uint8 //5bit
	SampleRateIndex uint8 //4bit
	ChannelConfig   uint8 //4bit
	OtherConfig     uint8 //3bit      这个为0
}

var soundFormat = map[uint8]string{
	0:  "PCM-platform-endian",
	1:  "ADPCM",
	2:  "MP3",
	3:  "PCM-little-endian",
	4:  "Nellymoser-16kHz",
	5:  "Nellymoser-8kHz",
	6:  "Nellymoser",
	7:  "G.711-A-law",
	8:  "G.711-mu-law",
	9:  "reserved",
	10: "AAC",
	11: "Speex",
	14: "MP3-8kHz",
	15: "Device-specific",
}
var soundRate = map[uint8]string{
	0: "5.5kHz",
	1: "11kHz",
	2: "22kHz",
	3: "44kHz",
}

var soundSize = map[uint8]string{
	0: "8-Bit",
	1: "16-Bit",
}

var soundType = map[uint8]string{
	0: "Mono",
	1: "Stereo",
}
var aacPacketType = map[uint8]string{
	0: "AAC-SeqHeader",
	1: "AAC-Raw",
}

var aacProfile = map[uint8]string{
	1: "AAC Main",
	2: "AAC LC",
	3: "AAC SSR",
}

var sampleRate = map[uint8]int{
	0:  96000,
	1:  88200,
	2:  64000,
	3:  48000,
	4:  44100,
	5:  32000,
	6:  24000,
	11: 8000,
}

var achannel = map[uint8]string{
	1: "单声道", //center front speaker
	2: "双声道", //left, right front speakers
	3: "三声道", //center, left, right front speakers
	4: "四声道", //center, left, right front speakers, rear surround speakers
}

func (a *AudioTagData) ParserTagBody(inData []byte) (err error,avc *AVCDecoderConfigurationRecord, aac* AudioSpecificConfig, outData []byte){
	fmt.Println("Audio Tag Parser..........")
	tmp := inData[0]

	a.SoundFormat = tmp >> 4 & 0x0f
	a.SoundRate = tmp >> 2 & 0x03
	a.SoundSize = tmp >> 1 & 0x01
	a.SoundType = tmp & 0x01

	if a.SoundFormat == 10 {
		a.AACPacketType = inData[1]

		if a.AACPacketType == 0 {
			//AudioSpecificConfig 2个字节，这个值就是faacEncGetDecoderSpecificInfo出来的
			//前5位，表示编码结构类型，AAC main编码为1，LOW低复杂度编码为2，SSR为3
			//4位，表示采样率。 按理说，应该是：0 ~ 96000， 1~88200， 2~64000， 3~48000， 4~44100， 5~32000， 6~24000， 7~ 22050， 8~16000...)
			// 通常aac固定选中44100，即应该对应为4，但是试验结果表明，当音频采样率小于等于44100时，应该选择3，而当音频采样率为48000时，应该选择2
			//接着4位，表示声道数。
			//最后3位，固定为0吧

			var config AudioSpecificConfig

			config.AacProfile = inData[2] >> 3 & 0x1f
			config.SampleRateIndex = ((inData[2] & 0x07) << 1) | (inData[3] >> 7)
			config.ChannelConfig = (inData[3] >> 3) & 0x0f
			config.OtherConfig = inData[3] & 0x03
			fmt.Printf("aac AacProfile:%v, SampleRateIndex:%v, ChannelConfig:%v, OtherConfig:%v\n",
				aacProfile[config.AacProfile], sampleRate[config.SampleRateIndex], achannel[config.ChannelConfig], config.OtherConfig)

			return nil, nil,&config,nil
		} else {
			//aac data
			//a.SoundData = inData[2:]
			outData = inData[2:]
			return nil, nil,nil,outData

		}

		//fmt.Printf("audio format:%v, rate:%v, size:%v, type:%v aactype:%v \n",
		//	soundFormat[a.SoundFormat], soundRate[a.SoundRate], soundSize[a.SoundSize], soundType[a.SoundType], aacPacketType[a.AACPacketType])

	} else {

	}

	return nil,nil,nil,nil
}
