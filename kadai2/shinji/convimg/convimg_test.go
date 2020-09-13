package convimg_test

import (
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/shinnosuke/gopherdojo-studyroom/kadai2/shinji/convimg"
)

var patherr *os.PathError

func TestMain(m *testing.M) {
	// reset testdata
	out, err := exec.Command("../testdata.sh").Output()
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Print(string(out))

	// run test
	code := m.Run()

	os.Exit(code)
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name    string
		srcPath string
		isErr   bool
		want    error
	}{
		{name: "decode jpg", srcPath: "../testdata/img/azarashi.jpg", isErr: false, want: nil},
		{name: "decode png", srcPath: "../testdata/osaru.png", isErr: false, want: nil},
		{name: "no such file or dir", srcPath: "../testdata/img/dontexist.jpg", isErr: true, want: patherr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := convimg.Decode(test.srcPath)

			// エラーチェック
			errCheck(t, err, test.isErr, test.want)
		})
	}
}

func TestConvExt(t *testing.T) {
	tests := []struct {
		name     string
		srcPath  string
		to       convimg.Ext
		expected string
	}{
		{name: "jpg_to_png", srcPath: "testdata/img/azarashi.jpg", to: ".png", expected: "testdata/img/azarashi.png"},
		{name: "jpg_to_gif", srcPath: "testdata/img/azarashi.jpg", to: ".gif", expected: "testdata/img/azarashi.gif"},
		{name: "png_to_jpg", srcPath: "testdata/img/azarashi.png", to: ".jpg", expected: "testdata/img/azarashi.jpg"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := convimg.ConvExt(test.srcPath, test.to)
			if actual != test.expected {
				t.Errorf("Test Failed: got %v, want %v", actual, test.expected)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	//テスト用のimageを生成
	img, _ := decodeForTest(t, "../testdata/img/azarashi.jpg")
	//var empty-img image.Image

	tests := []struct {
		name    string
		dstPath string
		img     image.Image
		to      convimg.Ext
		isErr   bool
		want    error
	}{
		{name: "encode gif", dstPath: "../testdata/img/azarashi.gif", img: img, to: ".gif", isErr: false, want: nil},
		{name: "encode png", dstPath: "../testdata/img/azarashi.png", img: img, to: ".png", isErr: false, want: nil},
		//{name: "empty img", dstPath: "../testdata/img/azarashi.gif", img: empty-img, to: ".gif", isErr: true, want: ""}, //対策を未実装
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := convimg.Encode(test.dstPath, test.img, test.to)

			// エラーチェック
			errCheck(t, err, test.isErr, test.want)

			// ファイルの有無をテスト
			if !test.isErr {
				if _, existerr := os.Stat(test.dstPath); os.IsNotExist(existerr) {
					t.Errorf("File %v should exist, but does not exist", test.dstPath)
				}
			}
		})
	}
}

func TestDo(t *testing.T) {
	tests := []struct {
		name    string
		srcPath string
		to      convimg.Ext
		dstPath string
		rmSrc   bool
		isErr   bool
		want    error
	}{
		{name: "jpg to png dont remove", srcPath: "../testdata/img/azarashi.jpg", to: ".png", dstPath: "../testdata/img/azarashi.png", rmSrc: false, isErr: false, want: nil},
		{name: "png to jpg dont remove", srcPath: "../testdata/osaru.png", to: ".jpg", dstPath: "../testdata/osaru.jpg", rmSrc: false, isErr: false, want: nil},
		{name: "jpg to png remove", srcPath: "../testdata/img/azarashi.jpg", to: ".png", dstPath: "../testdata/img/azarashi.png", rmSrc: true, isErr: false, want: nil},
		{name: "no such file or dir", srcPath: "../testdata/img/dontexist.jpg", to: ".png", dstPath: "../testdata/img/azarashi.png", rmSrc: true, isErr: true, want: patherr},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := convimg.Do(test.srcPath, test.to, test.rmSrc)

			// エラーチェック
			errCheck(t, err, test.isErr, test.want)

			if !test.isErr {
				// ファイルの有無をテスト
				existCheck(t, test.dstPath, true)

				// ファイルの削除についてテスト
				existCheck(t, test.srcPath, !test.rmSrc) //rmSrcがtrue → shouldExistはfalse
			}

		})
	}
}

func decodeForTest(t *testing.T, srcPath string) (image.Image, error) {
	t.Helper()

	// ファイルオープン
	src, openerr := os.Open(filepath.Clean(srcPath))
	if openerr != nil {
		return nil, openerr
	}
	var closeerr error
	defer func() {
		if closeerr = src.Close(); closeerr != nil {
			fmt.Fprintln(os.Stderr, "ERROR:", closeerr)
		}
	}()

	// ファイルオブジェクトを画像オブジェクトに変換
	img, _, decodeerr := image.Decode(src)
	if decodeerr != nil {
		return nil, decodeerr
	}

	return img, nil
}

func errCheck(t *testing.T, err error, isErr bool, want error) {
	t.Helper()

	switch {
	case !isErr && err != nil: //正常系なのにエラー
		t.Errorf("want no err, but got [%v]", reflect.TypeOf(err))
	case isErr && reflect.TypeOf(err) != reflect.TypeOf(want): //異常系なのにエラー無し、もしくは想定外のエラー
		t.Errorf("want [%v], but got [%v]", reflect.TypeOf(want), reflect.TypeOf(err))
	}
}

func existCheck(t *testing.T, path string, shouldExist bool) {
	t.Helper()

	_, existerr := os.Stat(path)
	switch {
	case shouldExist && os.IsNotExist(existerr): // あるはずなのにない
		fmt.Printf("File %v should exist, but does not exist", path)
	case !shouldExist && !os.IsNotExist(existerr): // ないはずなのにある
		fmt.Printf("File %v should be removed, but does exist", path)
	}
}
