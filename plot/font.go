package plot

import (
	"embed"
	"io"
	"sync"

	"gitee.com/quant1x/gox/api"
	"gitee.com/quant1x/pkg/freetype/truetype"
)

const (
	fontSimHei = "SimHei.ttf"
)

const (
	// ResourcesPath 资源路径
	ResourcesPath = "fonts"
)

//go:embed fonts/*
var fonts embed.FS

var (
	_defaultFontLock sync.Mutex
	_defaultFont     *truetype.Font
)

// GetDefaultFont returns the default font (Roboto-Medium).
func GetDefaultFont() (*truetype.Font, error) {
	if _defaultFont == nil {
		_defaultFontLock.Lock()
		defer _defaultFontLock.Unlock()
		if _defaultFont == nil {
			f, err := api.OpenEmbed(fonts, ResourcesPath+"/"+fontSimHei)
			if err != nil {
				return nil, err
			}
			data, err := io.ReadAll(f)
			if err != nil {
				return nil, err
			}
			font, err := truetype.Parse(data)
			if err != nil {
				return nil, err
			}
			_defaultFont = font
		}
	}
	return _defaultFont, nil
}

// DefaultFont 默认字体
func DefaultFont() *truetype.Font {
	font, _ := GetDefaultFont()
	return font
}
