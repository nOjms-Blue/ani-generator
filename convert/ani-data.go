package convert


type AniHeader struct {
	Size        uint32
	Frames      uint32
	Steps       uint32
	Width       uint32
	Height      uint32
	BitDepth    uint32
	Planes      uint32
	DefaultRate uint32
	Flags       uint32
}

type Rate = uint32
type Sequence = uint32
