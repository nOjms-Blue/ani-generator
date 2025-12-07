package main

import (
	"fmt"
	"os"
	"strconv"
	
	"ani-generator/settings"
	"ani-generator/loader"
	"ani-generator/generator"
)


const VERSION = "1.0.0"

func main() {
	var images []loader.ImageData
	var hotSpotXs []int16
	var hotSpotYs []int16
	var frameIndexes []uint32
	var rates []uint32
	var outputAniPath string
	
	fmt.Printf("ani-generator %s\n", VERSION)
	fmt.Println("Copyright (c) 2025 nOjms-Blue")
	fmt.Println("https://github.com/nOjms-Blue/ani-generator")
	fmt.Println()
	
	optionSettings := settings.OptionSettings{
		KeySettings: []settings.OptionKeySetting{
			{ LongKey: "input", ShortKey: "i", Required: false, Flag: false, Multiple: true, Description: "input image files", Example: "--input input.png" },
			{ LongKey: "output", ShortKey: "o", Required: false, Flag: false, Multiple: false, Description: "output ani file", Example: "--output output.ani" },
			{ LongKey: "hotSpotX", ShortKey: "hx", Required: false, Flag: false, Multiple: true, Description: "hot spot x coordinates", Example: "--hotSpotX 10 --hotSpotX 12" },
			{ LongKey: "hotSpotY", ShortKey: "hy", Required: false, Flag: false, Multiple: true, Description: "hot spot y coordinates", Example: "--hotSpotY 10 --hotSpotY 12" },
			{ LongKey: "frameIndex", ShortKey: "f", Required: false, Flag: false, Multiple: true, Description: "frame indexes", Example: "--frameIndex 0 --frameIndex 1" },
			{ LongKey: "rate", ShortKey: "r", Required: false, Flag: false, Multiple: true, Description: "frame rates (1/60s)", Example: "--rate 10 --rate 12" },
			{ LongKey: "json", ShortKey: "", Required: false, Flag: false, Multiple: false, Description: "apply settings from json file", Example: "--json settings.json" },
		},
		NoKeySettings: []settings.OptionNoKeySetting{},
	}
	checkResult, err := settings.CheckArguments(os.Args[1:], optionSettings)
	if err != nil {
		fmt.Println("Error: ", err)
		settings.PrintHelp(optionSettings)
		os.Exit(0)
	}
	
	if json, ok := checkResult.KeyValues["json"]; ok {
		settings, err := settings.LoadSettingsJson(json[0])
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		
		for _, imageSetting := range settings.Images {
			// 入力画像の読み込み
			image, err := loader.LoadImage(imageSetting.Path)
			if err != nil {
				fmt.Println("Error: ", err)
				os.Exit(1)
			}
			images = append(images, image)
			
			// ホットスポットの読み込み
			hotSpotXs = append(hotSpotXs, imageSetting.HotSpotX)
			hotSpotYs = append(hotSpotYs, imageSetting.HotSpotY)
		}
		frameIndexes = settings.FrameIndexes
		rates = settings.Rates
		outputAniPath = settings.Output
	} else {
		// 入力画像の読み込み
		inputImages, ok := checkResult.KeyValues["input"]
		if !ok {
			fmt.Println("Error: input images are not set")
			settings.PrintHelp(optionSettings)
			os.Exit(0)
		}
		images = []loader.ImageData{}
		for _, inputImage := range inputImages {
			image, err := loader.LoadImage(inputImage)
			if err != nil {
				fmt.Println("Error: ", err)
				settings.PrintHelp(optionSettings)
				os.Exit(0)
			}
			images = append(images, image)
		}
		
		// 出力先パスの読み込み
		outputAnis, ok := checkResult.KeyValues["output"]
		if !ok {
			fmt.Println("Error: output ani file is not set")
			settings.PrintHelp(optionSettings)
			os.Exit(0)
		}
		if len(outputAnis) == 0 {
			fmt.Println("Error: output ani file is not set")
			settings.PrintHelp(optionSettings)
			os.Exit(0)
		}
		outputAniPath = outputAnis[0]
		
		// ホットスポットの読み込み
		hotSpotXStrings, ok := checkResult.KeyValues["hotSpotX"]
		if !ok {
			fmt.Println("Error: hot spot x coordinates are not set")
			settings.PrintHelp(optionSettings)
			os.Exit(0)
		}
		hotSpotYStrings, ok := checkResult.KeyValues["hotSpotY"]
		if !ok {
			fmt.Println("Error: hot spot y coordinates are not set")
			settings.PrintHelp(optionSettings)
			os.Exit(0)
		}
		hotSpotXs = []int16{}
		for _, hotSpotXString := range hotSpotXStrings {
			hotSpotX, err := strconv.ParseInt(hotSpotXString, 10, 16)
			if err != nil {
				fmt.Println("Error: ", err)
				settings.PrintHelp(optionSettings)
				os.Exit(0)
			}
			hotSpotXs = append(hotSpotXs, int16(hotSpotX))
		}
		hotSpotYs = []int16{}
		for _, hotSpotYString := range hotSpotYStrings {
			hotSpotY, err := strconv.ParseInt(hotSpotYString, 10, 16)
			if err != nil {
				fmt.Println("Error: ", err)
				settings.PrintHelp(optionSettings)
				os.Exit(0)
			}
			hotSpotYs = append(hotSpotYs, int16(hotSpotY))
		}
		
		// フレームインデックスの読み込み
		frameIndexStrings, ok := checkResult.KeyValues["frameIndex"]
		if !ok {
			fmt.Println("Error: frame indexes are not set")
			settings.PrintHelp(optionSettings)
			os.Exit(0)
		}
		frameIndexes = []uint32{}
		for _, frameIndexString := range frameIndexStrings {
			frameIndex, err := strconv.ParseInt(frameIndexString, 10, 32)
			if err != nil {
				fmt.Println("Error: ", err)
				settings.PrintHelp(optionSettings)
				os.Exit(0)
			}
			frameIndexes = append(frameIndexes, uint32(frameIndex))
		}
		
		// フレームレートの読み込み
		rateStrings, ok := checkResult.KeyValues["rate"]
		if !ok {
			fmt.Println("Error: frame rates are not set")
			settings.PrintHelp(optionSettings)
			os.Exit(0)
		}
		rates = []uint32{}
		for _, rateString := range rateStrings {
			rate, err := strconv.ParseUint(rateString, 10, 32)
			if err != nil {
				fmt.Println("Error: ", err)
				settings.PrintHelp(optionSettings)
				os.Exit(0)
			}
			rates = append(rates, uint32(rate))
		}
	}
	
	// アニメーションの作成
	ani, err := generator.ConvertToAni(generator.CursorResource, images, hotSpotXs, hotSpotYs, frameIndexes, rates)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	
	// アニメーションの保存
	err = os.WriteFile(outputAniPath, ani, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(2)
	}
	fmt.Println("Success: animation saved to ", outputAniPath)
}
