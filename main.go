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
  sizex int
  sizey int
  length int
  history []Pos
  food Pos
  difficulty int
  style tcell.Style
}

// Shows the menu at the start
func menu(s tcell.Screen, snake Snake) {
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
    writeToScreen(s,snake.style,((x/2)-(len(str)/2)),y/3+(i*2),str)
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
          writeToScreen(s,snake.style,((x/2)-(len(strings[3])/2)),y/3+6,strings[3])
          s.Sync()
        case '2':
          snake.difficulty = 2
          strings[3] = fmt.Sprintf("Difficulty: %v", snake.difficulty)
          writeToScreen(s,snake.style,((x/2)-(len(strings[3])/2)),y/3+6,strings[3])
          s.Sync()
        case '3':
          snake.difficulty = 3
          strings[3] = fmt.Sprintf("Difficulty: %v", snake.difficulty)
          writeToScreen(s,snake.style,((x/2)-(len(strings[3])/2)),y/3+6,strings[3])
          s.Sync()
        default:
          game(s, snake)
        }
      default:
        game(s, snake)
      }
    }
  }
}

// The main game section
func game(s tcell.Screen, snake Snake) {
  // Handles Keyboard input
  go func() {
    for {
      switch ev := s.PollEvent().(type) {
      case *tcell.EventResize:
        s.Sync()
      case *tcell.EventKey:
        switch ev.Key() {
        case tcell.KeyCtrlC, tcell.KeyEscape:
          gameExit(s, snake.length-1)
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
          case 'Q', 'q':
            gameExit(s, snake.length-1)
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

    checkLose(s, snake)

    update(s, &snake)
    time.Sleep(time.Millisecond * time.Duration(delay))
  }
}

// Checks if the head of the snake has hit the wall, or ran
// into itself.
func checkLose(s tcell.Screen, snake Snake) {
  x, y := s.Size()
  if snake.x < 0 || snake.x > x || snake.y < 0 || snake.y > y-1 {
      gameOver(s, snake.length-1)
  }

  if len(snake.history) > 1 {
    for _, v := range snake.history[:len(snake.history)-1] {
      if snake.x == v.x && snake.y == v.y {
        gameOver(s, snake.length-1)
      }
    }
  }
}

// Updates the snake. Checks if it has eaten, if it has, grows by one. 
func update(s tcell.Screen, snake *Snake){
  food := checkFood(*snake)

  newHistory := Pos{snake.x, snake.y}
  snake.history = append(snake.history, newHistory)
  snake.history = snake.history[len(snake.history)-snake.length:len(snake.history)]

  if food {
    snake.length++
    snake.food = newFood(s, snake.style)
  }

  draw(s, *snake)
}

func draw(s tcell.Screen, snake Snake) {
  s.Clear()
  // Diagnostics:
  //strings := []string{
    //fmt.Sprintf("X:%v",snake.x),
    //fmt.Sprintf("Y:%v",snake.y),
    //fmt.Sprintf("Hist:%v",snake.history),
    //fmt.Sprintf("length:%v",snake.length),
  //}

  //for i, str := range strings {
    //writeToScreen(s,snake.style,1,2+i,str)
  //}

  writeToScreen(s,snake.style,1,1,fmt.Sprintf("Score: %v", snake.length-1))
  for i, _ := range snake.history {
    s.SetContent(snake.history[i].x, snake.history[i].y, snake.char, nil, snake.style)
  }
  s.SetContent(snake.food.x, snake.food.y, '@', nil, snake.style)
  s.Sync()
}

// Picks a new random position for the food. Should
// maybe update it so that it doesn't pick a spot that
// the snake occupies.
func newFood(s tcell.Screen, style tcell.Style) Pos {
  rand.Seed(time.Now().UnixNano())

  x, y := s.Size()

  var food Pos
  food.x = rand.Intn(x)
  food.y = rand.Intn(y)

  return food
}

// Check if the head of snake has gotten the food.
func checkFood(snake Snake) bool {
  return snake.x == snake.food.x && snake.y == snake.food.y
}

// Write a string to the screen.
func writeToScreen(s tcell.Screen, style tcell.Style, x int, y int, str string) {
  for i, char := range str {
    s.SetContent(x+i, y, rune(char), []rune{}, style)
  }
}

// Function to run if user quits
func gameExit(s tcell.Screen, score int) {
  s.Fini()
  fmt.Printf("Thanks for playing!\nScore: %v\n", score)
  os.Exit(0)
}

// Function to run when game over.
func gameOver(s tcell.Screen, score int) {
  s.Fini()
  fmt.Printf("You Lose!\nScore: %d\n", score)
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
  // Starts in center, heading right
  snake := Snake{
    char: tcell.RuneBlock,
    x: x/2,
    y: y/2,
    xspeed: 1,
    yspeed: 0,
    sizex: x,
    sizey: y,
    length: 1,
    history: []Pos{},
    food: newFood(s, style),
    difficulty: 1,
    style: style,
  }

  menu(s, snake)
}
