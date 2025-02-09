// Copyright 2017 Eric Zhou. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package base64Captcha

import (
	"image"
)

//DriverDigit config for captcha-engine-digit.
type DriverDigit struct {
	// Height png height in pixel.
	Height int
	// Width Captcha png width in pixel.
	Width int
	// DefaultLen Default number of digits in captcha solution.
	Length int
	// MaxSkew max absolute skew factor of a single digit.
	MaxSkew float64
	// DotCount Number of background circles.
	DotCount int
}

//NewDriverDigit creates a driver of digit
func NewDriverDigit(height int, width int, length int, maxSkew float64, dotCount int) *DriverDigit {
	return &DriverDigit{Height: height, Width: width, Length: length, MaxSkew: maxSkew, DotCount: dotCount}
}

//DefaultDriverDigit is a default driver of digit
var DefaultDriverDigit = NewDriverDigit(80, 240, 5, 0.7, 80)

//GenerateQuestionAnswer creates captcha content and answer
func (d *DriverDigit) GenerateQuestionAnswer() (q, a string) {
	digits := randomDigits(d.Length)
	a = parseDigitsToString(digits)
	return a, a
}

//GenerateItem creates digit captcha item
func (d *DriverDigit) GenerateItem(content string) (item Item, err error) {
	// Initialize PRNG.
	itemDigit := NewItemDigit(d.Width, d.Height, d.DotCount, d.MaxSkew)
	//parse digits to string
	digits := stringToFakeByte(content)
	itemDigit.rng.Seed(deriveSeed(imageSeedPurpose, randomId(), digits))

	itemDigit.Paletted = image.NewPaletted(image.Rect(0, 0, d.Width, d.Height), itemDigit.getRandomPalette())
	itemDigit.calculateSizes(d.Width, d.Height, len(digits))
	// Randomly position captcha inside the image.
	maxx := d.Width - (itemDigit.width+itemDigit.dotSize)*len(digits) - itemDigit.dotSize
	maxy := d.Height - itemDigit.height - itemDigit.dotSize*2
	var border int
	if d.Width > d.Height {
		border = d.Height / 5
	} else {
		border = d.Width / 5
	}
	x := itemDigit.rng.Int(border, maxx-border)
	y := itemDigit.rng.Int(border, maxy-border)
	// Draw digits.
	for _, n := range digits {
		itemDigit.drawDigit(digitFontData[n], x, y)
		x += itemDigit.width + itemDigit.dotSize
	}
	// Draw strike-through line.
	itemDigit.strikeThrough()
	// Apply wave distortion.
	itemDigit.distort(itemDigit.rng.Float(5, 10), itemDigit.rng.Float(100, 200))
	// Fill image with random circles.
	itemDigit.fillWithCircles(d.DotCount, itemDigit.dotSize)
	return itemDigit, nil
}
