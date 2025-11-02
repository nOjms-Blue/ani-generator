package loader

import (
	"os"
	"image"
	"image/color"
	
	_ "image/png"
	_ "image/jpeg"
	_ "image/gif"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/webp"
)


func getRgbAt(img image.Image, x int, y int) (r uint8, g uint8, b uint8, a uint8) {
	c := img.At(x, y)
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return rgba.R, rgba.G, rgba.B, rgba.A
}

func LoadImage(path string) (ImageData, error) {
	file, err := os.Open(path)
	if err != nil { return ImageData{}, err }
	defer file.Close()
	
	img, _, err := image.Decode(file)
	if err != nil { return ImageData{}, err }
	
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	
	
	imageData := ImageData{
		Width: int64(w),
		Height: int64(h),
		Pixels: make([]PixelData, w * h),
	}
	
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := getRgbAt(img, x, y)
			imageData.SetPixel(int64(x), int64(y), PixelData{R: r, G: g, B: b, A: a})
		}
	}
	
	return imageData, nil
}
