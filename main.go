package main

import (
	"fmt"
	"os"
	
	"ani-converter/loader"
	"ani-converter/convert"
)


func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ani-converter <image-file> <ani-file>")
		return
	}
	
	imageData, err := loader.LoadImage(os.Args[1])
	if err != nil { panic(err) }
	fmt.Printf("Image size: %d x %d\n", imageData.Width, imageData.Height)
	
	bytes, err := os.ReadFile(os.Args[2])
	if err != nil { panic(err) }
	riff := convert.Riff{}
	err = riff.Import(bytes)
	if err != nil { panic(err) }
	riff.Print()
	
	exportedBytes := riff.Export()
	err = os.WriteFile("img/export.ani", exportedBytes, 0644)
	if err != nil { panic(err) }
}
