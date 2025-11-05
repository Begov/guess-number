package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Result struct {
	Date         time.Time `json:"date"`
	Win          bool      `json:"win"`
	AttemptsUsed int       `json:"attemptsUsed"`
}

func NewResult(d time.Time, win bool, attemptsUsed int) *Result {
	return &Result{
		Date:         d,
		Win:          win,
		AttemptsUsed: attemptsUsed,
	}
}

var filename string
var results []Result

func init() {
	filename, _ = os.Getwd()
	filename += "/data/results.json"
	data, err := os.ReadFile(filename)
	if err == nil {
		json.Unmarshal(data, &results)
	}
}

var green = color.New(color.FgGreen).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

type Game struct {
	minNum       int
	maxNum       int
	secretNumber int
	attemptsLeft int
	attempts     []int
}

func NewGame() *Game {
	return &Game{
		minNum:       0,
		maxNum:       0,
		secretNumber: 0,
		attemptsLeft: 0,
		attempts:     make([]int, 0),
	}
}

func (g *Game) setDifficulty(num int) {
	g.minNum = 1
	switch num {
	case 1:
		g.maxNum = 50
		g.secretNumber = RandInt(1, 50)
		g.attemptsLeft = 15
	case 2:
		g.maxNum = 100
		g.secretNumber = RandInt(1, 100)
		g.attemptsLeft = 10
	case 3:
		g.maxNum = 200
		g.secretNumber = RandInt(1, 200)
		g.attemptsLeft = 5
	}
	fmt.Printf("Игра %s - от %s до %s началась!\n", green("\"Угадай число\""), yellow(g.minNum), yellow(g.maxNum))
	fmt.Printf("Угадайте число за %s попыток!\n", yellow(g.attemptsLeft))
}

var game Game

func StartGame() {
	play := true
	for play {
		difficulty, err := ChooseDifficulty()
		if err != nil {
			fmt.Println(red(err))
			continue
		}

		game = *NewGame()
		game.setDifficulty(difficulty)

		d, win, attempts := CheckGuess()
		results = append(results, *NewResult(d, win, attempts))
		saveToFile()
		play = AskPlayAgain()
	}
}

func ChooseDifficulty() (int, error) {
	fmt.Println(green("Easy: 1–50, 15 попыток;"), yellow("Medium: 1–100, 10 попыток;"), red("Hard: 1–200, 5 попыток"))
	fmt.Printf("Выберите сложность игры: %s; %s; %s 👆: ", green("1 - Easy"), yellow("2 - Medium"), red("3 - Hard"))
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	difficulty, err := strconv.Atoi(input)

	if err != nil || difficulty < 1 || difficulty > 3 {
		return 0, fmt.Errorf("Ошибка: уровень указан неверно")
	}

	return difficulty, nil
}

func CheckGuess() (time.Time, bool, int) {
	var win bool

	reader := bufio.NewReader(os.Stdin)
	for i := 1; true; i++ {

		if len(game.attempts) > 0 {
			attempts := make([]string, len(game.attempts))
			for i, v := range game.attempts {
				attempts[i] = strconv.Itoa(v)
			}
			fmt.Println("Вы уже вводили числа:", strings.Join(attempts, ", "))
		}

		fmt.Printf("Попытка #%d - Введите число: ", i)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		usernum, err := strconv.Atoi(input)

		if err != nil {
			fmt.Println(red("Введённое значение не является числом"))
			continue
		}

		game.attempts = append(game.attempts, usernum)
		difference := math.Abs(float64(usernum - game.secretNumber))

		if usernum == game.secretNumber {
			fmt.Println(green("Ураа! Победа!"))
			win = true
			break
		}

		if difference <= 5 {
			fmt.Println("🔥 Горячо")
		} else if difference <= 15 {
			fmt.Println("🙂 Тепло")
		} else {
			fmt.Println("❄️  Холодно")
		}

		if usernum > game.secretNumber {
			fmt.Println("Секретное число меньше 👇")
		} else {
			fmt.Println("Секретное число больше 👆")
		}

		game.attemptsLeft--
		fmt.Println("Осталось попыток:", game.attemptsLeft)

		if game.attemptsLeft == 0 {
			fmt.Println(red("Попытки закончились — вы проиграли"))
			win = false
			break
		}
	}

	return time.Now(), win, len(game.attempts)
}

func AskPlayAgain() bool {
	fmt.Printf("Хотите сыграть ещё раз? %s; %s: ", green("1 - Да"), red("0 - Нет"))
	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		answer, err := strconv.Atoi(input)

		if err != nil || answer < 0 || answer > 1 {
			fmt.Print(yellow("Укажите 1 или 0: "))
			continue
		}

		if answer == 0 {
			return false
		} else {
			return true
		}
	}
}

func saveToFile() {

	dataJson, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println("Не удалось получить json")
		return
	}

	if err := os.WriteFile(filename, dataJson, 0644); err != nil {
		fmt.Println("Не удалось сохранить результат")
		return
	}
}

func RandInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}
