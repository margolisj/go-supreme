# Go_Supreme

Buys some supreme stuff.

## Getting Started

### Running Tests

Individual Test:

~~~sh

~~~

Individual File:

~~~sh

~~~

Integration Tests:


### Runnings Application

1. Make sure you have a task file and settings file somewhere to use on the command line.
2. Build and run.

~~~~sh
#go build
./supreme ./pathto/tasks.json ./optionalFilePath/Settings.json
~~~~

### Task Setup

The task file is json.

~~~json
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
    "api": "mobile",
    "waitSettings": {
      "refreshWait": 150,
      "atcWait": 586,
      "checkoutWait": 776
    }
  }
]
~~~

### Settings Setup

~~~json
{
  "startTime": "2018-10-10T14:59:30.000Z",
  "refreshWait": 300,
  "atcWait": 800,
  "checkoutWait": 800,
}
~~~

## Building Multiple OS Targets

To build different targets you will need to install goreleaser (it is on homebrew). Then Run:

~~~~sh
goreleaser --snapshot
~~~~

### Build Windows Only

~~~sh
GOOS=windows GOARCH=386 go build -ldflags="-s -w" -gcflags="-trimpath=$GOPATH/src" -asmflags="-trimpath=$GOPATH/src" -o supreme-windows.exe
~~~

## TODO

### Updating / Cleaning

* Fix application settings in task.go

### Current

* Fix checkout for mobile and desktop
* Replace checkout waiting with total time from ATC to make sure we are above X ms
* Look to change grequest code with better http transport?
  * Change http transport to have a bigger pool
* Beta version of unify / pool initial item search
* Bugs - 10/18
  * Make sure checkout retries correctly after decline
* Restock monitor
  * Move API code to library

### Pipeline

* Test go obfuscate on something simple / check on stripping linking
* Look to improve algorithm resiliency
* Skip ATC
  * Testing Current Viability
  * Implementations:
    * Desktop
    * Mobile
* Add store credit option to task
* Add get time from some source to calculate computer time drift
* Optimizations
  * Move any tolower processing to task creation / verification?
  * Retry function should telescope to ~ 200 ms, add a setting, but still start maybe 20 or 50 ms
* Self deleting binary when I want the beta over or some http block
  * Probably need some sort of check, look more into keygen code
* Figure out how to set this up - http://www.akins.org/posts/vscode-go/
  * https://github.com/alecthomas/gometalinter
* Metrics server
* After working version
  * Set policy for keygen for only single copy per key
  * Add proxy support for each task
  * Add any size keyword
* Extra security:
  * Have users with e-mail and password
  * Background thread to periodically validate a single copy is open
* Clean up code and model an interface for mobile and desktop
* Auto update
  * https://github.com/tj/go-update
  * https://github.com/inconshreveable/go-update
* https://sequencediagram.org/ Diagram calls
* UI
  * UI Text
    * https://github.com/gizak/termui
    * https://github.com/jroimartin/gocui
    * https://github.com/gdamore/tcell
  * Full UI Version
    * Electron
      * https://nucleus.sh/docs/sell - For the entire js application
    * Add ability to kill one of these go routines via select statement from either cancel channel
      * <https://chilts.org/2017/06/12/cancelling-multiple-goroutines>
    * Add gcapture code in case they revert
    * https://github.com/murlokswarm/app
    * https://github.com/asticode/go-astilectron
    * https://github.com/golang-ui/nuklear
    * https://github.com/zserge/webview
* Command line to feed in file different commands
  * https://github.com/spf13/cobra

### Completed

* Finish API - 9/19/18
  * Finish queuing, see js version - 9/19/18
* Add worker pool - 9/20/18
* Add task importing - 9/20/18
* Create task generator in json - 9/20/18 (in python)
* Add better waiting functionality - 9/20/18
* Attempt to purge credentials from old commits - 9/23/18
* Clean up testing - 9/23/18
* Add go module system - 9/23/18
* Add better logging - 9/23/18
* Fix weird task struct formatting - 9/24/18
  * Update go code - 9/24/18
  * Update python code - 9/24/18
* Add a verify task function - 9/24/18
* Add better retrying - 9/24/18
* Fixed bug that only selected medium size - 9/24/18
* Fixed code to allow for selecting single size items - 9/24/18
* Supreme API HTML source Tests - Look for st and s on product page - 9/24/18
* Supreme API HTML source Tests - Look for articles on page - 9/26/18
* Bugs - 9/29/18
  * Bug when category is incorrect, empty category but was unable to advance - Fixed the pick sizing algo but not sure how it happened
  * Apparent checkout bug, unicode issue when printing out return string - Issue was trying to checkout when wasn't in cart
  * Bug when no size is specified but it is a sized item - Same as first bug listed
* Replace logging with: https://github.com/rs/zerolog - 9/29/18
* Figure out how much time jitter adds to retry, we probably want this kept to a minimum - 9/29/18
* Task Update: - 10/2/18
  * Add and use task Id in logging
  * Set task status during everything
  * Move checkout make it work off task log instead of with ID int
* Add support to generate multiple binaries - 10/2/18
* Go over with spell check - 10/2/18
* Settings support for different checkout sleep speeds via settings file - 10/2/18
* Add Licensing and Server Authentication - 10/3/18
* Add map of all categories - 10/3/18
* Update log-stats to use taskID and command line - 10/4/18
* Bug fixes - 10/4/18
  * Fix unlabled logging that should be there
  * Add log to print item information
  * If ATC is false should kill task
  * Add validity check for items and other fields
* Add Mobile API - 10/6/18
  * Figure out if mobile can also skip captcha - 10/6/18
* Finish task SupremeMobileCheckout - 10/9/18
* Add scheduling for start - 10/9/18
* Add start time to settings - 10/10/18
* Add different API selection to each task - 10/10/18
* Add unit tests for mobile API - 10/10/18
* Test Amex parsing and make sure that is valid - 10/12/18
* Review security code - 10/12/18
  * Add key versioning
  * Add date added
* Embedded timezone information to get rid of windows timezone bug - 10/12/18
* Bugs from 10/11/12 drop
  * New category couldn't get item on mobile - keyword is "new" not "New" - 10/12/18
  * Add check to make sure new is only with mobile - 10/12/18
  * Add increased checkout retry logic - 10/14/18
* Test if I can add cookie and skip ATC - 10/13/18
* Fix mobile for restocks - 10/14/18
* Add status code to all not okay logs - 10/14/18
* Optimizations - 10/14/18
  * use EqualFold for direct comparisons
  * Use pointers more freely in API
* Setup logs to use actual timestamp - 10/15/18
* Increased security - 10/15/18
  * Remove my computers path from errors
  * https://stackoverflow.com/questions/25062696/what-about-protection-for-golang-source-code
  * -s when building
* Add task specific delays - 10/15/18
* Figure out if new works for mobile 10/17/18
* Bugs from 10/18/18
  * Queue bug if queues more than once, see logs 13, 18 mac 1, 14 windows - 10/21/18
  * Add task name to log output - 10/22/18
  * Add API as a log variable and update log stats file - 10/22/18
* Unify redundant return queue logic and checkout logic - 10/21/18
* Adds skipATCMobile into full version - 10/22/18
* Replaces Panic calls with Fatal calls to stop leaking GOPATH - 10/22/18
* Updates tests, removes dead code, removes previous security and fixed desktop item identification - 3/2/20

## Objectives

### 9/20/18

* Test - Successful, 3 Liquid Tees but ended in crash

### 9/27/18

* Failed - Massive amount of user error, set up almost all of the item incorrectly

### 10/4/18

* Failed - Unsure but I believe the mobile API dropped first, some got to checkout but were denied
* Test queuing properly - worked
* Run more versions:
  * Find beta testers - Found 4, 2.5 ran
  * Run at parents and locally - Did not set up, used 
  * Test on Google cloud - Did not set up

### 10/11/18

* Add beta testers
  * Added 3 more as of now
* Test mobile API
* Figure out different ATC responses and replace this with grequests
* Stretch - Test unified search pool
* Results:
  * Failed - I believe delays were wrong
    * Rohit - Delay? checkout issue
    * Me - Card declines, delay? checkout issue, google
    * butch - wrong category
    * wolf - 108 / 150 taks declined probably because of delays
  * Success
    * Mobile API works?
    * Scheduler works
    * Store credit works

### 10/18/2018

* Results:
  * Working tasks and times
    * 13, mac - desktop, 150 786 902 (queued)
    * 18, mac - desktop, 150 792 998 (queued)
    * 14, windows - desktop, 150 645 976 (checkedout and queued)
    * 1, windows - desktop, 150 606 714 (queued)
  * Bugs:
    * Queue is incorrect, after the first response there isn't a slug and it should continue with the old slug
  * Success:
    * 1 should bag windows, 14
    * Looks like everything picked up and worked well today

## 10/25/2018

* Results
  * 6 Tee Checkouts
    * 28, windows, psx6 - philly, here
    * 40, windows, 5d4k - philly, here
    * 34, windows, wxeg - philly, here
    * 35, jake, qxdz - frances
    * 31, windows, vz8s -philly, here
    * 30, windows, 57y0 - philly, missing
  * Still a good number of declines

## Log Greps

~~~sh
grep "\"success\": false," logs/10-11-AllLogs/*wolf*/*
grep "\"success\": true," logs/10-11-AllLogs/*wolf*/*
grep "declined" logs/10-11-AllLogs/*wolf*/* | wc -l
grep "declined" logs/10-11-AllLogs/*wolf*/* | wc -l
grep "Success?" logs/10-11-AllLogs/10-11-mac/* | wc -l
~~~

## Libraries and Code Examples

### Libraries

* https://godoc.org/github.com/levigross/grequests
* https://godoc.org/github.com/PuerkitoBio/goquery
* https://godoc.org/github.com/stretchr/testify
* https://godoc.org/github.com/rs/zerolog
* https://github.com/goreleaser/goreleaser
* https://godoc.org/4d63.com/tz

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
* https://github.com/keygen-sh/example-go-program
* https://www.digitalocean.com/community/questions/how-to-efficiently-compare-strings-in-go
