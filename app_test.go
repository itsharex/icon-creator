package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreviewDataURLUsesInputMIMEType(t *testing.T) {
	path := filepath.Join(t.TempDir(), "preview")
	if err := os.WriteFile(path, []byte("preview-data"), 0600); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		format string
		prefix string
	}{
		{format: "webp", prefix: "data:image/webp;base64,"},
		{format: "svg", prefix: "data:image/svg+xml;base64,"},
	}
	for _, test := range tests {
		t.Run(test.format, func(t *testing.T) {
			got, err := previewDataURL(path, test.format)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.HasPrefix(got, test.prefix) {
				t.Fatalf("preview data URL = %q, want prefix %q", got, test.prefix)
			}
		})
	}
}
