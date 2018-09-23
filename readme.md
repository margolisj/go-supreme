# Go_Supreme
Buys some supreme stuff

## TODO:
### Current
* Attempt to purge credentials from old commits?
* Add better logging
* Add better testing
  * Clean up current Tests and separate networked and non-networked tests
* Add better retrying
* Move checkout into Tasks and make it a class function
* Add better waiting functionality
  * https://nathanleclaire.com/blog/2014/02/15/how-to-wait-for-all-goroutines-to-finish-executing-before-continuing/
* Fix weird task idea
  * Update go code
  * Update python code
* Add proxy support
  * Figure out best practice for maybe options

### Pipeline
* Restock monitor
* Add mobile API
  * Figure out if mobile can also skip captcha
  * Model an interface for current and former system
* Add gcapture code in case they revert

### Completed
* Finish API - 8/19
  * Finish queuing, see js version - 8/19
* Add worker pool - 8/20
* Add task importing - 8/20
* Create task generator in json - 8/20 (in python)

## Objectives

### 9/20/18
* Test - SUCESSFUL, 3 Liquid Tees but ended in crash

### 9/20/18
* Test queuing properly
* Run at parents and locally
* Test on google cloud

## Libraries and Code Examples

### Libraries
* https://godoc.org/github.com/levigross/grequests
* https://godoc.org/github.com/PuerkitoBio/goquery

### Code Examples
* http://polyglot.ninja/golang-making-http-requests/
* https://blog.alexellis.io/golang-writing-unit-tests/