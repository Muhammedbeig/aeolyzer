package attachments

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	DefaultMaxFileBytes   = 10 << 20
	DefaultMaxTotalBytes  = 20 << 20
	DefaultMaxImagePixels = 25_000_000
	DefaultMaxPDFPages    = 100
)

var (
	ErrEmptyFile       = errors.New("attachment is empty")
	ErrFileTooLarge    = errors.New("attachment is too large")
	ErrInvalidFilename = errors.New("attachment filename is invalid")
	ErrUnsupportedFile = errors.New("attachment type is unsupported")
	ErrInvalidFile     = errors.New("attachment content is invalid")
)

type File struct {
	Name        string
	ContentType string
	Data        []byte
	Size        int64
	SHA256      [sha256.Size]byte
}

type Processor struct {
	MaxFileBytes   int64
	MaxImagePixels int64
	MaxPDFPages    int
}

func NewProcessor() *Processor {
	return &Processor{
		MaxFileBytes:   DefaultMaxFileBytes,
		MaxImagePixels: DefaultMaxImagePixels,
		MaxPDFPages:    DefaultMaxPDFPages,
	}
}

func (p *Processor) Process(name string, data []byte) (File, error) {
	if p == nil {
		return File{}, errors.New("attachment processor is nil")
	}
	cleanName, err := cleanFilename(name)
	if err != nil {
		return File{}, err
	}
	if len(data) == 0 {
		return File{}, ErrEmptyFile
	}
	if p.MaxFileBytes < 1 || int64(len(data)) > p.MaxFileBytes {
		return File{}, ErrFileTooLarge
	}

	contentType, err := p.validateContent(cleanName, data)
	if err != nil {
		return File{}, err
	}
	return File{
		Name:        cleanName,
		ContentType: contentType,
		Data:        bytes.Clone(data),
		Size:        int64(len(data)),
		SHA256:      sha256.Sum256(data),
	}, nil
}

func (p *Processor) validateContent(name string, data []byte) (string, error) {
	detected := strings.ToLower(strings.TrimSpace(strings.Split(http.DetectContentType(data), ";")[0]))
	switch detected {
	case "image/png", "image/jpeg", "image/gif":
		if err := p.validateImage(data, detected); err != nil {
			return "", err
		}
		return detected, nil
	case "application/pdf":
		if err := p.validatePDF(data); err != nil {
			return "", err
		}
		return detected, nil
	case "text/plain", "application/json":
		return validateText(name, data)
	default:
		if isTextExtension(filepath.Ext(name)) && utf8.Valid(data) && !bytes.ContainsRune(data, '\x00') {
			return validateText(name, data)
		}
		return "", ErrUnsupportedFile
	}
}

func (p *Processor) validateImage(data []byte, contentType string) error {
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("%w: image header", ErrInvalidFile)
	}
	expectedFormat := map[string]string{
		"image/png":  "png",
		"image/jpeg": "jpeg",
		"image/gif":  "gif",
	}[contentType]
	if format != expectedFormat || config.Width < 1 || config.Height < 1 {
		return fmt.Errorf("%w: image format", ErrInvalidFile)
	}
	pixels := int64(config.Width) * int64(config.Height)
	if p.MaxImagePixels < 1 || pixels > p.MaxImagePixels {
		return fmt.Errorf("%w: image dimensions", ErrInvalidFile)
	}
	if contentType == "image/gif" && bytes.Count(data, []byte{0x2c}) > 200 {
		return fmt.Errorf("%w: too many gif frames", ErrInvalidFile)
	}
	return nil
}

func (p *Processor) validatePDF(data []byte) error {
	if len(data) < 8 || !bytes.HasPrefix(data, []byte("%PDF-")) {
		return fmt.Errorf("%w: pdf header", ErrInvalidFile)
	}
	tailStart := max(len(data)-2048, 0)
	if !bytes.Contains(data[tailStart:], []byte("%%EOF")) {
		return fmt.Errorf("%w: pdf trailer", ErrInvalidFile)
	}
	lower := bytes.ToLower(data)
	for _, marker := range [][]byte{
		[]byte("/javascript"),
		[]byte("/launch"),
		[]byte("/embeddedfile"),
		[]byte("/xfa"),
	} {
		if bytes.Contains(lower, marker) {
			return fmt.Errorf("%w: active pdf content", ErrInvalidFile)
		}
	}
	pageCount := bytes.Count(lower, []byte("/type/page"))
	if p.MaxPDFPages < 1 || pageCount > p.MaxPDFPages {
		return fmt.Errorf("%w: pdf page limit", ErrInvalidFile)
	}
	return nil
}

func validateText(name string, data []byte) (string, error) {
	if !utf8.Valid(data) || bytes.ContainsRune(data, '\x00') {
		return "", fmt.Errorf("%w: text encoding", ErrInvalidFile)
	}
	switch strings.ToLower(filepath.Ext(name)) {
	case ".json":
		if !json.Valid(data) {
			return "", fmt.Errorf("%w: json syntax", ErrInvalidFile)
		}
		return "application/json", nil
	case ".csv":
		return "text/csv", nil
	case ".md", ".markdown":
		return "text/markdown", nil
	case ".html", ".htm":
		return "text/html", nil
	default:
		return "text/plain", nil
	}
}

func cleanFilename(name string) (string, error) {
	if name == "" || filepath.IsAbs(name) || filepath.Base(name) != name || len(name) > 255 {
		return "", ErrInvalidFilename
	}
	cleaned := strings.TrimSpace(name)
	if cleaned == "" || cleaned == "." {
		return "", ErrInvalidFilename
	}
	for _, character := range cleaned {
		if unicode.IsControl(character) {
			return "", ErrInvalidFilename
		}
	}
	return cleaned, nil
}

func isTextExtension(extension string) bool {
	switch strings.ToLower(extension) {
	case ".txt", ".md", ".markdown", ".csv", ".json", ".html", ".htm", ".css", ".js", ".ts", ".tsx", ".jsx", ".go", ".py", ".java", ".rs", ".xml", ".yaml", ".yml":
		return true
	default:
		return false
	}
}
