package convert

import (
	"fmt"
	"strings"
	"errors"
	"encoding/binary"
)


type Chunk struct {
	ChunkID [4]byte
	DataSize uint32
	SubChunks []int32 // index of SubChunks
	Data []byte
}

type Riff struct {
	Chunks []Chunk
}

type chunkReadData struct {
	Parent int32 // index of Parent
	RemainBytes []byte
}

func (self *Riff) Import(bytes []byte) error {
	var remains []chunkReadData = []chunkReadData{
		{Parent: -1, RemainBytes: bytes},
	}
	
	for len(remains) > 0 {
		remain := remains[0]
		remains = remains[1:]
		
		if remain.Parent < 0 {
			// RIFFチャンク(ルートチャンク)を読み込む
			chunk := Chunk{}
			chunk.ChunkID = [4]byte(remain.RemainBytes[:4])
			chunk.DataSize = uint32(binary.LittleEndian.Uint32(remain.RemainBytes[4:8]))
			if chunk.ChunkID == [4]byte{'R', 'I', 'F', 'F'} {
				chunk.Data = remain.RemainBytes[8:8+4]
				remains = append(remains, chunkReadData{Parent: 0, RemainBytes: remain.RemainBytes[8+4:]})
			} else {
				return errors.New("invalid chunk id")
			}
			self.Chunks = []Chunk{chunk}
		} else {
			// 親チャンクを取得
			parentIndex := remain.Parent
			parent := &self.Chunks[parentIndex]
			
			// 新しいチャンクを作成
			chunkIndex := (int32)(len(self.Chunks))
			chunk := Chunk{}
			chunk.ChunkID = [4]byte(remain.RemainBytes[:4])
			chunk.DataSize = uint32(binary.LittleEndian.Uint32(remain.RemainBytes[4:8]))
			if chunk.ChunkID == [4]byte{'L', 'I', 'S', 'T'} {
				chunk.Data = remain.RemainBytes[8:8+4]
				remains = append(remains, chunkReadData{Parent: chunkIndex, RemainBytes: remain.RemainBytes[8+4:8 + chunk.DataSize]})
			} else {
				chunk.Data = remain.RemainBytes[8:8 + chunk.DataSize]
			}
			parent.SubChunks = append(parent.SubChunks, chunkIndex)
			self.Chunks = append(self.Chunks, chunk)
			
			// 未使用バイト列があれば、次のチャンク読み取りを追加
			remain_bytes := remain.RemainBytes[8 + chunk.DataSize:]
			if len(remain_bytes) > 0 {
				remains = append(remains, chunkReadData{Parent: parentIndex, RemainBytes: remain_bytes})
			}
		}
	}
	return nil
}

func (self Riff) exportChunk(index int32) []byte {
	var bytes []byte = []byte{self.Chunks[index].ChunkID[0], self.Chunks[index].ChunkID[1], self.Chunks[index].ChunkID[2], self.Chunks[index].ChunkID[3], 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(bytes[4:8], uint32(self.Chunks[index].DataSize))
	bytes = append(bytes, self.Chunks[index].Data...)
	for _, subChunkIndex := range self.Chunks[index].SubChunks {
		bytes = append(bytes, self.exportChunk(subChunkIndex)...)
	}
	return bytes
}
func (self Riff) Export() []byte {
	var bytes []byte = []byte{}
	bytes = append(bytes, self.exportChunk(0)...)
	return bytes
}

// チャンクを再起的に標準出力へ
func (self Riff) printChunk(index int32, indent int32) {
	chunk := &self.Chunks[index]
	indentString := strings.Repeat("  ", int(indent))
	
	fmt.Printf("%s%c%c%c%c %d %d\n", indentString, chunk.ChunkID[0], chunk.ChunkID[1], chunk.ChunkID[2], chunk.ChunkID[3], len(chunk.Data), chunk.DataSize)
	for _, subChunkIndex := range chunk.SubChunks {
		self.printChunk(subChunkIndex, indent + 1)
	}
}

// チャンク情報を標準出力へ
func (self Riff) Print() {
	self.printChunk(0, 0)
}
