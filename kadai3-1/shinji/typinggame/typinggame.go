package typinggame

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var Answers []string

func init() {
	setAns()
	rand.Seed(time.Now().UnixNano())
}

func setAns() {
	Answers = append(Answers, "apple", "grape", "pineapple", "peach", "kiwi", "banana", "melon")
}

func GetInput(r io.Reader) <-chan string {
	input := make(chan string)

	go func() {
		s := bufio.NewScanner(r)
		for s.Scan() {
			input <- s.Text()
		}
		close(input)
	}()

	return input
}

func Ask() string {
	// 質問をランダムに生成
	n := rand.Intn(len(Answers))
	answer := Answers[n]

	// 質問を表示
	fmt.Println(answer)
	fmt.Print(">")

	return answer
}

func Check(input string, answer string, score int, total int) (int, int) {
	total++
	if input == answer {
		score++
		fmt.Printf("correct! score: %v/%v \n\n", score, total)
	} else {
		fmt.Printf("incorrect! score: %v/%v \n\n", score, total)
	}
	return score, total
}

func Do() {
	// チャネルを用意
	in := GetInput(os.Stdin)
	timelimit := time.After(15 * time.Second)

	// スコアを初期化
	var score, total int

	for {
		// 出題
		answer := Ask()

		select {
		// 入力があれば判定
		case input := <-in:
			score, total = Check(input, answer, score, total)

		// 時間切れならスコアを表示して終了
		case <-timelimit:
			fmt.Printf("\n\n Time out! score: %v/%v \n\n", score, total)
			os.Exit(0)
		}
	}
}
