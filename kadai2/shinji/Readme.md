# <課題2-1> 【TRY】 io.Readerとio.Writer
## 概要
### io.Readerとio.Writerについて調べてみよう

- 標準パッケージでどのように使われているか
- io.Readerとio.Writerがあることで、どういう利点があるのか具体例を挙げて考えてみる

## 答案

- 標準パッケージでどのように使われているか
  - bufio.Readerやos.Fileに実装されている  
  https://golang.org/pkg/bufio/#Reader.Read  
  https://golang.org/pkg/os/#File.Read  
  ※Reader/Writerインターフェイスに規定されているメソッドはRead/Write関数のみなので、<br>
  Read/Write関数を実装するだけで、Reader/Writerインターフェイスを実装する事が出来る
  - 上記より、io.Readerまたはio.Writerを引数に指定している関数にbufio.Readerやos.Fileを渡すことができる

- io.Readerとio.Writerがあることで、どういう利点があるのか具体例を挙げて考えてみる
  - io.Reader/io.Writer型の変数に対して同じ処理を行う関数を一般化出来る
  
（io.Readerとio.Writerが無かった場合）
```golang
func sayb(r *Buffer) {
	fmt.Print(r)
}

func sayf(r *File) {
  fmt.Print(r)
}

func main() {
	r := bytes.NewBufferString("hello")
	// r, _ := os.Open("hoge.txt") // ファイルでもOK
	say(r)
}
```
（io.Readerとio.Writerがある場合）
```golang
func say(r io.Reader) {
	fmt.Print(r)
}

func main() {
	r := bytes.NewBufferString("hello")
	// r, _ := os.Open("hoge.txt") // ファイルでもOK
	say(r)
}
```

# <課題2-2> 【TRY】 テストを書いてみよう
## 概要
### 1回目の課題のテストを作ってみて下さい

- テストのしやすさを考えてリファクタリングしてみる
- テストのカバレッジを取ってみる
- テーブル駆動テストを行う
- テストヘルパーを作ってみる

## 答案

|課題|答案|
|:---:|:---:|
|テストのしやすさを考えてリファクタリングしてみる|関数を切り出し、単体テストしやすいようにした|
|テストのカバレッジを取ってみる|-coverオプションを使用してカバレッジを取得した <br> （coverage: 81.8% of statements）|
|テーブル駆動テストを行う|テストは全てテーブル駆動テストで実装|
|テストヘルパーを作ってみる|テストヘルパーとして、 <br> decodeForTest(), errCheck(), existCheck(）を実装した|

## 実行手順
```shell
$ cd {Path_To_Repository}/kadai2/shinji

# テスト実行（カバレッジを表示）
$ go test ./... -cover 

# テスト実行（HTMLにカバレッジを保存）
$ go test ./... -coverprofile=cover.out 
$ go tool cover -html=cover.out -o cover.html
```

## ディレクトリ構造
```
shinji
├── Readme.md
├── convimg
│   ├── convimg.go
│   ├── convimg_test.go
│   └── export_test.go
├── cover.html
├── cover.out
├── go.mod
├── main
├── main.go
├── testdata
└── testdata.sh

```

## テストデータのディレクトリ構造
```
testdata
├── azarashi.jpg
├── tanuki.jpg
├── osaru.png
└── img
    ├── azarashi.jpg
    └── tanuki.jpg
```
