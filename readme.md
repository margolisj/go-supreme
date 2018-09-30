# Go_Supreme
Buys some supreme stuff

## Getting Started
1. Make sure you have a task file somewhere and its being pointed at in code.
2. Build and go. Pretty easy.
~~~~
export GO111MODULE=on
go build
./supreme
~~~~

## TODO:
### Current
* Task Update:
  * Add request object
  * Set task status during everything
  * Move checkout make it work off task
* Add proxy support
* Commandline to feed in file different commands
  * https://github.com/spf13/cobra

### Pipeline
* Add authentication
  * https://github.com/denisbrodbeck/machineid  
* https://sequencediagram.org/ Diagram calls
* Restock monitor
* Utilize default ro in api better
* Unify / pool inital item search
* Add mobile API
  * Figure out if mobile can also skip captcha
  * Model an interface for current and former system
* Add gcapture code in case they revert
* UI
  * make GUI to keep track of status
  * https://github.com/gizak/termui
  * https://github.com/jroimartin/gocui
  * https://github.com/gdamore/tcell

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
* Figure out how much time jitter adds to retrys, we probably want this kept to a minimum - 9/29

## Objectives

### 9/20/18
* Test - SUCESSFUL, 3 Liquid Tees but ended in crash

### 9/27/18
* Failed - Massive amount of user error, set up almost all of the item incorrectly

### 10/6/18
* Test queueing properly
* Run more versions:
  * Find beta testers
  * Run at parents and locally
  * Test on google cloud

## Libraries and Code Examples

### Libraries
* https://godoc.org/github.com/levigross/grequests
* https://godoc.org/github.com/PuerkitoBio/goquery
* https://godoc.org/github.com/stretchr/testify
* https://godoc.org/github.com/rs/zerolog

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
