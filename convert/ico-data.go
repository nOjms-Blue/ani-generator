package convert

import (
	"errors"
	"encoding/binary"
)


type IconFileHeader struct {
	IcoReserved      uint16
	IcoResourceType  uint16
	IcoResourceCount uint16
}

func (self *IconFileHeader) Import(bytes []byte) error {
	if len(bytes) < 6 {
		return errors.New("not enough bytes to import IconFileHeader")
	}
	self.IcoReserved = binary.LittleEndian.Uint16(bytes[:2])
	self.IcoResourceType = binary.LittleEndian.Uint16(bytes[2:4])
	self.IcoResourceCount = binary.LittleEndian.Uint16(bytes[4:6])
	return nil
}

func (self IconFileHeader) Export() []byte {
	var bytes []byte = make([]byte, 6)
	binary.LittleEndian.PutUint16(bytes[:2], self.IcoReserved)
	binary.LittleEndian.PutUint16(bytes[2:4], self.IcoResourceType)
	binary.LittleEndian.PutUint16(bytes[4:6], self.IcoResourceCount)
	return bytes
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

func (self *IconInfoHeader) Import(bytes []byte) error {
	if len(bytes) < 16 {
		return errors.New("not enough bytes to import IconInfoHeader")
	}
	self.Width = bytes[0]
	self.Height = bytes[1]
	self.ColorCount = bytes[2]
	self.Reserved1 = bytes[3]
	self.Reserved2 = binary.LittleEndian.Uint16(bytes[4:6])
	self.Reserved3 = binary.LittleEndian.Uint16(bytes[6:8])
	self.IcoDIBSize = binary.LittleEndian.Uint32(bytes[8:12])
	self.IcoDIBOffset = binary.LittleEndian.Uint32(bytes[12:16])
	return nil
}

func (self IconInfoHeader) Export() []byte {
	var bytes []byte = make([]byte, 16)
	bytes[0] = self.Width
	bytes[1] = self.Height
	bytes[2] = self.ColorCount
	bytes[3] = self.Reserved1
	binary.LittleEndian.PutUint16(bytes[4:6], self.Reserved2)
	binary.LittleEndian.PutUint16(bytes[6:8], self.Reserved3)
	binary.LittleEndian.PutUint32(bytes[8:12], self.IcoDIBSize)
	binary.LittleEndian.PutUint32(bytes[12:16], self.IcoDIBOffset)
	return bytes
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

func (self *BitmapInfoHeader) Import(bytes []byte) error {
	if len(bytes) < 40 {
		return errors.New("not enough bytes to import BitmapInfoHeader")
	}
	self.BcSize = binary.LittleEndian.Uint32(bytes[:4])
	self.BcWidth = binary.LittleEndian.Uint32(bytes[4:8])
	self.BcHeight = int32(binary.LittleEndian.Uint32(bytes[8:12]))
	self.BcPlanes = binary.LittleEndian.Uint16(bytes[12:14])
	self.BcBitCount = binary.LittleEndian.Uint16(bytes[14:16])
	self.BiCompression = binary.LittleEndian.Uint32(bytes[16:20])
	self.BiSizeImage = binary.LittleEndian.Uint32(bytes[20:24])
	self.BiXPixPerMeter = binary.LittleEndian.Uint32(bytes[24:28])
	self.BiYPixPerMeter = binary.LittleEndian.Uint32(bytes[28:32])
	self.BiClrUsed = binary.LittleEndian.Uint32(bytes[32:36])
	self.BiClrImportant = binary.LittleEndian.Uint32(bytes[36:40])
	return nil
}

func (self BitmapInfoHeader) Export() []byte {
	var bytes []byte = make([]byte, 40)
	binary.LittleEndian.PutUint32(bytes[:4], self.BcSize)
	binary.LittleEndian.PutUint32(bytes[4:8], self.BcWidth)
	binary.LittleEndian.PutUint32(bytes[8:12], uint32(self.BcHeight))
	binary.LittleEndian.PutUint16(bytes[12:14], self.BcPlanes)
	binary.LittleEndian.PutUint16(bytes[14:16], self.BcBitCount)
	binary.LittleEndian.PutUint32(bytes[16:20], self.BiCompression)
	binary.LittleEndian.PutUint32(bytes[20:24], self.BiSizeImage)
	binary.LittleEndian.PutUint32(bytes[24:28], self.BiXPixPerMeter)
	binary.LittleEndian.PutUint32(bytes[28:32], self.BiYPixPerMeter)
	binary.LittleEndian.PutUint32(bytes[32:36], self.BiClrUsed)
	binary.LittleEndian.PutUint32(bytes[36:40], self.BiClrImportant)
	return bytes
}

type Palette struct {
	Blue      uint8
	Green     uint8
	Red       uint8
	Reserved  uint8
}

func (self *Palette) Import(bytes []byte) error {
	if len(bytes) < 4 {
		return errors.New("not enough bytes to import Palette")
	}
	self.Blue = bytes[0]
	self.Green = bytes[1]
	self.Red = bytes[2]
	self.Reserved = bytes[3]
	return nil
}

func (self Palette) Export() []byte {
	var bytes []byte = []byte{self.Blue, self.Green, self.Red, self.Reserved}
	return bytes
}
