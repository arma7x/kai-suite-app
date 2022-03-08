package custom_theme

import(
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	
)

type LightMode struct{}

func (LightMode) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameBackground:
		return color.NRGBA{0xff, 0xff, 0xff, 0xff}
	case theme.ColorNameButton:
		return color.Transparent
	case theme.ColorNameDisabledButton:
		return color.NRGBA{0xe5, 0xe5, 0xe5, 0xff}
	case theme.ColorNameDisabled:
		return color.NRGBA{0x0, 0x0, 0x0, 0x42}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xf4, G: 0x43, B: 0x36, A: 0xff}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0x7f}
	case theme.ColorNameForeground:
		return color.NRGBA{0x21, 0x21, 0x21, 0xff}
	case theme.ColorNameHover:
		return color.NRGBA{0x0, 0x0, 0x0, 0x0f}
	case theme.ColorNameInputBackground:
		return color.NRGBA{0x0, 0x0, 0x0, 0x19}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{0x88, 0x88, 0x88, 0xff}
	case theme.ColorNamePressed:
		return color.NRGBA{0x0, 0x0, 0x0, 0x19}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x21, G: 0x96, B: 0xf3, A: 0xff}
	case theme.ColorNameScrollBar:
		return color.NRGBA{0x0, 0x0, 0x0, 0x99}
	case theme.ColorNameShadow:
		return color.NRGBA{0x0, 0x0, 0x0, 0x33}
	default:
		return theme.DefaultTheme().Color(c, v)
	}
}

func (LightMode) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return theme.DefaultTheme().Font(s)
	}
	if s.Bold {
		if s.Italic {
			return theme.DefaultTheme().Font(s)
		}
		return theme.DefaultTheme().Font(s)
	}
	if s.Italic {
		return theme.DefaultTheme().Font(s)
	}
	return theme.DefaultTheme().Font(s)
}

func (LightMode) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (LightMode) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 3
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInputBorder:
		return 2
	case theme.SizeNameHeadingText:
		return 15
	default:
		return theme.DefaultTheme().Size(s)
	}
}
