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
  history []Pos
  food Pos
}

func menu(s tcell.Screen, style tcell.Style) {
  x, y := s.Size()
  str1 := "UnclassedPenguin Snake"
  str2 := "Press any key to start game"
  str3 := "Esc or Ctrl-C to quit"

  writeToScreen(s,style,((x/2)-(len(str1)/2)),y/3,str1)
  writeToScreen(s,style,((x/2)-(len(str2)/2)),y/3+2,str2)
  writeToScreen(s,style,((x/2)-(len(str3)/2)),y/3+4,str3)

  for {
    switch ev := s.PollEvent().(type) {
    case *tcell.EventResize:
      s.Sync()
    case *tcell.EventKey:
      switch ev.Key() {
      case tcell.KeyCtrlC, tcell.KeyEscape:
        s.Fini()
        os.Exit(0)
      default:
        game(s, style)
      }
    }
  }

}

func game(s tcell.Screen, style tcell.Style) {
  x, y := s.Size()

  // Create initial snake
  var snake Snake
  snake.char = tcell.RuneBlock
  snake.x = x/2
  snake.y = y/2
  snake.xspeed = 1
  snake.yspeed = 0
  snake.length = 1
  snake.history = []Pos{}
  snake.tail = []Pos{}
  snake.food = newFood(s, style)

  // Handles Keyboard input
  go func() {
    for {
      switch ev := s.PollEvent().(type) {
      case *tcell.EventResize:
        s.Sync()
      case *tcell.EventKey:
        switch ev.Key() {
        case tcell.KeyCtrlC, tcell.KeyEscape:
          s.Fini()
          os.Exit(0)
        case tcell.KeyDown:
          snake.yspeed = 1
          snake.xspeed = 0
        case tcell.KeyUp:
          snake.yspeed = -1
          snake.xspeed = 0
        case tcell.KeyLeft:
          snake.xspeed = -1
          snake.yspeed = 0
        case tcell.KeyRight:
          snake.xspeed = 1
          snake.yspeed = 0
        case tcell.KeyRune:
          switch ev.Rune() {
          case 'J', 'j':
            snake.yspeed = 1
            snake.xspeed = 0
          case 'K', 'k':
            snake.yspeed = -1
            snake.xspeed = 0
          case 'H', 'h':
            snake.xspeed = -1
            snake.yspeed = 0
          case 'L', 'l':
            snake.xspeed = 1
            snake.yspeed = 0
          }
        default:
          fmt.Println(ev.Key())
        }
      }
    }
  }()

  // Main loop
  for {
    snake.x += snake.xspeed
    snake.y += snake.yspeed
    if snake.x < 0 || snake.x > x || snake.y < 0 || snake.y > y {
      gameOver(s, snake)
    }
    for _, v := range snake.tail {
      if snake.x == v.x && snake.y == v.y {
        gameOver(s, snake)
      }
    }
    update(&snake, s, style)
    time.Sleep(time.Millisecond * 200)
  }

}

func draw(snake Snake, s tcell.Screen, style tcell.Style) {
  s.SetContent(snake.x, snake.y, snake.char, []rune{}, style)
  for i, _ := range snake.history {
    s.SetContent(snake.history[i].x, snake.history[i].y, snake.char, []rune{}, style)
  }
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

func writeToScreen(s tcell.Screen, style tcell.Style, x int, y int, str string) {
  for i, char := range str {
    s.SetContent(x+i, y, rune(char), []rune{}, style)
  }
}

func update(snake *Snake, s tcell.Screen, style tcell.Style){
  s.Clear()
  food := checkFood(*snake)
  newHistory := Pos{snake.x, snake.y}
  snake.history = append(snake.history, newHistory)
  snake.history = snake.history[len(snake.history)-snake.length:len(snake.history)]
  if len(snake.history) > 0 {
    snake.tail = snake.history[:len(snake.history)-1]
  }

  if food {
    snake.length++
    snake.food = newFood(s, style)
  }
  draw(*snake, s, style)
  s.Sync()
  //fmt.Println("X:", snake.x)
  //fmt.Println("Y:", snake.y)
  //fmt.Println(snake.history)
  //fmt.Println("TAIL: ", snake.tail)
  //fmt.Println(snake.length)
}

func gameOver(s tcell.Screen, snake Snake) {
  s.Fini()
  fmt.Println("You Lose!")
  fmt.Printf("Score: %v\n", len(snake.tail))
  os.Exit(0)
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

  style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

  menu(s, style)
}
