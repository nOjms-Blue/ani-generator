package generator

import (
	"errors"
	"encoding/binary"
)


type AniHeader struct {
	Size        uint32  // ヘッダのサイズ
	Frames      uint32  // アイコン数
	Steps       uint32  // 表示コマ数
	Width       uint32  // 画像幅     (生データ時のみ)
	Height      uint32  // 画像高さ   (生データ時のみ)
	BitDepth    uint32  // ビット深度 (生データ時のみ)
	Planes      uint32  // プレーン数 (生データ時のみ)
	DefaultRate uint32  // コマの表示時間のデフォルト値
	Flags       uint32  // フラグ (bit0: アイコンorカーソルか  bit1: シーケンスデータを含むか)
}

func (self *AniHeader) Import(bytes []byte) error {
	if len(bytes) < 36 {
		return errors.New("not enough bytes to import AniHeader")
	}
	self.Size = binary.LittleEndian.Uint32(bytes[:4])
	self.Frames = binary.LittleEndian.Uint32(bytes[4:8])
	self.Steps = binary.LittleEndian.Uint32(bytes[8:12])
	self.Width = binary.LittleEndian.Uint32(bytes[12:16])
	self.Height = binary.LittleEndian.Uint32(bytes[16:20])
	self.BitDepth = binary.LittleEndian.Uint32(bytes[20:24])
	self.Planes = binary.LittleEndian.Uint32(bytes[24:28])
	self.DefaultRate = binary.LittleEndian.Uint32(bytes[28:32])
	self.Flags = binary.LittleEndian.Uint32(bytes[32:36])
	return nil
}

func (self AniHeader) Export() []byte {
	var bytes []byte = make([]byte, 36)
	binary.LittleEndian.PutUint32(bytes[:4], self.Size)
	binary.LittleEndian.PutUint32(bytes[4:8], self.Frames)
	binary.LittleEndian.PutUint32(bytes[8:12], self.Steps)
	binary.LittleEndian.PutUint32(bytes[12:16], self.Width)
	binary.LittleEndian.PutUint32(bytes[16:20], self.Height)
	binary.LittleEndian.PutUint32(bytes[20:24], self.BitDepth)
	binary.LittleEndian.PutUint32(bytes[24:28], self.Planes)
	binary.LittleEndian.PutUint32(bytes[28:32], self.DefaultRate)
	binary.LittleEndian.PutUint32(bytes[32:36], self.Flags)
	return bytes
}

type Rate = uint32
type Sequence = uint32
