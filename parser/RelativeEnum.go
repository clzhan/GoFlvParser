package parser

type EnumPreProcessing int

func (v EnumPreProcessing) String() (s string) {
	switch v {
	case 0:
		s = "NonPrepProcessing"
	case 1:
		s = "PreProcessing"

	}
	return
}

type EnumTagType int

func (v EnumTagType) String() (s string) {
	switch v {
	case 8:
		s = "audio"
	case 9:
		s = "video"
	case 18:
		s = "script"
	}
	return
}

func GetSoundFormat(v uint8) (s string) {
	switch v {
	case 1:
		s = "ADPCM"
	case 2:
		s = "MP3"
	case 10:
		s = "AAC"
	case 11:
		s = "Speex"
	}
	return s
}

func GetSoundRate(v uint8) (s string) {
	switch v {
	case 0:
		s = "5.5kHz"
	case 1:
		s = "11kHz"
	case 2:
		s = "22kHz"
	case 3:
		s = "44kHz"
	}
	return s
}

func GetSoundSize(v uint8) (s string) {
	switch v {
	case 0:
		s = "8bit"
	case 1:
		s = "16bit"
	}
	return s
}

func GetSoundType(v uint8) (s string) {
	switch v {
	case 0:
		s = "Mono"
	case 1:
		s = "Stereo"
	}
	return s
}

// if sound format = 10
func GetAACPacketType(v uint8) (s string) {
	switch v {
	case 0:
		s = "Sequence header"
	case 1:
		s = "raw"
	}
	return s
}


func GetAACProfile(v uint8) (s string) {
	switch v {
	case 1:
		s = "AAC Main"
	case 2:
		s = "AAC LC"
	case 3:
		s = "AAC SSR"
	}
	return s
}

func GetSampleRate(v uint8) (s string) {
	switch v {
	case 0:
		s = "96000"
	case 1:
		s = "88200"
	case 2:
		s = "64000"
	case 3:
		s = "64000"
	case 4:
		s = "44100"
	case 5:
		s = "32000"
	case 6:
		s = "24000"
	case 7:
		s = "22050"
	case 8:
		s = "16000"
	case 9:
		s = "12000"
	case 10:
		s = "11025"
	case 11:
		s = "8000"
	case 12:
		s = "reserved"
	case 13:
		s = "reserved"
	case 14:
		s = "reserved"
	case 15:
		s = "escape value"
	}
	return s
}

func GetChannel(v uint8) (s string) {
	switch v {
	case 1:
		s = "单声道" //center front speaker
	case 2:
		s = "双声道" //left, right front speakers
	case 3:
		s = "三声道" //center, left, right front speakers
	case 4:
		s = "四声道" //center, left, right front speakers, rear surround speakers
	case 5:
		s = "五声道" //center, left, right front speakers, left surround, right surround rear speakers
	case 6:
		s = "5.1声道" //center, left, right front speakers, left surround, right surround rear speakers, front low frequency effects speaker
	case 7:
		s = "7.1声道" //center, left, right center front speakers, left, right outside front speakers, left surround, right surround rear speakers, front low frequency effects speaker
	default:
		s = "reserved" //0x08-0x0F - reserved
	}
	return s
}

func GetFrameType(v uint8) (s string) {
	switch v {
	case 1:
		s = "key frame"
	case 2:
		s = "inter frame"
	case 3:
		s = "disposable inter frame"
	case 4:
		s = "generated key frame"
	case 5:
		s = "video info/command frame"
	}
	return s

}

func GetCodecId(v uint8) (s string) {

	switch v {
	case 2:
		s = "H263"
	case 3:
		s = "Screen video"
	case 4:
		s = "VP6"
	case 5:
		s = "VP6 with alpha channel"
	case 6:
		s = "Screen video 2"
	case 7:
		s = "AVC"
	}
	return s
}


func GetAVCPacketType(v uint8) (s string) {

	switch v {
	case 0:
		s = "Sequence header"
	case 1:
		s = "NALU"
	case 2:
		s = "End of sequence"
	}
	return s
}

