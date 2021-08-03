package validation

import "testing"

func TestValidateExtension(t *testing.T) {
	// Arrange

	supportedExtensions := []string{".png", ".jpeg", ".jpg"}
	tests := []struct {
		name      string
		extension string
		want      bool
	}{
		{
			name:      "PNG",
			extension: ".png",
			want:      true,
		},
		{
			name:      "JPEG",
			extension: ".jpeg",
			want:      true,
		},
		{
			name:      "JPG",
			extension: ".jpg",
			want:      true,
		},
		{
			name:      "WEBP",
			extension: ".webp",
			want:      false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// `t.Parallel()` will run all tests that include that line in parallel.
			// Any that don't include the line will run serially
			t.Parallel()

			// Act
			valid := validateExtension(testCase.extension, supportedExtensions)

			// Assert
			if valid != testCase.want {
				t.Errorf("wanted %t; got %t", testCase.want, valid)
			}
		})
	}
}

func TestValidateImageDimensions(t *testing.T) {
	// Arrange
	maxWidth := 500
	maxHeight := 500
	tests := []struct {
		name        string
		imageWidth  int
		imageHeight int
		want        bool
	}{
		{
			name:        "Valid Image",
			imageWidth:  500,
			imageHeight: 500,
			want:        true,
		},
		{
			name:        "Invalid Width",
			imageWidth:  1000,
			imageHeight: 500,
			want:        false,
		},
		{
			name:        "Invalid Height",
			imageWidth:  500,
			imageHeight: 1000,
			want:        false,
		},
		{
			name:        "Invalid Width & Height",
			imageWidth:  1000,
			imageHeight: 1000,
			want:        false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			img, _ := newTestImage(t, testCase.imageWidth, testCase.imageHeight)
			imageReader := convertImageToByteReader(t, img)

			// Act
			valid := validateImageDimensions(imageReader, maxWidth, maxHeight)

			// Assert
			if valid != testCase.want {
				t.Errorf("wanted %t; got %t", testCase.want, valid)
			}
		})
	}
}
