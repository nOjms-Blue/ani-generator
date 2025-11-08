package convert

import (
	"errors"
	"unsafe"
	"ani-converter/loader"
)


type ResourceType = uint16
const (
	IconResource ResourceType = 1
	CursorResource ResourceType = 2
)

func ConvertImageDataToIcon(image loader.ImageData, resourceType ResourceType) ([]byte, error) {
	if image.Width != 32 && image.Width != 64 && image.Width != 128 && image.Width != 256 {
		return nil, errors.New("unsupported image width")
	}
	if image.Height != 32 && image.Height != 64 && image.Height != 128 && image.Height != 256 {
		return nil, errors.New("unsupported image height")
	}
	
	iconFileHeader := IconFileHeader{
		IcoReserved: 0,
		IcoResourceType: resourceType,
		IcoResourceCount: 1,
	}
	iconInfoHeader := IconInfoHeader{
		Width: uint8(image.Width),
		Height: uint8(image.Height),
		ColorCount: 0,
		Reserved1: 0,
		Reserved2: 0,
		Reserved3: 0,
		IcoDIBSize: 0,
		IcoDIBOffset: uint32(unsafe.Sizeof(IconFileHeader{}) + unsafe.Sizeof(IconInfoHeader{})),
	}
	bitmapInfoHeader := BitmapInfoHeader{
		BcSize: uint32(unsafe.Sizeof(BitmapInfoHeader{})),
		BcWidth: uint32(image.Width),
		BcHeight: int32(image.Height) * 2,
		BcPlanes: 1,
		BcBitCount: 32,
		BiCompression: 0,
		BiSizeImage: 0,
		BiXPixPerMeter: 0,
		BiYPixPerMeter: 0,
		BiClrUsed: 0,
		BiClrImportant: 0,
	}
	
	pixels := make([]byte, image.Width * image.Height * 4)
	for y := image.Height - 1; y >= 0; y-- {
		for x := int64(0); x < image.Width; x++ {
			index := ((image.Height - 1 - y) * image.Width + x) * 4
			rgba := image.GetPixel(x, y)
			
			pixels[index + 0] = rgba.B
			pixels[index + 1] = rgba.G
			pixels[index + 2] = rgba.R
			pixels[index + 3] = rgba.A
		}
	}
	i := 0
	alphaMasks := make([]byte, image.Width / 8 * image.Height)
	for y := image.Height - 1; y >= 0; y-- {
		mask := uint8(0b10000000)
		value := uint8(0)
		for x := int64(0); x < image.Width; x++ {
			rgba := image.GetPixel(x, y)
			if rgba.A == 0 {
				value = value | mask
			}
			mask = mask >> 1
			if mask == 0 {
				alphaMasks[i] = value
				i++
				value = 0
				mask = 0b10000000
			}
		}
	}
	
	iconInfoHeader.IcoDIBSize = uint32(int64(unsafe.Sizeof(BitmapInfoHeader{})) + int64(len(pixels)) + int64(len(alphaMasks)))
	bitmapInfoHeader.BiSizeImage = uint32(len(pixels)) + uint32(len(alphaMasks))
	
	bytes := make([]byte, 0)
	bytes = append(bytes, iconFileHeader.Export()...)
	bytes = append(bytes, iconInfoHeader.Export()...)
	bytes = append(bytes, bitmapInfoHeader.Export()...)
	bytes = append(bytes, pixels...)
	bytes = append(bytes, alphaMasks...)
	
	return bytes, nil
}
