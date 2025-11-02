package loader


type PixelData struct {
	R uint8
	G uint8
	B uint8
	A uint8
}

type ImageData struct {
	Width int64
	Height int64
	Pixels []PixelData
}

func (img ImageData) GetPixel(x int64, y int64) PixelData {
	return img.Pixels[y*img.Width+x]
}

func (img *ImageData) SetPixel(x int64, y int64, pixel PixelData) {
	img.Pixels[y*img.Width+x] = pixel
}
