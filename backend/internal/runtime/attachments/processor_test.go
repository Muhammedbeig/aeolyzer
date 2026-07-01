package attachments

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/png"
	"testing"
)

func TestProcessorProcess(t *testing.T) {
	t.Parallel()

	processor := NewProcessor()
	pngData := encodePNG(t, 2, 3)
	pdfData := []byte("%PDF-1.7\n1 0 obj <</Type/Page>> endobj\n%%EOF")

	tests := []struct {
		name        string
		filename    string
		data        []byte
		contentType string
		wantErr     error
	}{
		{name: "png", filename: "chart.png", data: pngData, contentType: "image/png"},
		{name: "pdf", filename: "report.pdf", data: pdfData, contentType: "application/pdf"},
		{name: "markdown", filename: "notes.md", data: []byte("# Notes"), contentType: "text/markdown"},
		{name: "json", filename: "data.json", data: []byte(`{"ok":true}`), contentType: "application/json"},
		{name: "invalid json", filename: "data.json", data: []byte(`{`), wantErr: ErrInvalidFile},
		{name: "active pdf", filename: "bad.pdf", data: []byte("%PDF-1.7\n/JavaScript\n%%EOF"), wantErr: ErrInvalidFile},
		{name: "path", filename: "../secret.txt", data: []byte("x"), wantErr: ErrInvalidFilename},
		{name: "unsupported", filename: "archive.zip", data: []byte("PK\x03\x04payload"), wantErr: ErrUnsupportedFile},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := processor.Process(test.filename, test.data)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Process() error = %v, want %v", err, test.wantErr)
			}
			if err == nil && got.ContentType != test.contentType {
				t.Fatalf("Process() content type = %q, want %q", got.ContentType, test.contentType)
			}
		})
	}
}

func TestProcessorRejectsImagePixelBomb(t *testing.T) {
	t.Parallel()

	processor := NewProcessor()
	processor.MaxImagePixels = 100
	data := encodePNG(t, 11, 10)
	_, err := processor.Process("large.png", data)
	if !errors.Is(err, ErrInvalidFile) {
		t.Fatalf("Process() error = %v, want %v", err, ErrInvalidFile)
	}
}

func TestProcessorRejectsOversizedFile(t *testing.T) {
	t.Parallel()

	processor := NewProcessor()
	processor.MaxFileBytes = 3
	_, err := processor.Process("file.txt", []byte("four"))
	if !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("Process() error = %v, want %v", err, ErrFileTooLarge)
	}
}

func encodePNG(t *testing.T, width, height int) []byte {
	t.Helper()
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, image.NewRGBA(image.Rect(0, 0, width, height))); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}

func TestMalformedPNGDimensionsAreRejected(t *testing.T) {
	t.Parallel()

	data := encodePNG(t, 1, 1)
	binary.BigEndian.PutUint32(data[16:20], 100_000)
	processor := NewProcessor()
	_, err := processor.Process("bomb.png", data)
	if !errors.Is(err, ErrInvalidFile) {
		t.Fatalf("Process() error = %v, want %v", err, ErrInvalidFile)
	}
}
