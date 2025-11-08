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


// 画像のピクセルデータを取得
func getRgbaAt(img image.Image, x int, y int) (r uint8, g uint8, b uint8, a uint8) {
	c := img.At(x, y)
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return rgba.R, rgba.G, rgba.B, rgba.A
}

func LoadImage(path string) (ImageData, error) {
	// 画像ファイルを開く
	file, err := os.Open(path)
	if err != nil { return ImageData{}, err }
	defer file.Close()
	
	// 画像をデコード
	img, _, err := image.Decode(file)
	if err != nil { return ImageData{}, err }
	
	// 画像の幅と高さを取得
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	
	// 画像データを作成
	imageData := ImageData{
		Width: int64(w),
		Height: int64(h),
		Pixels: make([]PixelData, w * h),
	}
	
	// 画像データをピクセルデータに変換
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := getRgbaAt(img, x, y)
			imageData.SetPixel(int64(x), int64(y), PixelData{R: r, G: g, B: b, A: a})
		}
	}
	
	return imageData, nil
}
