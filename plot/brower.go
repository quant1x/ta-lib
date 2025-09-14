package plot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"runtime"

	"gitee.com/quant1x/pkg/chart"
)

// OpenURL 用系统默认的浏览器打开一个url
func OpenURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// macOS
		cmd = exec.Command("open", url)
	case "windows":
		// Windows
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		// Linux 或其他操作系统
		cmd = exec.Command("xdg-open", url)
	}
	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

const (
	imgStrPrefix = "data:image/png;base64,"
)

func imgToStr(imgContent bytes.Buffer) string {
	return imgStrPrefix + base64.StdEncoding.EncodeToString(imgContent.Bytes())
}

func ImageContent(c chart.Chart) (string, error) {
	var imgContent bytes.Buffer
	err := c.Render(chart.PNG, &imgContent)
	if err != nil {
		return "", err
	}
	return imgToStr(imgContent), nil
}

// OpenImage TODO: 打开chart有点问题, 原因可能是生成的URL过长
func OpenImage(c chart.Chart) error {
	img, err := ImageContent(c)
	if err != nil {
		return err
	}
	fmt.Println(img)
	return OpenURL(img)
}
