package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/GitJuski/timer-buddy/dbhandler"
)

//#include "clicalendar.h"
import "C"

// ANSI Escapes
const (
  GREEN_BG = "\033[42m"
  BLACK_TEXT = "\033[30m"
  COLOR_RED = "\033[31m"
  COLOR_ORANGE = "\033[33m"
  COLOR_RESET = "\033[0m"
  CURSOR_DOWN = "\033[3B" // Down 3
  CURSOR_UP = "\033[3A" // Up 3
  CURSOR_START = "\033[H"
  HIDE_CURSOR = "\033[?25l"
  SHOW_CURSOR = "\033[?25h"
  BELL = "\a" // c-escape for bell ASCII
)

// Changed to uint since values shoudn't be under 0. ParseUint doesn't return uint so we have to use uint64
type Timer struct {
  Hours uint64
  Minutes uint64
  Seconds uint64
}

var stop bool = false
var timerMode string

// My ascii masterpiece a robot of some sort(Create your own ascii buddy. It's a slice so you can add different faces)
var asciiArt = []string {`
     ###########      
    ##._.###._.##     
   ##/ o \#/ o \##    
  ###\___/#\___/###   
  #################   
   ######   ######    
    #############     
     ###########      
       /*****\        
      @_______@       
     /|*******|\      
    / |_______| \     
       _/  \_         
`, `
     ###########      
    ##._.###._.##     
   ##/ > \#/ < \##    
  ###\___/#\___/###   
  #################   
   #####\___/#####    
    #############     
     ###########      
    \  /*****\  /     
     \@_______@/      
      |*******|       
      |_______|       
        |  |          
`}

// Function for handling the operations in the interactive mode
func handleInteractive() {
  fmt.Print(SHOW_CURSOR)
  var month int
  var year int
  var day int

  for {
    fmt.Print("Month (1-12): ")
    fmt.Scanln(&month)
    fmt.Print("Year (yyyy): ")
    fmt.Scanln(&year)

    C.printCalendar(C.int(month), C.int(year)) // Use the C code to print the calendar of the inputted month and year
    times := dbhandler.GetTimes(fmt.Sprintf("%02d", month), fmt.Sprintf("%d", year)) // Query the database for that month's data

    for {
      fmt.Print("Day to check (Insert 0 to go back): ")
      fmt.Scanln(&day)
      if day == 0 {
        break
      }

      for i := range times {
        // If the supplied date matches the structs Date field's last two chars "yyyy-mm-dd" <- The day. Then it prints the values stored
        if fmt.Sprintf("%02d", day) == times[i].Date[8:] {
          fmt.Printf("\n%sID:\t%d\n", COLOR_ORANGE, times[i].ID)
          fmt.Printf("%sLength:\t%s\n", COLOR_ORANGE, times[i].Length)
          fmt.Printf("%sDate:\t%s\n", COLOR_ORANGE, times[i].Date)
          fmt.Printf("%sNote:\t%s\n", COLOR_ORANGE, times[i].Note)
          fmt.Println(COLOR_RESET)
        }
      }
    }
  }
}

// Presents the option to save the time into SQLITE
func handleSave(length string, currentDate string) {
  var choice string

  fmt.Print("\nSave the time (Y/n) ")
  fmt.Scanln(&choice)

  switch(strings.ToLower(choice)) {
  case "", "y":
    // We use bufio to get an input that can have spaces. fmt does not allow this it breaks the input in whitespace
    input := bufio.NewReader(os.Stdin)
    fmt.Print("\nWrite a note (leave empty to not write one) ")
    note, err := input.ReadString('\n')
    if err != nil {
      log.Fatal(err)
    }
    note = strings.TrimSpace(note) // Remove leading and trailing whitepaces and the newline at the end

    dbhandler.InsertTime(length, currentDate, note)
  case "n":
    fmt.Println("DONT INSERT")
  default:
    log.Fatal("Give a proper value")
  }
}

// Stopwatch functionality. Stop it with Ctrl+C since there is no stopping functionality atm. There is a Ctrl+C handling in other goroutine
func stopWatch() string {

  timer := Timer{0, 0, 0} // Initialize a timer variable from the Timer struct.
  
// Golangs while loop
for {
  if stop {
    break
  }
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
  return fmt.Sprintf("%02d:%02d:%02d", timer.Hours, timer.Minutes, timer.Seconds)
}

// Countdown timer functionality. Gets a pointer to a string slice as a parameter so we don't have to create a copy. I'm pretty sure passing arrays as arguments works a bit differently than in C which is why I did this. I could be wrong tho. In C an array as parameter is a pointer to the array[o]
func countdown(parts *[]string) string {
  hours, err := strconv.ParseUint((*parts)[0], 10, 0)
  if err != nil { 
    log.Fatal(err)
  }
  minutes, err := strconv.ParseUint((*parts)[1], 10, 0)
  if err != nil {
    log.Fatal(err)
  }
  seconds, err := strconv.ParseUint((*parts)[2], 10, 0)
  if err != nil {
    log.Fatal(err)
  }
  timer := Timer{hours, minutes, seconds} // A new instance of Timer with the inputted values
  
// Golang's while loop
  for {
    if stop {
      break
    }
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
      stop = true
      fmt.Println(SHOW_CURSOR) // Show cursor again
      break
    }
    // Simple countdown timer functionality -> decrease seconds, do checks etc...
    if timer.Seconds == 0 && timer.Minutes > 0 {
      timer.Minutes --
      timer.Seconds = 59
    }
    if timer.Minutes == 0 && timer.Hours > 0 {
      timer.Hours --
      timer.Minutes = 59
      timer.Seconds = 59
    }
    fmt.Printf("\r%02d:%02d:%02d", timer.Hours, timer.Minutes, timer.Seconds)
    time.Sleep(time.Second)
    timer.Seconds --
  }
  return fmt.Sprintf("%02d:%02d:%02d", timer.Hours, timer.Minutes, timer.Seconds)
}

func main() {
  // Command line argument. Returns a string pointer.
  countTime := flag.String("t", "00:00:00", "Insert time in hh:mm:ss format")
  interactiveMode := flag.Bool("i", false, "Enters to interactive mode.")
  flag.Parse()
  // This is for checking Ctrl+C termination and then doing cleanup.
  channel := make(chan os.Signal, 1) // Creates a channel with os.signal types with buffer size 1
  signal.Notify(channel, os.Interrupt) //Go runtime sends os.Interrupt signals to the channel
  go func () { // Start a goroutine that waits for data from the channel and do the following.
    <-channel // Wait for data from the channle.
    stop = true
    fmt.Print(SHOW_CURSOR)
    fmt.Println("\nExiting...")
    if *interactiveMode {
      os.Exit(0)
    }
  }() // Start the goroutine right away

  dbhandler.CreateDB() // Creates the database and table if they don't exist

  fmt.Print(HIDE_CURSOR)

  if *interactiveMode {
    handleInteractive()
  } else {
    // Split the string into slice and check that it's length is 3.
    parts := strings.Split(*countTime, ":")
    if len(parts) != 3 {
      fmt.Print(SHOW_CURSOR)
      log.Fatal("Your input is incorrectly formatted")
    }

    var stopwatchTime string
    stopwatchResultChan := make(chan string) // A channel for getting the return value of the stopWatch function that is ran in goroutine

    var countdownTime string
    countdownResultChan := make(chan string) // A channel for getting the return value of the countdown function that is ran in goroutine
    
    // Dereference the value and start the stopWatch if no argument were given or argument 00:00:00.
    if *countTime == "00:00:00" {
      // Start the stopWatch in goroutine in a way that gets the return value with the use of the channel from before
      go func() {
        stopwatchTime = stopWatch() 
        stopwatchResultChan <- stopwatchTime
      }()
      timerMode = "stopWatch"
    } else {
      // Same as the above but with the countdown function
      go func() {
        countdownTime = countdown(&parts) 
        countdownResultChan <- countdownTime
      }()
      timerMode = "countdown"
    }

    // A while loop for the ascii art
    for {
      if stop {
        break
      }
      for _, art := range asciiArt {
        fmt.Printf("%s%s%s%s%s", CURSOR_START, GREEN_BG, BLACK_TEXT, art, COLOR_RESET)
        time.Sleep(2 * time.Second)
      }
    }
    currentDate := time.Now().Format(time.DateOnly) // SQLITE doesn't have time types so we'll get the current date in string

    if timerMode == "stopWatch" {
      stopwatchTime = <- stopwatchResultChan // Save the return value of stopWatch function to stopwatchTime
      handleSave(stopwatchTime, currentDate)
    } else if timerMode == "countdown" {
      countdownTime = <- countdownResultChan // Save the return value of countdown function to countdownTime
      // If the countdown didn't run till the end (If for example you stop it prematurely with Ctrl+C)
      if countdownTime != "00:00:00" {
        // Split the inputted time and the time at the moment you stopped the program into their own slices
        timePassedSplit := strings.Split(countdownTime, ":")
        setTimeSplit := strings.Split(*countTime, ":")
        // Wanted to do most of this in one line, but since Golang's err handling is what it is, I had to do it this way. No hate to Golang err handling I like it most of the time
        setHours, err := strconv.Atoi(setTimeSplit[0])
        if err != nil {
          log.Fatal(err)
        }
        setMinutes, err := strconv.Atoi(setTimeSplit[1])
        if err != nil {
          log.Fatal(err)
        }
        setSeconds, err := strconv.Atoi(setTimeSplit[2])
        if err != nil {
          log.Fatal(err)
        }
        passedHours, err := strconv.Atoi(timePassedSplit[0])
        if err != nil {
          log.Fatal(err)
        }
        passedMinutes, err := strconv.Atoi(timePassedSplit[1])
        if err != nil {
          log.Fatal(err)
        }
        passedSeconds, err := strconv.Atoi(timePassedSplit[2])
        if err != nil {
          log.Fatal(err)
        }

        // Both times in seconds
        totalSetTime := ((setHours * 60) * 60) + (setMinutes * 60) + setSeconds
        totalPassedTime := ((passedHours * 60) * 60) + (passedMinutes * 60) + passedSeconds

        // Transfer the substraction into hours, minutes and seconds
        n := totalSetTime - totalPassedTime
        hours := n / 3600
        n %= 3600
        minutes := n / 60
        n %= 60
        seconds := n

        realPassedTime := fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds) // Format a new string "HH:MM:SS" of the real passed time
        handleSave(realPassedTime, currentDate)
      } else {
        handleSave(*countTime, currentDate) // If you let the countdown run till the end it just saves the time you inputted
      } 
    } else {
      log.Fatal("Something went wrong")
    }
  }
}
