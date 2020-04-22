package modules

import (
	"bytes"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/golang/freetype"
	"image"
	"image/png"
	"log"
	"strconv"
)

func createImage(textName string) *bytes.Buffer {
	resultPNG := bytes.NewBufferString("")
	//创建位图,坐标x,y,长宽x,y
	img := image.NewNRGBA(image.Rect(0, 0, 120, 33))
	// 字体数据，为了方便打包进执行文件，这里直接暴力写入代码中。
	fontBytes := constant.CONSOLE_FONT
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return nil
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(40)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)
	//设置字体显示位置
	pt := freetype.Pt(5, 20+int(c.PointToFixed(40)>>8))
	_, err = c.DrawString(textName, pt)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = png.Encode(resultPNG, img)
	return resultPNG
}

// 启动时的标题
func DrawBanner(text string) string {
	pngBuf := createImage(text)
	if pngBuf == nil {
		return ""
	}
	//灰度替换字符
	base := "$TIS"
	image1, _ := png.Decode(pngBuf)
	bounds := image1.Bounds() //获取图像的边界信息
	logo := ""                //存储最终的字符画string
	for y := 0; y < bounds.Dy(); y += 2 {
		for x := 0; x < bounds.Dx(); x++ {
			pixel := image1.At(x, y)   //获取像素点
			r, g, b, _ := pixel.RGBA() //获取像素点的rgb
			r = r & 0xFF
			g = g & 0xFF
			b = b & 0xFF
			//灰度计算
			gray := 0.299*float64(r) + 0.578*float64(g) + 0.114*float64(b)
			temp := fmt.Sprintf("%.0f", gray*float64(len(base)+1)/255)
			index, _ := strconv.Atoi(temp)
			//根据灰度索引字符并保存
			if index >= len(base) {
				logo += " "
			} else {
				logo += string(base[index])
			}

		}
		logo += "\r\n"
	}
	//输出字符画
	return fmt.Sprintf("\n\033[31;1m%s", logo)
}
