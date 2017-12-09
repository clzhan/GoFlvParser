package parser

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type FlvHeader struct {
	SignatrueF         uint8  //1 byte  F
	SignatrueL         uint8  //1 byte  L
	SignatrueV         uint8  //1 byte  V
	Version            uint8  //1 byte 版本号
	TypeFlagsReserved5 uint8  //5 bit 保留的， 必须为0
	TypeFlagsAudio     uint8  //1bit 是否存在音频tag
	TypeFlagsReserved1 uint8  //1bit 保留 为0
	TypeFlagsVideo     uint8  //1bit  是否存在视频tag
	DataOffset         uint32 //4byte flv头尾长度 version1中一般为9
}


func (header *FlvHeader) ParseFlvHeader(r io.Reader) (err error) {
	//1. read signatrue 签名

	tmp := [5]uint8{}
	if err := binary.Read(r, binary.BigEndian, &tmp); err != nil {
		fmt.Println("read flv heard failed, err is " + err.Error())
		return err
	}
	//b := make([]byte, len(tmp))
	//for i, v := range tmp {
	//	b[i] = byte(v)
	//}
	//fmt.Println(hex.Dump(b))
	if string(tmp[0:3]) != "FLV" {
		fmt.Println("exp flv header=FLV, but actual = %v", string(tmp[0:3]))
		return errors.New("No header")
	}
	header.SignatrueF = tmp[0]
	header.SignatrueL = tmp[1]
	header.SignatrueV = tmp[2]

	header.Version = tmp[3]
	//fmt.Println("flv heard Version  is %v " , header.Version)

	header.TypeFlagsReserved5 = tmp[4] >> 3 & 0x1f

	header.TypeFlagsAudio = tmp[4] >> 2 & 0x01
	header.TypeFlagsReserved1 = tmp[4] >> 1 & 0x01
	header.TypeFlagsVideo = tmp[4] & 0x01

	binary.Read(r, binary.BigEndian, &header.DataOffset)

	fmt.Printf("==========================FLV Header==========================\n")
	fmt.Printf("SignatureF:		%c\n", header.SignatrueF)
	fmt.Printf("SignatureL:		%c\n", header.SignatrueL)
	fmt.Printf("SignatureV:		%c\n", header.SignatrueV)
	fmt.Printf("Version:	    0x%x\n", header.Version)
	fmt.Printf("TypeFlagsReserved5:	0x%x\n", header.TypeFlagsReserved5)
	fmt.Printf("TypeFlagsAudio:	0x%x\n", header.TypeFlagsAudio)
	fmt.Printf("TypeFlagsVideo:	0x%x\n", header.TypeFlagsVideo)
	fmt.Printf("DataOffset:		0x%x\n", header.DataOffset)
	fmt.Printf("==========================FLV Body============================\n")
	return nil

}
