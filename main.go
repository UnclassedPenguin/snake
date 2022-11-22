package main

import (
  "os"
  "fmt"
  "time"
  "math/rand"
  "github.com/gdamore/tcell/v2"
)

type Pos struct {
  x int
  y int
}

type Snake struct {
  char rune
  x int
  y int
  xspeed int
  yspeed int
  length int
  tail []Pos
  food Pos
}

func draw(snake Snake, s tcell.Screen, style tcell.Style) {
  s.SetContent(snake.x, snake.y, snake.char, []rune{}, style)
  s.SetContent(snake.food.x, snake.food.y, '@', []rune{}, style)
}

func newFood(s tcell.Screen, style tcell.Style) Pos {
  rand.Seed(time.Now().UnixNano())

  x, y := s.Size()

  var food Pos
  food.x = rand.Intn(x)
  food.y = rand.Intn(y)

  return food
}

func checkFood(snake Snake) bool {
  if snake.x == snake.food.x && snake.y == snake.food.y {
    return true
  } else {
    return false
  }
}

func update(snake *Snake, s tcell.Screen, style tcell.Style){
  s.Clear()
  food := checkFood(*snake)
  if food {
    snake.food = newFood(s, style)
  }
  draw(*snake, s, style)
  s.Sync()
}

func main() {
  s, err := tcell.NewScreen()
  if err != nil {
    fmt.Println("Error: ", err)
    os.Exit(1)
  }

  if err := s.Init(); err != nil {
    fmt.Println("Error initializing New Screen: ", err)
    os.Exit(1)
  }

  s.Clear()

  //style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
  style := tcell.StyleDefault.Foreground(tcell.ColorBlue)

  x, y := s.Size()

  // Create initial snake
  var snake Snake
  snake.char = tcell.RuneBlock
  snake.x = x/2
  snake.y = y/2
  snake.xspeed = 1
  snake.yspeed = 0
  snake.length = 0
  snake.tail = []Pos{}
  snake.food = newFood(s, style)

  // Draw snake
  draw(snake, s, style)

  for {
    switch ev := s.PollEvent().(type) {
    case *tcell.EventResize:
      s.Sync()
    case *tcell.EventKey:
      switch ev.Key() {
      case tcell.KeyCtrlC, tcell.KeyEscape:
        s.Fini()
        os.Exit(0)
      case tcell.KeyRune:
        switch ev.Rune() {
        case 'J', 'j':
          snake.y += 1
          update(&snake, s, style)
        case 'K', 'k':
          snake.y -= 1
          update(&snake, s, style)
        case 'H', 'h':
          snake.x -= 1
          update(&snake, s, style)
        case 'L', 'l':
          snake.x += 1
          update(&snake, s, style)
        }
      default:
        fmt.Println(ev.Key())
      }
    }
  }
}
