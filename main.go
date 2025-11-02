package main

import (
	"fmt"
	"os"
	
	"ani-converter/loader"
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ani-converter <image-file>")
		return
	}
	
	imageData, err := loader.LoadImage(os.Args[1])
	if err != nil { panic(err) }
	
	fmt.Printf("Image size: %d x %d\n", imageData.Width, imageData.Height)
}
