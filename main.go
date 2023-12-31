package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
	"image/jpeg"
	"os"
)

var actualScreen int
var defaultQualityScreenshot = 80

var myApp = app.New()
var myWindow = myApp.NewWindow("My Screen")

func main() {
	//myApp := app.New()
	//myWindow := myApp.NewWindow("My Screen")
	myWindow.Resize(fyne.NewSize(520, 320))
	myWindow.SetFixedSize(true)
	myWindow.SetContent(container.NewVBox(selectWindowContainer(), selectQualityContainer(), selectOutputFolder(), captureWindowContainer()))
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
	responseContainer := container.NewVBox(widget.NewLabel(""))
	return container.NewVBox(
		widget.NewLabel("Output file name"),
		output,
		widget.NewButton("Capture", func() {
			msg, err := captureScreenshot(actualScreen, output.Text)
			if err != nil {
				responseContainer.Objects[0] = widget.NewLabel(err.Error())
				responseContainer.Refresh()
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
			defaultQualityScreenshot = 10
		case "Medium":
			defaultQualityScreenshot = 50
		case "High":
			defaultQualityScreenshot = 80
		}
	})
	quality.Selected = "High"
	return container.NewVBox(
		widget.NewLabel("Select a quality"),
		quality,
	)
}

var outputFolder fyne.ListableURI
var outputFolderLabel = widget.NewLabel("Output folder: No selected")

func selectOutputFolder() *fyne.Container {
	openDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if uri == nil {
			dialog.ShowError(fmt.Errorf("No folder selected"), myWindow)
			return
		}
		outputFolder = uri
		outputFolderLabel.SetText("Output folder: " + uri.String())
	}, myWindow)

	return container.NewVBox(
		widget.NewLabel("Select a folder"),
		outputFolderLabel,
		widget.NewButton("Open Folder Dialog", func() {
			openDialog.Show()
		}),
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

	outputPath := outputFolder.String() + "/" + fileName
	fmt.Println(outputPath)

	bounds := screenshot.GetDisplayBounds(screen)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			panic(err.Error())
		}
	}(file)

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: defaultQualityScreenshot})
	if err != nil {
		return "", err
	}

	return "Screenshot saved to " + fileName, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}
