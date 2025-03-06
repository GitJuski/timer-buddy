#include <stdio.h>
#include <time.h>
#include <stdbool.h>
#include "clicalendar.h"

// Check if the supplied year is a leap year or not
bool isLeap(int* year) {
  if (*year % 4 == 0 && (*year % 100 != 0 || *year % 400 == 0)) {
    return true;
  } else {
    return false;
  }
}

// Function for printing out the calendar. Originally I had ptrs to these values a parameter but this way we don't have to use Unsafe in GO.
void printCalendar(int month, int year) {
  int daysInMonth[12] = {31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31};
  if (isLeap(&year)) {
    daysInMonth[1] = 29;
  }

  const char* monthNames[] = {"January", "February", "March", "april", "May", "June", "july", "August", "September", "October", "November", "december"};

  printf("\n    %s %d\n", monthNames[month - 1], year); // We have to substract 1 from the supplied month since arrays start at 0 so January is 0
  printf("Su Mo Tu We Th Fr Sa\n");

  struct tm timeinfo = {0}; // New instance of tm struct
  timeinfo.tm_year = year - 1900; // Set the tm_year value to year - 1900 because tm_year is designed to hold an offset from 1900. 2025 - 1900 = 125. Later when we use mktime() it uses this 125 value
  timeinfo.tm_mon = month - 1; // tm_mon has the same thing going as the monthNames. So January is 0 etc..
  timeinfo.tm_mday = 1; // Set the day of the month to 1
  mktime(&timeinfo); // Converts the struct to time_t. It gets the tm_wday for example that we'll use

  // After calling mktime() The day of the month is some day of the week -> 0: sunday, 1: Monday etc... The following loop prints out 3 whitespaces as many times as the weekday index is. For example if the The first day of the month is a Tuesday it will print 2 times 3 whitespaces because Tuesday is index 2. So the result is a calendar view where the first day starts at the position of the Tuesday in line 25. Tried my best to explain this.
  for (int i = 0; i < timeinfo.tm_wday; i++) {
    printf("   ");
  }
  // Loop through the other days in month
  for (int day = 1; day <= daysInMonth[month - 1]; day++) {
    printf("%02d ", day);
    // If the next printable day is Sunday -> Prints a newline
    if ((day + timeinfo.tm_wday) % 7 == 0) {
      printf("\n");
    }
  }
  printf("\n\n");
}
