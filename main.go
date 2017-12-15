package main

import (
	"fmt"

	"os"

	"encoding/hex"

	"github.com/clzhan/GoFlvParser/log"
	"github.com/clzhan/GoFlvParser/parser"
)

const FileNmae = "flv_test.flv"

func main() {

	//summary := pprof.GCSummary()

	fmt.Println("This is a flv parser  ")

	f, err := os.Open(FileNmae)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()

	fout, err1 := os.OpenFile("1.264", os.O_CREATE|os.O_TRUNC, 0666) //打开文件
	if err1 != nil {
		fmt.Println(err1.Error())
		return
	}
	defer fout.Close()

	log := log.NewLog(log.LogInfo, os.Stdout)

	flv := parser.InitReader(f, log)

	if _, err = flv.FlvReadHeader(); err != nil {
		log.Error("Read Flv Header error:%v", err.Error())
	}

	//flvHeader := &parser.FlvHeader{}
	//
	//if err = flvHeader.ParseFlvHeader(f); err != nil{
	//	fmt.Println("main err : " + err.Error())
	//	return
	//}
	//
	//for{
	//	tag := &parser.FlvTag{}
	//	if err= tag.ParseFlvTag(f);err != nil{
	//		fmt.Println("err : "+ err.Error())
	//		break;
	//	}
	//}
	//

	//var AdstHeader [7]byte
	var tag *parser.FlvTag
	var avc *parser.AVCDecoderConfigurationRecord
	var aac *parser.AudioSpecificConfig
	//var Ainfo  *parser.AudioSpecificConfig
	//var Vinfo  *parser.AVCDecoderConfigurationRecord
	var sps []byte
	var pps []byte
	var h264Prefix [4]byte
	h264Prefix[0] = 0x00;
	h264Prefix[1] = 0x00;
	h264Prefix[2] = 0x00;
	h264Prefix[3] = 0x01;

	for {
		tag, avc, aac, err = flv.FlvReadTag()
		if err != nil {
			fmt.Println("errt : " + err.Error())
			break
		}
		if avc != nil {
			//Vinfo = avc
			fmt.Println("avc..... ")

			for p := avc.SPS.Front(); p != nil; p = p.Next() {
				sps = p.Value.([]byte)
				fmt.Println("a sps...", hex.Dump(sps))
			}

			for p := avc.PPS.Front(); p != nil; p = p.Next() {
				pps = p.Value.([]byte)
				fmt.Println("a pps...", hex.Dump(pps))
			}

		}

		//output[0] = 0xFF;
		//output[1] = 0xF1; // 0xF9 (MPEG2)
		////output[2] = 0x40 | (GetSamplingFrequencyIndex(sampling_frequency) << 2) | (channel_configuration >> 2);
		//output[2] = 0x40 | (3 << 2) | (1 >> 2);
		////output[3] = ((channel_configuration & 0x3) << 6) | ((frame_size + 7) >> 11);
		//output[3] = ((1 & 0x3) << 6) | ((frame_size + 7) >> 11);
		//output[4] = ((frame_size + 7) >> 3) & 0xFF;
		//output[5] = (((frame_size + 7) << 5) & 0xFF) | 0x1F;
		//output[6] = 0xFC;


		if aac != nil {
			//Ainfo = aac
			fmt.Println("aac..... ")
		}

		if tag.Header.TagType == parser.FLVAudio {
			//fmt.Print("a tag.....  %p", tag)

			if tag.Buffer != nil {
				//fmt.Println("a auido...", hex.Dump(tag.Buffer))
				//AdstHeader[0] = 0xFF;
				//AdstHeader[1] = 0xF1; // 0xF9 (MPEG2);
				//AdstHeader[2] =  0x40 | (Ainfo.SampleRateIndex << 2) | (Ainfo.ChannelConfig >> 2);
				//var size int = len(tag.Buffer)
				//AdstHeader[3] =uint8((Ainfo.ChannelConfig  & 0x3) << 6) | uint8(( size + 7) >> 11);
				//AdstHeader[4] = uint8((size + 7) >> 3) & 0xFF;
				//AdstHeader[5] =  uint8(((size + 7) << 5) & 0xFF) | 0x1F;
				//AdstHeader[6] = 0xFC;
				//
				//fout.Write(AdstHeader[:])
				//fout.Write(tag.Buffer[:])
			}

		}

		if tag.Header.TagType == parser.FLVVideo{

			if tag.Buffer != nil {
				//fmt.Println("a video...", hex.Dump(tag.Buffer))
				fout.Write(h264Prefix[:])
				fout.Write(sps[:])
				fout.Write(h264Prefix[:])
				fout.Write(pps[:])
				fout.Write(h264Prefix[:])
				fout.Write(tag.Buffer[:])
			}

		}

	}

	fmt.Println("decode tag finished")

}
