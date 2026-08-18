package generator

import (
	"encoding/binary"
	"strings"
	"testing"

	"ani-generator/loader"
)

func validImage() loader.ImageData {
	return loader.ImageData{Width: 32, Height: 32, Pixels: make([]loader.PixelData, 32*32)}
}

func TestConvertImageDataToIconRejectsMalformedInput(t *testing.T) {
	tests := []struct {
		name     string
		image    loader.ImageData
		resource ResourceType
		x, y     int16
		want     string
	}{
		{"invalid resource", validImage(), 99, 0, 0, "invalid resource type"},
		{"short pixel buffer", loader.ImageData{Width: 32, Height: 32}, IconResource, 0, 0, "invalid pixel count"},
		{"negative hotspot", validImage(), CursorResource, -1, 0, "outside the image"},
		{"large hotspot", validImage(), CursorResource, 0, 32, "outside the image"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ConvertImageDataToIcon(tt.image, tt.resource, tt.x, tt.y)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want error containing %q", err, tt.want)
			}
		})
	}
}

func TestConvertToAniRejectsOutOfRangeSequence(t *testing.T) {
	_, err := ConvertToAni(IconResource, []loader.ImageData{validImage()}, []int16{0}, []int16{0}, []Sequence{1}, []Rate{1})
	if err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("error = %v, want out-of-range error", err)
	}
}

func TestRiffImportRejectsMalformedDataWithoutPanicking(t *testing.T) {
	malformed := [][]byte{
		nil,
		[]byte("RIFF"),
		append([]byte("NOPE\x04\x00\x00\x00"), []byte("ACON")...),
	}

	// A valid root containing a child whose declared payload is truncated.
	truncatedChild := []byte("RIFF\x0c\x00\x00\x00ACONrate\x08\x00\x00\x00")
	malformed = append(malformed, truncatedChild)

	for i, data := range malformed {
		var riff Riff
		if err := riff.Import(data); err == nil {
			t.Errorf("case %d: expected an error", i)
		}
	}
}

func TestRiffImportHandlesOddChunkPadding(t *testing.T) {
	data := []byte("RIFF\x16\x00\x00\x00ACONone \x01\x00\x00\x00x\x00two \x00\x00\x00\x00")
	// Keep this fixture readable while asserting its RIFF size is correct.
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(data)-8))

	var riff Riff
	if err := riff.Import(data); err != nil {
		t.Fatalf("Import() error = %v", err)
	}
	if len(riff.Chunks) != 3 {
		t.Fatalf("got %d chunks, want 3", len(riff.Chunks))
	}
}
