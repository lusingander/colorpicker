![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/lusingander/colorpicker)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/lusingander/colorpicker?tab=doc)](https://pkg.go.dev/github.com/lusingander/colorpicker?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/lusingander/colorpicker)](https://goreportcard.com/report/github.com/lusingander/colorpicker)
![GitHub](https://img.shields.io/github/license/lusingander/colorpicker)

# colorpicker

Color picker component for [Fyne](https://fyne.io/)

<img src="./resource/popup.gif" width=300>

## Usage

```go
picker := colorpicker.New(200 /* height */, colorpicker.StyleHue /* Style */)
picker.SetOnChanged(func(c color.Color) {
    // called when the color is changed on the picker
    fmt.Println(c)
})

// you can use it just like any other Fyne widget
fyne.NewContainer(picker)
```

## Documentation

See [pkg.go.dev](https://pkg.go.dev/github.com/lusingander/colorpicker?tab=doc)

## Example

### colorpicker

You can see all the styles implemented.

[colorpicker/cmd/colorpicker/](./cmd/colorpicker/)

<img src="./resource/image.png" width=500>

----

### colorpicker-popup

Example of embedding in Fyne's custom dialog.

[colorpicker/cmd/colorpicker-popup/](./cmd/colorpicker-popup/)

<img src="./resource/image2.png" width=400>
