package convert

import (
	"time"
	"errors"
	"unsafe"
	"encoding/binary"
	
	"ani-converter/loader"
)


type AniInfo struct {
	Name string  // INAM Chunk
	Artist string  // IART Chunk
	Copyright string  // ICOP Chunk
	Comment string  // ICMT Chunk
}

func ConvertToAni(resourceType ResourceType, images []loader.ImageData, hotSpotX []int16, hotSpotY []int16, frameIndexes []Sequence, rates []Rate, info *AniInfo) ([]byte, error) {
	// バリデーション
	if len(images) == 0 {
		return nil, errors.New("images must not be empty")
	}
	if len(frameIndexes) == 0 {
		return nil, errors.New("frameIndexes must not be empty")
	}
	if len(rates) == 0 {
		return nil, errors.New("rates must not be empty")
	}
	if len(images) != len(hotSpotX) || len(images) != len(hotSpotY) {
		return nil, errors.New("images and hotSpotX, hotSpotY must have the same length")
	}
	if len(frameIndexes) != len(rates) {
		return nil, errors.New("frameIndexes and rates must have the same length")
	}
	
	aniHeader := AniHeader{
		Size: uint32(unsafe.Sizeof(AniHeader{})),
		Frames: uint32(len(images)),
		Steps: uint32(len(frameIndexes)),
		Width: 0,
		Height: 0,
		BitDepth: 0,
		Planes: 0,
		DefaultRate: rates[0],
		Flags: 0b00000000000000000000000000000011,
	}
	
	imagesBytesSize := uint32(0)
	imagesBytes := make([][]byte, len(images))
	for i, image := range images {
		hsX := hotSpotX[i]
		hsY := hotSpotY[i]
		
		bytes, err := ConvertImageDataToIcon(image, resourceType, hsX, hsY)
		if err != nil { return nil, err }
		imagesBytes[i] = bytes
		imagesBytesSize += uint32(len(bytes))
	}
	
	ratesBytes := make([]byte, len(rates) * 4)
	for i, rate := range rates {
		binary.LittleEndian.PutUint32(ratesBytes[i*4:i*4+4], rate)
	}
	
	seqBytes := make([]byte, len(frameIndexes) * 4)
	for i, frameIndex := range frameIndexes {
		binary.LittleEndian.PutUint32(seqBytes[i*4:i*4+4], uint32(frameIndex))
	}
	
	anihSize := uint32(unsafe.Sizeof(aniHeader))
	listSize := 4 + uint32(len(images)) * 8 + imagesBytesSize
	rateSize := uint32(len(ratesBytes))
	seqSize := uint32(len(seqBytes))
	riffSize := (8 + anihSize) + (8 + listSize) + (8 + rateSize) + (8 + seqSize)
	riff := Riff{
		Chunks: []Chunk{
			{ ChunkID: [4]byte{'R', 'I', 'F', 'F'}, DataSize: riffSize, SubChunks: []int32{1, 2, 3, 4}, Data: []byte{'A', 'C', 'O', 'N'} },
			{ ChunkID: [4]byte{'a', 'n', 'i', 'h'}, DataSize: anihSize, SubChunks: []int32{}, Data: aniHeader.Export() },
			{ ChunkID: [4]byte{'L', 'I', 'S', 'T'}, DataSize: listSize, SubChunks: []int32{}, Data: []byte{'f', 'r', 'a', 'm'} },
			{ ChunkID: [4]byte{'r', 'a', 't', 'e'}, DataSize: rateSize, SubChunks: []int32{}, Data: ratesBytes },
			{ ChunkID: [4]byte{'s', 'e', 'q', ' '}, DataSize: seqSize, SubChunks: []int32{}, Data: seqBytes },
		},
	}
	for _, imageBytes := range imagesBytes {
		chunk := Chunk{
			ChunkID: [4]byte{'i', 'c', 'o', 'n'},
			DataSize: uint32(len(imageBytes)),
			SubChunks: []int32{},
			Data: imageBytes,
		}
		
		riff.Chunks = append(riff.Chunks, chunk)
		riff.Chunks[2].SubChunks = append(riff.Chunks[2].SubChunks, int32(len(riff.Chunks) - 1))
	}
	
	infoChunks := []Chunk{}
	if info != nil {
		if info.Name != "" {
			infoChunks = append(infoChunks, Chunk{
				ChunkID: [4]byte{'I', 'N', 'A', 'M'},
				DataSize: uint32(len(info.Name) + 1),
				SubChunks: []int32{},
				Data: []byte(info.Name + "\x00"),
			})
		}
		if info.Artist != "" {
			infoChunks = append(infoChunks, Chunk{
				ChunkID: [4]byte{'I', 'A', 'R', 'T'},
				DataSize: uint32(len(info.Artist) + 1),
				SubChunks: []int32{},
				Data: []byte(info.Artist + "\x00"),
			})
		}
		if info.Copyright != "" {
			infoChunks = append(infoChunks, Chunk{
				ChunkID: [4]byte{'I', 'C', 'O', 'P'},
				DataSize: uint32(len(info.Copyright) + 1),
				SubChunks: []int32{},
				Data: []byte(info.Copyright + "\x00"),
			})
		}
		if info.Comment != "" {
			infoChunks = append(infoChunks, Chunk{
				ChunkID: [4]byte{'I', 'C', 'M', 'T'},
				DataSize: uint32(len(info.Comment) + 1),
				SubChunks: []int32{},
				Data: []byte(info.Comment + "\x00"),
			})
		}
	}
	today := time.Now().Format("2006-01-02") + "\x00"
	softwareName := "AniConvert (temporary)" + "\x00"
	infoChunks = append(infoChunks, Chunk{
		ChunkID: [4]byte{'I', 'S', 'F', 'T'},
		DataSize: uint32(len(softwareName)),
		SubChunks: []int32{},
		Data: []byte(softwareName),
	})
	infoChunks = append(infoChunks, Chunk{
		ChunkID: [4]byte{'I', 'C', 'R', 'D'},
		DataSize: uint32(len(today)),
		SubChunks: []int32{},
		Data: []byte(today),
	})
	infoChunks = append(infoChunks, Chunk{
		ChunkID: [4]byte{'I', 'E', 'N', 'R'},
		DataSize: uint32(len("nOjms-Blue") + 1),
		SubChunks: []int32{},
		Data: []byte("nOjms-Blue" + "\x00"),
	})
	
	infoListSize := uint32(4)
	for _, infoChunk := range infoChunks {
		infoListSize += 8 + infoChunk.DataSize
	}
	riff.Chunks = append(riff.Chunks, Chunk{
		ChunkID: [4]byte{'L', 'I', 'S', 'T'},
		DataSize: infoListSize,
		SubChunks: []int32{},
		Data: []byte{'I', 'N', 'F', 'O'},
	})
	infoListChunkIndex := int32(len(riff.Chunks) - 1)
	riff.Chunks[0].SubChunks = append(riff.Chunks[0].SubChunks, infoListChunkIndex)
	for _, infoChunk := range infoChunks {
		riff.Chunks = append(riff.Chunks, infoChunk)
		riff.Chunks[infoListChunkIndex].SubChunks = append(riff.Chunks[infoListChunkIndex].SubChunks, int32(len(riff.Chunks) - 1))
	}
	return riff.Export(), nil
}