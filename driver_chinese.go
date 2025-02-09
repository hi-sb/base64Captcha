package base64Captcha

import (
	"github.com/golang/freetype/truetype"
	"image/color"
	"math/rand"
	"strings"
)

//DriverChinese is a driver of unicode Chinese characters.
type DriverChinese struct {
	//Height png height in pixel.
	Height int
	//Width Captcha png width in pixel.
	Width int

	//NoiseCount text noise count.
	NoiseCount int

	//ShowLineOptions := OptionShowHollowLine | OptionShowSlimeLine | OptionShowSineLine .
	ShowLineOptions int

	//Length random string length.
	Length int

	//Source is a unicode which is the rand string from.
	Source string

	//BgColor captcha image background color (optional)
	BgColor *color.RGBA
	//Fonts loads by name see fonts.go's comment
	Fonts []string

	fontsArray []*truetype.Font
}

//NewDriverChinese creates a driver of Chinese characters
func NewDriverChinese(height int, width int, noiseCount int, showLineOptions int, length int, source string, bgColor *color.RGBA, fonts []string) *DriverChinese {
	tfs := []*truetype.Font{}
	for _, fff := range fonts {
		tf := loadFontByName("fonts/" + fff)
		tfs = append(tfs, tf)
	}
	if len(tfs) == 0 {
		tfs = fontsAll
	}
	return &DriverChinese{Height: height, Width: width, NoiseCount: noiseCount, ShowLineOptions: showLineOptions, Length: length, Source: source, BgColor: bgColor, fontsArray: tfs}
}

//ConvertFonts loads fonts by names
func (d *DriverChinese) ConvertFonts() *DriverChinese {
	tfs := []*truetype.Font{}
	for _, fff := range d.Fonts {
		tf := loadFontByName("fonts/" + fff)
		tfs = append(tfs, tf)
	}
	if len(tfs) == 0 {
		tfs = fontsAll
	}
	d.fontsArray = tfs
	return d
}

//GenerateQuestionAnswer generates captcha content and its answer
func (d *DriverChinese) GenerateQuestionAnswer() (content, answer string) {

	ss := strings.Split(d.Source, ",")
	length := len(ss)
	if length == 1 {
		c := randText(d.Length, ss[0])
		return c, c
	}
	if length <= d.Length {
		c := randText(d.Length, TxtNumbers+TxtAlphabet)
		return c, c
	}

	res := make([]string, d.Length)
	for k := range res {
		res[k] = ss[rand.Intn(length)]
	}

	content = strings.Join(res, "")
	return content, content
}

//GenerateItem generates captcha item(image)
func (d *DriverChinese) GenerateItem(content string) (item Item, err error) {

	var bgc color.RGBA
	if d.BgColor != nil {
		bgc = *d.BgColor
	} else {
		bgc = randLightColor()
	}
	itemChar := NewItemChar(d.Width, d.Height, bgc)

	//draw hollow line
	if d.ShowLineOptions&OptionShowHollowLine == OptionShowHollowLine {
		itemChar.drawHollowLine()
	}

	//draw slime line
	if d.ShowLineOptions&OptionShowSlimeLine == OptionShowSlimeLine {
		itemChar.drawSlimLine(3)
	}

	//draw sine line
	if d.ShowLineOptions&OptionShowSineLine == OptionShowSineLine {
		itemChar.drawSineLine()
	}

	//draw noise
	if d.NoiseCount > 0 {
		source := TxtNumbers + TxtAlphabet + ",.[]<>"
		noise := randText(d.NoiseCount, strings.Repeat(source, d.NoiseCount))
		err = itemChar.drawNoise(noise, d.fontsArray)
		if err != nil {
			return
		}
	}

	//draw content
	err = itemChar.drawText(content, d.fontsArray)
	if err != nil {
		return
	}

	return itemChar, nil
}
