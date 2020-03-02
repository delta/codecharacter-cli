package main

var COMPILER_IMAGE = "deltanitt/codecharacter-compiler-2019:latest"
var RUNNER_IMAGE = "deltanitt/codecharacter-runner-2019:latest"
var key = []byte(`Key`)

//Default game map
var gameMap = []byte(`L L L L L L L L L L L L L L L L L L L L L L L L L L L L L W \n
W W W W W W W W W W W W W W W W W W W W W W W W W L L L W W \n
L W L L L L W W L L L W W L G G L W W L L L W W L L L W W L \n
L W L W L L L W W L L L W W L L W W L L L W W L L L W W W L \n
L W L W W L L L W W L L L L L W W L L L W W L L L W W L W L \n
L W L L W W L L L W W L L L W W L L L W W L L L W W L L W L \n
L W L L L W W L L L W W L L W L L L W W L L L W W L L L W L \n
L W W L L L W W L L L W W L L L L W W L L L W W L L L L W L \n
L W W W L L L W W L L L W W L L W W L L L W W L L L W L W L \n
L W L W W L L L W W L L L L L W W L L L W W L L L W W L W L \n
L W L L W W L L L W W L L L W W L L L W W L L L W W L L W L \n
L W L L L W W L L L W W L L L L L L W W L L L W W L L L W L \n
L W L L L L W W L L L W W L L L L W W L L L W W L L L W W L \n
L W L W L L L W W L L L W W L L W W L L L W W L L L W W W L \n
L W G W W L L L W W L L L W L L W L L L W W L L L W W G W L \n
L W G W W L L L W W L L L W L L W L L L W W L L L W W G W L \n
L W W W L L L W W L L L W W L L W W L L L W W L L L W L W L \n
L W W L L L W W L L L W W L L L L W W L L L W W L L L L W L \n
L W L L L W W L L L W W L L L L L L W W L L L W W L L L W L \n
L W L L W W L L L W W L L L W W L L L W W L L L W W L L W L \n
L W L W W L L L W W L L L W W L L L L L W W L L L W W L W L \n
L W L W L L L W W L L L W W L L W W L L L W W L L L W W W L \n
L W L L L L W W L L L W W L L L L W W L L L W W L L L W W L \n
L W L L L W W L L L W W L L L W L L W W L L L W W L L L W L \n
L W L L W W L L L W W L L L W W L L L W W L L L W W L L W L \n
L W L W W L L L W W L L L W W L L L L L W W L L L W W L W L \n
L W W W L L L W W L L L W W L L W W L L L W W L L L W L W L \n
L W W L L L W W L L L W W L G G L W W L L L W W L L L L W L \n
W W L L L W W W W W W W W W W W W W W W W W W W W W W W W W \n
W L L L L L L L L L L L L L L L L L L L L L L L L L L L L L \n
`)
