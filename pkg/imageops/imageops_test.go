package imageops

import "testing"

func TestResizeImage(t *testing.T) {
	// Arrange

	// If changing `originalWidth` or `originalHeight`, you will also need to update
	// the values for `expectedWidth` and `expectedHeight` in the following test cases:
	// "Only Height Specified", "Only Width Specified"
	originalWidth := 400
	originalHeight := 500
	_, testImageBytes := newTestImage(t, originalWidth, originalHeight)

	tests := []struct {
		name           string
		desiredWidth   int
		desiredHeight  int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "Reduce Size",
			desiredWidth:   200,
			desiredHeight:  200,
			expectedWidth:  200,
			expectedHeight: 200,
		},
		{
			name:           "Increase Size",
			desiredWidth:   700,
			desiredHeight:  700,
			expectedWidth:  700,
			expectedHeight: 700,
		},
		{
			name:           "Only Height Specified",
			desiredWidth:   0,
			desiredHeight:  300,
			expectedWidth:  240, // (desiredHeight / originalHeight) * originalWidth
			expectedHeight: 300,
		},
		{
			name:           "Only Width Specified",
			desiredWidth:   300,
			desiredHeight:  0,
			expectedWidth:  300,
			expectedHeight: 375, // (desiredWidth / originalWidth) * originalHeight
		},
		{
			name:           "Height & Width Zero",
			desiredWidth:   0,
			desiredHeight:  0,
			expectedWidth:  originalWidth,
			expectedHeight: originalHeight,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Act
			resizedImageBytes, _ := ResizeImage(testImageBytes, testCase.desiredHeight, testCase.desiredWidth)

			// Assert
			resizedImage, _ := DecodeImage(resizedImageBytes, "image/jpeg")

			if resizedImage.Bounds().Max.X != testCase.expectedWidth {
				t.Errorf("wanted width to equal %d, got %d", testCase.expectedWidth, resizedImage.Bounds().Max.X)
			}

			if resizedImage.Bounds().Max.Y != testCase.expectedHeight {
				t.Errorf("wanted height to equal %d, got %d", testCase.expectedHeight, resizedImage.Bounds().Max.Y)
			}
		})
	}
}
