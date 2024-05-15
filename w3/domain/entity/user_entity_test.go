package entity_test

import (
	"halosuster/domain/entity"
	"testing"
)

func TestIsImageURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "valid http jpg", url: "https://www.example.com/image.jpg", want: true},
		{name: "valid https png", url: "http://static.website.com/photo.png", want: true},
		{name: "invalid schema", url: "ftp://ftp.server.com/image.jpg", want: false},
		{name: "invalid extension", url: "https://www.example.com/document.pdf", want: false},
		{name: "invalid url format", url: "invalidurl", want: false},
		{name: "uppercase extension", url: "https://domain.com/image.GIF", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.ValidateIdentityCardScanImageURL(tt.url)
			if got != tt.want {
				t.Errorf("ValidateIdentityCardScanImageURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}
