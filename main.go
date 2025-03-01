package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

// ANSI Escapes
const (
  GREEN_BG = "\033[42m"
  BLACK_TEXT = "\033[30m"
  COLOR_RED = "\033[31m"
  COLOR_RESET = "\033[0m"
  CURSOR_DOWN = "\033[3B" // Down 3
  CURSOR_UP = "\033[3A" // Up 3
  CURSOR_START = "\033[H"
  HIDE_CURSOR = "\033[?25l"
  SHOW_CURSOR = "\033[?25h"
  BELL = "\a" // c-escape for bell ASCII
)

type Timer struct {
  Hours int
  Minutes int
  Seconds int
}

// var asciiArt = []string{">.<", "-.-", "^o^"}
//var asciiArt = []string {`
// /\_/\ 
//( O.O )
// <   > 
//`, `
// /\_/\ 
//( >.< )
// >   < 
//`}
// My ascii masterpiece (Create your own ascii buddy. It's a slice so you can add different faces)
var asciiArt = []string {`
 ()___()
 /     \
| O  O |
|  ~   |
 \____/ 
`, `
 ()___()
 /     \
| >  < |
|  .   |
 \____/ 
`}

// Stopwatch functionality. Stop it with Ctrl+C since there is no stopping functionality atm. There is a Ctrl+C handling in other goroutine
func stopWatch() {

  timer := Timer{0, 0, 0} // Initialize a timer variable from the Timer struct.
  
// Golangs while loop
for {
  fmt.Printf("\r%02d:%02d:%02d", timer.Hours, timer.Minutes, timer.Seconds, ) // Prints the values of the timer instance formatted HH:MM:SS
  // Basic stopwatch functionality. Add one second every second, add minute every 60 seconds etc...
    time.Sleep(time.Second) 
    timer.Seconds ++
    if timer.Seconds == 60 {
      timer.Minutes ++
      timer.Seconds = 0
    }
    if timer.Minutes == 60 {
      timer.Hours ++
      timer.Minutes = 0
    }
  }
}

// Countdown timer functionality. Gets a pointer to a string slice as a parameter so we don't have to create a copy. I'm pretty sure passing arrays as arguments works a bit differently than in C which is why I did this. I could be wrong tho. In C an array as parameter is a pointer to the array[o]
func countdown(parts *[]string) {
  hours, err := strconv.Atoi((*parts)[0]) // Same as strconv.ParseInt(string, 10, 0) where 10 is base and 0 is bits. Used ParseInt with 64 bit, int64 
  if err != nil { 
    log.Fatal(err)
  }
  minutes, err := strconv.Atoi((*parts)[1])
  if err != nil {
    log.Fatal(err)
  }
  seconds, err := strconv.Atoi((*parts)[2])
  if err != nil {
    log.Fatal(err)
  }
  timer := Timer{hours, minutes, seconds} // A new instance of Timer with the inputted values
  
// Golang's while loop
  for {
    // If all values are 0 -> timer ends
    if timer.Seconds == 0 && timer.Minutes == 0 && timer.Hours == 0 {
      // A for loop for three times to blink red and ring terminal bell
      for i:= 0; i < 3; i++ {
        fmt.Printf("\r%s00:00:00", COLOR_RED)
        fmt.Print("\a")
        time.Sleep(time.Second)
        fmt.Printf("\r%s00:00:00", COLOR_RESET)
        time.Sleep(time.Second)
      }
      fmt.Print(SHOW_CURSOR) // Show cursor again
      os.Exit(0) // Exit with exit status 0 meaning a normal exit
    }
    // Simple countdown timer functionality -> decrease seconds, do checks etc...
    if timer.Seconds == 0 && timer.Minutes > 0 {
      timer.Minutes --
      timer.Seconds = 59
    }
    if timer.Minutes == 0 && timer.Hours > 0 {
      timer.Hours --
      timer.Minutes = 59
    }
    fmt.Printf("\r%02d:%02d:%02d", timer.Hours, timer.Minutes, timer.Seconds)
    time.Sleep(time.Second)
    timer.Seconds --
  }
}

func main() {
  // This is for checking Ctrl+C termination and then doing cleanup.
  channel := make(chan os.Signal, 1) // Creates a channel with os.signal types with buffer size 1
  signal.Notify(channel, os.Interrupt) //Go runtime sends os.Interrupt signals to the channel
  go func () { // Start a goroutine that waits for data from the channel and do the following.
    <-channel // Wait for data from the channle.
    fmt.Print(SHOW_CURSOR)
    fmt.Println("\nExiting...")
    os.Exit(1) // Exit with code 1 -> unregular exit.
  }() // Start the goroutine right away

  fmt.Print(HIDE_CURSOR)
  // Command line argument. Returns a string pointer.
  countTime := flag.String("t", "00:00:00", "Insert time in hh:mm:ss format")
  flag.Parse()

  // Split the string into slice and check that it's length is 3.
  parts := strings.Split(*countTime, ":")
  if len(parts) != 3 {
    fmt.Print(SHOW_CURSOR)
    log.Fatal("Your input is incorrectly formatted")
  }
  
  // Dereference the value and start the stopWatch if no argument were given or argument 00:00:00.
  if *countTime == "00:00:00" {
    go stopWatch() // Start the stopWatch as a goroutine (separate thread).
  } else {
    go countdown(&parts) // Start the countdown as a goroutine and pass the slice's memory location as parameter.
  }

  // A while loop for the ascii art
  for {
    for _, art := range asciiArt {
      fmt.Printf("%s%s%s%s%s", CURSOR_START, GREEN_BG, BLACK_TEXT, art, COLOR_RESET)
      time.Sleep(2 * time.Second)
    }
  }
}
