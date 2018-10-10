# Go_Supreme
Buys some supreme stuff.

## Getting Started
1. Make sure you have a task file somewhere to use on the command line.
2. Build and go. Pretty easy.
~~~~
go build
./supreme ./pathto/tasks.json ./optionalFilePath/Settings.json
~~~~

### Task Setup
The task file is plain json.
```
[
  {
    "taskName": "shopsafe 0 0RU3",
    "item": {
      "keywords": [
        "Briefs",
        "Boxer"
      ],
      "size": "medium",
      "color": "white",
      "category": "accessories"
    },
    "account": {
      "person": {
        "firstname": "Jax",
        "lastname": "Blax",
        "email": "none+0RU3@gmail.com",
        "phoneNumber": "354-143-9568"
      },
      "address": {
        "address1": "0RU3 123 HoneySuckle Ave",
        "address2": "",
        "zipcode": "85542",
        "city": "Springfield",
        "state": "WA",
        "country": "USA"
      },
      "card": {
        "cardtype": "notneeded",
        "number": "1234 2541 2154 5487",
        "month": "09",
        "year": "2022",
        "cvv": "789"
      }
    }
  }
]
```

### Settings Setup
```
{
  "refreshWait": 300,
  "atcWait": 800,
  "checkoutWait": 800,
}
```

## Building Versions
To build different versions you will need to setup goreleaser using brew or some other installer. Then:
~~~~
goreleaser --snapshot
~~~~

### Build Windows Only
```
GOOS=windows GOARCH=386 go build -o supreme-storecredit.exe
```

## TODO:
### Current
* Add scheduling
* Review security code
  * Add key versioning
  * Add date added

### Pipeline
* Unify / pool initial item search
* Test if I can add cookie and skip ATC
* Clean up code and model an interface for mobile and desktop
* Add proxy support for each task
* Add any size keyword
* https://sequencediagram.org/ Diagram calls
* UI Text
  * https://github.com/gizak/termui
  * https://github.com/jroimartin/gocui
  * https://github.com/gdamore/tcell
* Full UI Version
  * Add ability to kill one of these go routines via select statement from either cancel channel
      * https://chilts.org/2017/06/12/cancelling-multiple-goroutines
  * Add gcapture code in case they revert
* Command line to feed in file different commands
  * https://github.com/spf13/cobra
* Restock monitor
* Increased security
  * https://nucleus.sh/docs/sell - For the entire js application
  * https://stackoverflow.com/questions/25062696/what-about-protection-for-golang-source-code


### Completed
* Finish API - 9/19
  * Finish queuing, see js version - 9/19
* Add worker pool - 9/20
* Add task importing - 9/20
* Create task generator in json - 9/20 (in python)
* Add better waiting functionality - 9/20
* Attempt to purge credentials from old commits - 9/23
* Clean up testing - 9/23
* Add go module system - 9/23
* Add better logging - 9/23
* Fix weird task struct formatting - 9/24
  * Update go code - 9/24
  * Update python code - 9/24
* Add a verify task function - 9/24
* Add better retrying - 9/24
* Fixed bug that only selected medium size - 9/24
* Fixed code to allow for selecting single size items - 9/24
* Supreme API HTML source Tests - Look for st and s on product page - 9/24
* Supreme API HTML source Tests - Look for articles on page - 9/26
* Bugs - 9/29
  * Bug when category is incorrect, empty category but was unable to advance - Fixed the pick sizing algo but not sure how it happened
  * Apparent checkout bug, unicode issue when printing out return string - Issue was trying to checkout when wasn't in cart
  * Bug when no size is specified but it is a sized item - Same as first bug listed
* Replace logging with: https://github.com/rs/zerolog - 9/29
* Figure out how much time jitter adds to retry, we probably want this kept to a minimum - 9/29
* Task Update: - 10/2
  * Add and use task Id in logging
  * Set task status during everything
  * Move checkout make it work off task log instead of with ID int
* Add support to generate multiple binaries - 10/2
* Go over with spell check - 10/2
* Settings support for different checkout sleep speeds via settings file - 10/2
* Add Licensing and Server Authentication - 10/3
* Add map of all categories - 10/3
* Update log-stats to use taskID and command line - 10/4
* Bug fixes - 10/4
  * Fix unlabled logging that should be there
  * Add log to print item information
  * If ATC is false should kill task
  * Add validity check for items and other fields
* Add Mobile API - 10/6
  * Figure out if mobile can also skip captcha - 10/6
* Finish task SupremeMobileCheckout - 10/9

## Objectives
### 9/20/18
* Test - Successful, 3 Liquid Tees but ended in crash

### 9/27/18
* Failed - Massive amount of user error, set up almost all of the item incorrectly

### 10/4/18
* Failed - Unsure but I belive the mobile API dropped first, some got to checkout but were denied
* Test queuing properly - worked
* Run more versions:
  * Find beta testers - Found 4, 2.5 ran
  * Run at parents and locally - Did not set up, used 
  * Test on Google cloud - Did not set up

### 10/11/18
* Add beta testers
* Test mobile API
* Test unified search pool
* Figure out different ATC responses and replace this with grequests


## Libraries and Code Examples

### Libraries
* https://godoc.org/github.com/levigross/grequests
* https://godoc.org/github.com/PuerkitoBio/goquery
* https://godoc.org/github.com/stretchr/testify
* https://godoc.org/github.com/rs/zerolog
* https://github.com/goreleaser/goreleaser

### Code Examples
* http://polyglot.ninja/golang-making-http-requests/
* https://blog.alexellis.io/golang-writing-unit-tests/
* https://help.github.com/articles/removing-sensitive-data-from-a-repository/
* https://gist.github.com/life1347/69b9f60410070b2609ad2d0779d30cbf
* https://gitlab.com/brandonryan/example/blob/master/logrus/async_test.go
* https://stackoverflow.com/questions/48305425/json-key-can-either-be-a-string-or-an-object
* https://smartystreets.com/blog/2015/02/go-testing-part-1-vanillla
* https://medium.com/@gosamv/using-gos-context-library-for-logging-4a8feea26690
* https://upgear.io/blog/simple-golang-retry-function/
* https://kaviraj.me/understanding-condition-variable-in-go/
* https://github.com/golang/go/wiki/WindowsCrossCompiling
* https://github.com/StefanSchroeder/Golang-Regex-Tutorial/blob/master/01-chapter2.markdown
* https://medium.com/@mlowicki/https-proxies-support-in-go-1-10-b956fb501d6b
