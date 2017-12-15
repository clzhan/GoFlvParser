package parser

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type VideoTagData struct {
	FrameType       uint8  //4bit 帧类型
	CodecID         uint8  //4bit 视频编码类型
	AVCPacketType   uint8  //8bit
	CompositionTime uint32 //24

	AVCPayload []byte
}

type AVCDecoderConfigurationRecord struct {
	ConfigurationVersion uint8 // 1 bytes
	AVCProfileIndication uint8 // 1 bytes
	ProfileCompatibility uint8 // 1 bytes
	AVCLevelIndication   uint8 // 1 bytes
	LengthSizeMinusOne   uint8
	//number_of_sequence_parameter_sets 1 bytes
	NumOfSequenceParameterSets int
	SPS                        *list.List
	NumOfPictureParameterSets  int
	PPS                        *list.List

	//SPS []byte //sequence parameter set length(2 bytes)

	// number of picture parameter sets(1 bytes)
	//PPS []byte // picture parameter set length (2 bytes)
}

var codecId = map[uint8]string{
	2: "H263",
	3: "Screen video",
	4: "VP6",
	5: "VP6 with alpha channel",
	6: "Screen video 2",
	7: "AVC",
}

var vframeType = map[uint8]string{
	1: "key frame",
	2: "inter frame",
	3: "disposable inter frame",
	4: "generated key frame",
	5: "video info/command frame",
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
func ReadUint32(b []byte) uint32 {
	var ui uint32

	sr := bytes.NewReader(b[0:4])
	binary.Read(sr, binary.BigEndian, &ui)

	return ui
}

func (v *VideoTagData) ParserTagBody(inData []byte) (err error, avc *AVCDecoderConfigurationRecord, aac *AudioSpecificConfig, outData []byte) {
	fmt.Println("Video Tag Parser..........")

	tmp := inData[0]

	v.FrameType = tmp >> 4 & 0x0f
	v.CodecID = tmp & 0x0f

	if v.CodecID == 7 { //avc 才有的字段
		v.AVCPacketType = inData[1]
		v.CompositionTime = Bytes3ToUint32(inData[2:5])

		if v.AVCPacketType == 0 {
			//AVC sequence header
			tmp = inData[5]
			var info AVCDecoderConfigurationRecord

			info.ConfigurationVersion = inData[5]
			info.AVCProfileIndication = inData[6]
			info.ProfileCompatibility = inData[7]
			info.AVCLevelIndication = inData[8]

			info.LengthSizeMinusOne = inData[9] & 0x03

			info.NumOfSequenceParameterSets = int(inData[10] & 0x1f)

			fmt.Printf("lengthSizeMinusOne : %v number_of_sequence_parameter_sets:%v\n", info.LengthSizeMinusOne, info.NumOfSequenceParameterSets)

			//var i uint8 = 0
			//for i = 0;i < number_of_sequence_parameter_sets;i++{
			//	fmt.Println(i)
			//
			//}
			info.SPS = list.New()
			var size int
			for i := 0; i < info.NumOfSequenceParameterSets; i++ {

				size = (int(inData[11]) << 8) | (int(inData[12]))
				if size == 0 {
					err = errors.New("invalid sps size")
					return err, nil, nil, nil
				}

				dataSps := inData[13 : 13+size]

				fmt.Printf("dataSps : %v \n", hex.Dump(dataSps))
				info.SPS.PushBack(dataSps)
			}

			datapps := inData[13+size:]
			info.PPS = list.New()
			info.NumOfPictureParameterSets = int(datapps[0])

			for i := 0; i < info.NumOfPictureParameterSets; i++ {

				size = (int(datapps[1]) << 8) | (int(datapps[2]))
				if size == 0 {
					err = errors.New("invalid pps size")
					return err, nil, nil, nil
				}

				dataPps := datapps[3:]

				fmt.Printf("dataPps : %v \n", hex.Dump(dataPps))
				info.PPS.PushBack(dataPps)
			}

			return nil, &info, nil, nil
		} else {
			var start uint32
			var end uint32
			var h264prefix = [4]byte{0x00,0x00,0x00,0x01}

			v.AVCPayload = inData[5:]
			end = uint32(len(v.AVCPayload))
			s :=make([]byte,)

			for start < end {
				len := ReadUint32(v.AVCPayload[start:])
				start += 4

				s1 :=v.AVCPayload[start:start+len]
				 s = append(s,s1...)

				start = start + len

				if start < end{
					s=append(s,h264prefix[:]...)
				}

				//nalu := new(AVCNalUnit)
				//_, err := nalu.Parser(r, b[start:start+len])
				//if err != nil {
				//	return media.Error, err
				//}

			}
			outData = s;

			//data
			//outData = inData[5 + 4:]
			return nil, nil, nil, outData

		}
	} else {
		fmt.Println("Not AVC Encode")
	}

	return nil, nil, nil, nil
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
