package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type Color int

const (
	rightBtn = "￪"
	leftBtn  = "￫"
	topBtn   = "￭"
	downBtn  = "￬"
)

var screenLines = [10][10]string{}
var snakePositions = [][]int{}
var xP, yP int = 10, 10
var keyPress string
var keyPressChar string
var food = []int{5, 5}
var snakeEat bool = false
var point int = 0
var startGame bool = false

func main() {
	go TrackClickGo()
	start()
}

func startMenu() {
	var input string

	fmt.Println("Если у тебя аллергия на мерцание, то не играй")
	fmt.Println("Нынешних знаний не хватает, что бы убрать это мерцание")
	fmt.Println("Начнем играть?")

	fmt.Println("1 - да")
	fmt.Println("2 - нет")

	fmt.Scanf("%s\n", &input)

	if input == "2" {
		fmt.Println("У меня есть выбор?")
	}

	if input == "1" {
		fmt.Println("Ну погнали")
	}

	startGame = true
	time.Sleep(1000 * time.Millisecond)

	resetGameData()
}

func resetGameData() {
	keyPress = downBtn
	snakePositions = [][]int{{0, 0}}
	food = []int{2, 2}
	snakeEat = false
	point = 0
}

func start() {
	startMenu()
	for {
		loops()
		if !startGame {
			break
		}
	}
	loseInfoPrint()
	start()
}

func loseInfoPrint() {
	fmt.Println("Упс! ты проиграл....")
	time.Sleep(1000 * time.Millisecond)
}

func loops() {
	witdrawGameScreenLogics()
	endLoopPrintInfo()
	oneStepForward()
}

func witdrawGameScreenLogics() {

	var symbol string = " "
	for x := 0; x < xP; x++ {
		for y := 0; y < yP; y++ {
			symbol = calculateSnakeBodyPositionSymbol(x, y)

			if food[0] == x && food[1] == y && !snakeEat {
				symbol = "●"
			}

			screenLines[x][y] = symbol

			if snakeEat {
				setNewFoodPositin()
			}

		}
	}

	for l := 0; l < len(screenLines); l++ {
		fmt.Println("|" + fmt.Sprint(strings.Join(screenLines[l][:], "")) + "☰")
	}
}

func calculateSnakeBodyPositionSymbol(x int, y int) (symbol string) {
	//Вычитываем есть ли на этой клетке змея
	symbol = " "
	for snakeBody := len(snakePositions) - 1; snakeBody >= 0; snakeBody-- {
		if snakePositions[snakeBody][0] == x && snakePositions[snakeBody][1] == y {
			if snakeBody == len(snakePositions)-1 {
				if food[0] == x && food[1] == y {
					snakeEat = true
				}
				symbol = "□"
			} else {
				symbol = "■"
			}
		}

		if snakePositions[snakeBody][0] > xP-1 ||
			snakePositions[snakeBody][0] < 0 ||
			snakePositions[snakeBody][1] > yP-1 ||
			snakePositions[snakeBody][1] < 0 {
			//Вышли за поле
			startGame = false
		}

		checkOnIntersectionBody(snakeBody)

	}
	return
}

func setNewFoodPositin() {
	//Высчитываем новую позицию еды
	freeBlockForFood := [][]int{}

	for x := 0; x < xP; x++ {
		for y := 0; y < yP; y++ {
			if screenLines[x][y] == " " {
				freeBlockForFood = append(freeBlockForFood, []int{x, y})
			}
		}
	}

	freePositionRand := rand.Intn(len(freeBlockForFood)-0) + 0

	food[0] = freeBlockForFood[freePositionRand][0]
	food[1] = freeBlockForFood[freePositionRand][1]
}

func endLoopPrintInfo() {
	fmt.Println("☰☰☰☰☰☰☰☰☰☰☰☰☰☰☰☰☰☰")
	fmt.Println("Очки :", point, "Длина : ", len(snakePositions))
	time.Sleep(750 * time.Millisecond)
}

func checkOnIntersectionBody(snakeBody int) {
	if snakeBody == len(snakePositions)-1 {
		for snakeCheckIntersection := 0; snakeCheckIntersection < len(snakePositions); snakeCheckIntersection++ {
			if snakeCheckIntersection != len(snakePositions)-1 {
				if snakePositions[snakeBody][0] == snakePositions[snakeCheckIntersection][0] &&
					snakePositions[snakeBody][1] == snakePositions[snakeCheckIntersection][1] {
					startGame = false
					//Съели себя
				}
			}

		}
	}
}

func oneStepForward() {
	//Высчитываем куда ползет змейка
	var xP int = snakePositions[len(snakePositions)-1][0]
	var yP int = snakePositions[len(snakePositions)-1][1]

	switch key := keyPress; key {
	case rightBtn:
		yP += 1
	case leftBtn:
		yP -= 1
	case topBtn:
		xP -= 1
	default:
		xP += 1
	}

	snakePositions = append(snakePositions, []int{xP, yP})

	if !snakeEat {
		//Удаляем последнее тело, если не ели
		snakePositions = snakePositions[1:]
	} else {
		point += 10
	}

	snakeEat = false
}

func TrackClickGo() {
	//Отслеживаем нажатие на клавиатуру
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if !(string(key) == leftBtn && keyPress == rightBtn ||
			string(key) == rightBtn && keyPress == leftBtn ||
			string(key) == topBtn && keyPress == downBtn ||
			string(key) == downBtn && keyPress == topBtn) {
			keyPress = string(key)
			keyPressChar = string(char)
		}
	}
}
