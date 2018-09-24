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
* Add better retrying
* Move checkout into Tasks and make it a class function
* Add proxy support
  * Figure out best practice for maybe options
* Utilize default ro in api better
* Unify / pool inital item search

### Pipeline
* Restock monitor
* Add mobile API
  * Figure out if mobile can also skip captcha
  * Model an interface for current and former system
* Add gcapture code in case they revert

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
* Fix weird task struct formatting
  * Update go code - 9/24
  * Update python code - 9/24

## Objectives

### 9/20/18
* Test - SUCESSFUL, 3 Liquid Tees but ended in crash

### 9/20/18
* Test queueing properly
* Run more versions:
  * Find beta testers
  * Run at parents and locally
  * Test on google cloud

## Libraries and Code Examples

### Libraries
* https://godoc.org/github.com/levigross/grequests
* https://godoc.org/github.com/PuerkitoBio/goquery

### Code Examples
* http://polyglot.ninja/golang-making-http-requests/
* https://blog.alexellis.io/golang-writing-unit-tests/
* https://help.github.com/articles/removing-sensitive-data-from-a-repository/
* https://gist.github.com/life1347/69b9f60410070b2609ad2d0779d30cbf
* https://gitlab.com/brandonryan/example/blob/master/logrus/async_test.go
* https://stackoverflow.com/questions/48305425/json-key-can-either-be-a-string-or-an-object
