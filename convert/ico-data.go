package convert


type IconFileHeader struct {
	IcoReserved      uint16
	IcoResourceType  uint16
	IcoResourceCount uint16
}

type IconInfoHeader struct {
	Width            uint8
	Height           uint8
	ColorCount       uint8
	Reserved1        uint8
	Reserved2        uint16
	Reserved3        uint16
	IcoDIBSize       uint32
	IcoDIBOffset     uint32
}

type BitmapInfoHeader struct {
	BcSize           uint32
	BcWidth          uint32
	BcHeight         int32
	BcPlanes         uint16
	BcBitCount       uint16
	BiCompression    uint32
	BiSizeImage      uint32
	BiXPixPerMeter   uint32
	BiYPixPerMeter   uint32
	BiClrUsed        uint32
	BiClrImportant   uint32
}

type Palette struct {
	Blue      uint8
	Green     uint8
	Red       uint8
	Reserved  uint8
}
