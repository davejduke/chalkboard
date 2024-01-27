package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"

	fatihColor "github.com/fatih/color"

	"golang.org/x/image/draw"
)

func main() {
	fmt.Print("CHALKBOARDIMAGE: by Dave Duke ")
	fatihColor.New(fatihColor.FgHiGreen).Print("dave@daveduke.co.uk")
	fmt.Println()

	// Check if an image file path is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_image_path> [threshold] [thickness] [-invert]")
		return
	}

	// Get the input image file path from the command-line argument
	inputFileName := os.Args[1]

	// Process the filename for spaces (if any) by using filepath.Glob
	files, err := filepath.Glob(inputFileName)
	if err != nil {
		log.Fatalf("Error processing file path: %s", err)
	}

	if len(files) < 1 {
		fmt.Println("No image files found.")
		return
	}

	inputFileName = files[0]

	// Set default threshold and thickness values
	threshold := 20
	thickness := 0

	// Check if threshold and thickness are provided as command-line arguments
	if len(os.Args) >= 4 {
		// Parse the threshold value
		threshold, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid threshold value. Using default (50).")
			threshold = 50
		}

		// Parse the thickness value
		thickness, err = strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Invalid thickness value. Using default (1).")
			thickness = 1
		}
	}

	// Check if the -invert flag is provided
	invert := false
	for _, arg := range os.Args[1:] {
		if arg == "-invert" {
			invert = true
			break
		}
	}

	// Replace "output.png" with the path where you want to save the output mask.
	outputFileName := "output.png"

	// Read the input image file
	inputImage, err := readImage(inputFileName)
	if err != nil {
		log.Fatalf("Error reading image: %s", err)
	}

	// Create a grayscale version of the input image (if needed)
	grayImage := image.NewGray(inputImage.Bounds())
	draw.Draw(grayImage, grayImage.Bounds(), inputImage, image.Point{}, draw.Src)

	// Apply edge detection or contour finding algorithm here to create the mask
	mask := createMask(grayImage, threshold, thickness)

	// Invert the mask if the -invert flag is provided
	if invert {
		invertImage(mask)
	}

	// Save the mask to an output file
	if err := saveMask(outputFileName, mask); err != nil {
		log.Fatalf("Error saving mask: %s", err)
	}

	fmt.Println("Mask created and saved successfully.")
}

// readImage reads an image file and returns the decoded image.
func readImage(fileName string) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// createMask applies an edge detection or contour finding algorithm to create a mask.
// Replace this function with your own image processing algorithm to generate the mask.
func createMask(img *image.Gray, threshold, thickness int) *image.Gray {
	bounds := img.Bounds()
	mask := image.NewGray(bounds)

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			// Calculate the gradient in the x and y directions
			dx := int(img.GrayAt(x+1, y).Y) - int(img.GrayAt(x-1, y).Y)
			dy := int(img.GrayAt(x, y+1).Y) - int(img.GrayAt(x, y-1).Y)

			// Calculate the gradient magnitude
			gradient := uint8(math.Sqrt(float64(dx*dx + dy*dy)))

			// Check if the gradient is greater than the threshold or is larger than its neighbors
			if int(gradient) > threshold ||
				int(gradient) > int(img.GrayAt(x-1, y-1).Y)+threshold ||
				int(gradient) > int(img.GrayAt(x, y-1).Y)+threshold ||
				int(gradient) > int(img.GrayAt(x+1, y-1).Y)+threshold ||
				int(gradient) > int(img.GrayAt(x-1, y).Y)+threshold ||
				int(gradient) > int(img.GrayAt(x+1, y).Y)+threshold ||
				int(gradient) > int(img.GrayAt(x-1, y+1).Y)+threshold ||
				int(gradient) > int(img.GrayAt(x, y+1).Y)+threshold ||
				int(gradient) > int(img.GrayAt(x+1, y+1).Y)+threshold {

				for i := y - thickness; i <= y+thickness; i++ {
					for j := x - thickness; j <= x+thickness; j++ {
						if i >= bounds.Min.Y && i < bounds.Max.Y && j >= bounds.Min.X && j < bounds.Max.X {
							// Increase the pixel intensity to make the lines brighter
							mask.SetGray(j, i, color.Gray{Y: 255})
						}
					}
				}
			}
		}
	}

	return mask
}

// saveMask saves the mask image to the specified file.
func saveMask(fileName string, mask *image.Gray) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := png.Encode(file, mask); err != nil {
		return err
	}

	return nil
}

// invertImage inverts the pixel intensities of the given grayscale image.
func invertImage(img *image.Gray) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldGray := img.GrayAt(x, y)
			newGray := color.Gray{Y: 255 - oldGray.Y} // Invert the pixel intensity
			img.SetGray(x, y, newGray)
		}
	}
}
