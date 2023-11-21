package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
	"image/jpeg"
	"os"
)

var actualScreen int
var qualityScreenshot int

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("My Screen")
	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.SetFixedSize(true)
	myWindow.SetContent(container.NewVBox(selectWindowContainer(), selectQualityContainer(), captureWindowContainer()))
	myWindow.ShowAndRun()
}

func selectWindowContainer() *fyne.Container {
	screensStr := getAvaliableScreens()
	windowSelect := widget.NewSelect(screensStr, func(value string) {
		for i, screen := range screensStr {
			if screen == value {
				actualScreen = i
				break
			}
		}
	})
	windowSelect.Selected = screensStr[0]
	return container.NewVBox(
		widget.NewLabel("Select a screen"),
		windowSelect,
	)
}

func captureWindowContainer() *fyne.Container {
	output := widget.NewEntry()
	output.SetPlaceHolder("Output file name (default: screenshot.jpg)")
	output.OnChanged = func(value string) {
		output.SetText(value)
	}
	responseContainer := container.NewVBox(widget.NewLabel(""))
	return container.NewVBox(
		widget.NewLabel("Output file name"),
		output,
		widget.NewButton("Capture", func() {
			msg, err := captureScreenshot(actualScreen, output.Text)
			if err != nil {
				widget.NewLabel(err.Error())
			}
			responseContainer.Objects[0] = widget.NewLabel(msg)
			responseContainer.Refresh()
		}), responseContainer,
	)
}

func selectQualityContainer() *fyne.Container {
	quality := widget.NewSelect([]string{"Low", "Medium", "High"}, func(value string) {
		switch value {
		case "Low":
			qualityScreenshot = 10
		case "Medium":
			qualityScreenshot = 50
		case "High":
			qualityScreenshot = 100
		}
	})
	quality.Selected = "High"
	return container.NewVBox(
		widget.NewLabel("Select a quality"),
		quality,
	)
}

func getAvaliableScreens() []string {
	n := screenshot.NumActiveDisplays()
	screensStr := make([]string, n)
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		infoScreen := fmt.Sprintf("Id: '%d', Bounds '%v'\n", i, bounds)
		screensStr[i] = infoScreen
	}
	return screensStr
}

func captureScreenshot(screen int, output string) (string, error) {
	if output == "" {
		outputFileNameDefault := fmt.Sprintf("screenshot_(%d)", screen)
		output = outputFileNameDefault
	}

	fileName := fmt.Sprintf("%s.jpg", output)
	counter := 1
	for fileExists(fileName) {
		fileName = fmt.Sprintf("%s(%d).jpg", output, counter)
		counter++
	}

	bounds := screenshot.GetDisplayBounds(screen)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err.Error())
		}
	}(file)

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: qualityScreenshot})
	if err != nil {
		return "", err
	}

	return "Screenshot saved to " + fileName, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}
