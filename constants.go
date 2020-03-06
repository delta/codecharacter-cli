package main

var COMPILER_IMAGE = "deltanitt/codecharacter-simulator-compiler-2020:latest"
var RUNNER_IMAGE = "deltanitt/codecharacter-simulator-runner-2020:latest"

//var COMPILER_IMAGE = "deltanitt/codecharacter-compiler-2019:latest"
//var RUNNER_IMAGE = "deltanitt/codecharacter-runner-2019:latest"
var key = []byte(`Key`)

//Default game map
var gameMap = []byte(`L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L F F F F L L L L L L L L L L L L L
L L L L L L L L L L L L L F F F F L L L L L L L L L L L L L
L L L L L L L L L L L L L F F F F L L L L L L L L L L L L L
L L L L L L L L L L L L L F F F F L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
L L L L L L L L L L L L L L L L L L L L L L L L L L L L L L
`)
