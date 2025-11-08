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
	
	icoCount := 0
	checkBytes := []byte{}
	for _, chunk := range riff.Chunks {
		if chunk.ChunkID == [4]byte{'i', 'c', 'o', 'n'} {
			err = os.WriteFile(fmt.Sprintf("img/icon-%d.ico", icoCount + 1), chunk.Data, 0644)
			if err != nil { panic(err) }
			
			if icoCount == 0 {
				checkBytes = make([]byte, len(chunk.Data))
				copy(checkBytes, chunk.Data)
			}
			icoCount++
		}
	}
	
	bytes, err = convert.ConvertImageDataToIcon(imageData, convert.CursorResource, 16, 16)
	if err != nil { panic(err) }
	err = os.WriteFile("img/icon.ico", bytes, 0644)
	if err != nil { panic(err) }
}
