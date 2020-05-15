package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
)

func _main(loadFilePath, saveFilePath *string, chipSize, chipMargin int) int {
	// ファイルよみこみ
	loadfile, err := os.Open(*loadFilePath)
	if err != nil {
		fmt.Println("[error] load error.")
		return 1
	}
	defer loadfile.Close()

	// 画像デコード
	srcImage, _, err := image.Decode(loadfile)
	if err != nil {
		fmt.Println("[error] decode error.")
		return 1
	}

	// リサイズ後の箱を作っとく(チップ間のマージンを考慮)
	newWidth := ((srcImage.Bounds().Dx() / chipSize) * chipMargin) + srcImage.Bounds().Dx()
	newHeight := ((srcImage.Bounds().Dy() / chipSize) * chipMargin) + srcImage.Bounds().Dy()
	destImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 1pxずつ走査しながら画像を加工。

	// y方向
	for destY, y := 0, srcImage.Bounds().Min.Y; y < srcImage.Bounds().Max.Y; y++ {
		// タイルの区切りで間あける(垂直)
		if y%chipSize == 0 {
			destY += chipMargin // ソーシャルディスタンス！！
		}

		// x方向
		for destX, x := 0, srcImage.Bounds().Min.X; x < srcImage.Bounds().Max.X; x++ {
			// タイルの区切りで間あける(水平)
			if x%chipSize == 0 {
				destX += chipMargin // ソーシャルディスタンス！！！！
			}

			// srcからdestに色をコピー
			r, g, b, a := srcImage.At(x, y).RGBA()
			destImage.Set(destX, destY, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
			destX++
		}

		destY++
	}

	// ほぞん
	destFile, err := os.Create(*saveFilePath)
	if err != nil {
		fmt.Println("[error] save error.")
		return 1
	}
	defer destFile.Close()

	if err := png.Encode(destFile, destImage); err != nil {
		fmt.Println("[error] encode error.")
		return 1
	}

	// よかったね。
	fmt.Println("[Success] Yattaze.")

	return 0
}

func main() {
	// 引数いろいろ
	var (
		// 読み込むpng
		loadFilePath = flag.String("i", "", "load file path(*.png)")
		// 書き出すpng
		saveFilePath = flag.String("o", "", "save file path(*.png)")
		// マップチップ1つのサイズ(px)
		chipSize = flag.Int("s", 32, "chip size(pixel)")
		// マップチップ間のマージン(px)
		chipMargin = flag.Int("m", 4, "chip margin(pixel)")
	)
	flag.Parse()

	// 実行
	code := _main(loadFilePath, saveFilePath, *chipSize, *chipMargin)

	// OSに実行結果返却
	os.Exit(code)
}
