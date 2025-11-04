package game

import (
	"bufio"
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var green = color.New(color.FgGreen).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

type Game struct {
	secretNumber int
	attemptsLeft int
	attempts     []int
}

func NewGame(secretNumber, attemptsLeft int) *Game {
	return &Game{
		secretNumber: secretNumber,
		attemptsLeft: attemptsLeft,
		attempts:     make([]int, 0),
	}
}

func (g *Game) setDifficulty(num int) {
	switch num {
	case 1:
		g.secretNumber = RandInt(1, 50)
		g.attemptsLeft = 15
	case 2:
		g.secretNumber = RandInt(1, 100)
		g.attemptsLeft = 10
	case 3:
		g.secretNumber = RandInt(1, 200)
		g.attemptsLeft = 5
	}
}

var game Game

func StartGame() {
	play := true

	for play {
		difficulty, err := ChooseDifficulty()
		if err != nil {
			fmt.Println(err)
			continue
		}

		game = *NewGame(0, 0)
		game.setDifficulty(difficulty)

		t, status, attempts := CheckGuess()
		fmt.Println(t, status, attempts)
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
	for {

		if len(game.attempts) > 0 {
			attempts := make([]string, len(game.attempts))
			for i, v := range game.attempts {
				attempts[i] = strconv.Itoa(v)
			}
			fmt.Println("Вы уже вводили числа:", strings.Join(attempts, ", "))
		}

		fmt.Print("Введите число: ")
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
		} else if difference <= 5 {
			fmt.Println("🔥 Горячо")
		} else if difference <= 15 {
			fmt.Println("🙂 Тепло")
		} else {
			fmt.Println("❄️  Холодно")
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
			fmt.Print("Укажите 1 или 0: ")
			continue
		}

		if answer == 0 {
			return false
		} else {
			return true
		}
	}
}

func RandInt(min, max int) int {
	return rand.IntN(max-min+1) + min
}
