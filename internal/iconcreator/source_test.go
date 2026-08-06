package iconcreator

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
)

func TestDecodeSourceWebP(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString("UklGRjIAAABXRUJQVlA4ICYAAACQAQCdASoCAAEAAgA0JYgCdLoAA5gA/us2/4oBSOHO8rH1KoAAAA==")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "source.webp")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	img, format, err := DecodeSource(path)
	if err != nil {
		t.Fatal(err)
	}
	if format != "webp" {
		t.Fatalf("format = %q, want webp", format)
	}
	if got := img.Bounds().Size(); got.X != 2 || got.Y != 1 {
		t.Fatalf("size = %dx%d, want 2x1", got.X, got.Y)
	}
	config, configFormat, err := DecodeSourceConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if configFormat != "webp" || config.Width != 2 || config.Height != 1 {
		t.Fatalf("config = %dx%d %q, want 2x1 webp", config.Width, config.Height, configFormat)
	}
}

func TestDecodeSourceSVG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source.svg")
	svg := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 100"><rect width="200" height="100" fill="#e3364a"/></svg>`
	if err := os.WriteFile(path, []byte(svg), 0600); err != nil {
		t.Fatal(err)
	}

	img, format, err := DecodeSource(path)
	if err != nil {
		t.Fatal(err)
	}
	if format != "svg" {
		t.Fatalf("format = %q, want svg", format)
	}
	if got := img.Bounds().Size(); got.X != 2048 || got.Y != 1024 {
		t.Fatalf("raster size = %dx%d, want 2048x1024", got.X, got.Y)
	}
	_, _, _, alpha := img.At(1024, 512).RGBA()
	if alpha == 0 {
		t.Fatal("expected the SVG rectangle to be rendered")
	}
	config, configFormat, err := DecodeSourceConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if configFormat != "svg" || config.Width != 2048 || config.Height != 1024 {
		t.Fatalf("config = %dx%d %q, want 2048x1024 svg", config.Width, config.Height, configFormat)
	}
}

func TestDecodeSourceSVGRequiresDimensions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "source.svg")
	if err := os.WriteFile(path, []byte(`<svg xmlns="http://www.w3.org/2000/svg"><circle r="10"/></svg>`), 0600); err != nil {
		t.Fatal(err)
	}

	if _, _, err := DecodeSource(path); err == nil {
		t.Fatal("expected an SVG without dimensions to fail")
	}
}
