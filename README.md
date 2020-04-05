colorpicker
====

Color picker component for [Fyne](https://fyne.io/)

## Usage

```go
picker := colorpicker.New(200 /* height */, colorpicker.StyleCircle /* Style */)
picker.SetOnChanged(func(c color.Color) {
    // called when the color is changed on the picker
    fmt.Println(c)
})

// you can use it just like any other Fyne widget
fyne.NewContainer(picker)
```

## Example

[colorpicker/cmd/colorpicker/](./cmd/colorpicker/)

<img src="./resource/image.png" width=500>
