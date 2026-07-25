package utils

import (
	"testing"
	"time"
)

func TestColonyObjectPaths(t *testing.T) {
	timestamp := time.Date(2026, 7, 18, 10, 51, 23, 0, time.UTC)
	tests := []struct {
		name       string
		sourcePath string
		wantImage  string
		wantRecord string
		wantError  bool
	}{
		{
			name:       "bmp",
			sourcePath: "/local/photo.bmp",
			wantImage:  "E8F60AE0CBEF/0/20260718-105123.bmp",
			wantRecord: "E8F60AE0CBEF/0/20260718-105123.txt",
		},
		{
			name:       "uppercase jpg",
			sourcePath: "/local/photo.JPG",
			wantImage:  "E8F60AE0CBEF/0/20260718-105123.jpg",
			wantRecord: "E8F60AE0CBEF/0/20260718-105123.txt",
		},
		{
			name:       "jpeg",
			sourcePath: "/local/photo.jpeg",
			wantImage:  "E8F60AE0CBEF/0/20260718-105123.jpeg",
			wantRecord: "E8F60AE0CBEF/0/20260718-105123.txt",
		},
		{
			name:       "surrounding whitespace",
			sourcePath: "  /local/photo.BMP  ",
			wantImage:  "E8F60AE0CBEF/0/20260718-105123.bmp",
			wantRecord: "E8F60AE0CBEF/0/20260718-105123.txt",
		},
		{name: "missing extension", sourcePath: "/local/photo", wantError: true},
		{name: "unsupported extension", sourcePath: "/local/photo.png", wantError: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			image, record, err := colonyObjectPaths("E8F60AE0CBEF", 0, timestamp, test.sourcePath)
			if test.wantError {
				if err == nil {
					t.Fatalf("colonyObjectPaths() error = nil, want an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("colonyObjectPaths() error = %v", err)
			}
			if image != test.wantImage {
				t.Errorf("image key = %q, want %q", image, test.wantImage)
			}
			if record != test.wantRecord {
				t.Errorf("record key = %q, want %q", record, test.wantRecord)
			}
		})
	}
}
