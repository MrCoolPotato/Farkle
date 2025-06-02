package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/theme"
    "fyne.io/fyne/v2/widget"
    "image/color"
)

type whiteTheme struct{}

func (whiteTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
    switch name {
    case theme.ColorNameInputBackground:
        return color.White
    case theme.ColorNameInputBorder,
        theme.ColorNameFocus,
        theme.ColorNameShadow,
        theme.ColorNameSeparator,
        theme.ColorNameHover,
        theme.ColorNameDisabled:
        return color.Transparent
    }
    return theme.LightTheme().Color(name, variant)
}
func (whiteTheme) Font(style fyne.TextStyle) fyne.Resource { return theme.LightTheme().Font(style) }
func (whiteTheme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.LightTheme().Icon(name) }
func (whiteTheme) Size(name fyne.ThemeSizeName) float32 {
    return theme.LightTheme().Size(name)
}

func whiteBg(obj fyne.CanvasObject) fyne.CanvasObject {
    bg := canvas.NewRectangle(color.White)
    return container.NewMax(bg, obj)
}

func main() {
    a := app.NewWithID("com.ap.farkle")
    a.Settings().SetTheme(whiteTheme{})
    w := a.NewWindow("Farkle")

    output := widget.NewMultiLineEntry()
    output.SetPlaceHolder("Welcome to Farkle GUI...\n")
    output.Wrapping = fyne.TextWrapWord
    output.Disable()

    input := widget.NewEntry()
    input.SetPlaceHolder("Type commands here…")
    input.OnSubmitted = func(s string) {
        if s != "" {
            output.SetText(output.Text + s + "\n")
            input.SetText("")
        }
    }

    leftTop := whiteBg(output)
    leftBottom := whiteBg(input)
    left := container.NewBorder(nil, leftBottom, nil, nil, leftTop)

    animBG := canvas.NewRectangle(color.White)
    animBG.SetMinSize(fyne.NewSize(200, 200))
    animText := canvas.NewText("🎲", theme.ForegroundColor())
    animText.TextSize = 64
    anim := container.NewCenter(animBG, animText)

    grid := container.NewGridWithColumns(2, left, anim)

    pad := container.New(layout.NewPaddedLayout(), grid)
    w.SetContent(pad)
    w.Resize(fyne.NewSize(800, 600))
    w.ShowAndRun()
}