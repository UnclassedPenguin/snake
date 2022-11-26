//-------------------------------------------------------------------------------
//-------------------------------------------------------------------------------
//
// Tyler(UnclassedPenguin) Snake 2022
//
// Author: Tyler(UnclassedPenguin)
//    URL: https://unclassed.ca
// GitHub: https://github.com/UnclassedPenguin
//
//-------------------------------------------------------------------------------
//-------------------------------------------------------------------------------

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
  difficulty int
}

func menu(s tcell.Screen, style tcell.Style, snake Snake) {
  x, y := s.Size()
  strings := []string{
    "UnclassedPenguin Snake",
    "Press any key to start game",
    "Esc or Ctrl-C to quit",
    fmt.Sprintf("Difficulty: %v", snake.difficulty),
    "Press 1, 2, or 3 to set difficulty",
    "1 being easiest and 3 being hardest",
  }

  for i, str := range strings {
    writeToScreen(s,style,((x/2)-(len(str)/2)),y/3+(i*2),str)
  }

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
        case '1':
          snake.difficulty = 1
          strings[3] = fmt.Sprintf("Difficulty: %v", snake.difficulty)
          writeToScreen(s,style,((x/2)-(len(strings[3])/2)),y/3+6,strings[3])
          s.Sync()
        case '2':
          snake.difficulty = 2
          strings[3] = fmt.Sprintf("Difficulty: %v", snake.difficulty)
          writeToScreen(s,style,((x/2)-(len(strings[3])/2)),y/3+6,strings[3])
          s.Sync()
        case '3':
          snake.difficulty = 3
          strings[3] = fmt.Sprintf("Difficulty: %v", snake.difficulty)
          writeToScreen(s,style,((x/2)-(len(strings[3])/2)),y/3+6,strings[3])
          s.Sync()
        default:
          game(s, style, snake)
        }
      default:
        game(s, style, snake)
      }
    }
  }

}

func game(s tcell.Screen, style tcell.Style, snake Snake) {
  x, y := s.Size()
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
          if snake.length > 0 {
            fmt.Println("Score:", snake.length-1)
          }
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
          //case ' ':
            //snake.xspeed = 0
            //snake.yspeed = 0
          }
        }
      }
    }
  }()

  var delay int
  switch snake.difficulty {
  case 1:
    delay = 200
  case 2:
    delay = 100
  case 3:
    delay = 50
  }

  // Main loop
  for {

    snake.x += snake.xspeed
    snake.y += snake.yspeed
    if snake.x < 0 || snake.x > x || snake.y < 0 || snake.y > y-1 {
      gameOver(s, snake)
    }
    for _, v := range snake.tail {
      if snake.x == v.x && snake.y == v.y {
        gameOver(s, snake)
      }
    }

    update(&snake, s, style)
    time.Sleep(time.Millisecond * time.Duration(delay))
    }
}

func draw(snake Snake, s tcell.Screen, style tcell.Style) {
  writeToScreen(s,style,1,1,fmt.Sprintf("Score: %v", snake.length-1))
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
  fmt.Printf("Score: %v\n", snake.length-1)
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
  snake.difficulty = 1

  menu(s, style, snake)
}
