package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
)

func handleContainerLogs(outStream io.ReadCloser, action string) {
	logs, _ := ioutil.ReadAll(outStream)
	outStream = ioutil.NopCloser(bytes.NewBuffer(logs))
	stdcopy.StdCopy(os.Stdout, os.Stderr, outStream)
	if contains(strings.Split(string(logs), " "), "error") || contains(strings.Split(string(logs), " "), "errors") {
		if action == "COMPILE" {
			log.Fatal("Compile error in player code")
		}
		if action == "EXECUTE" {
			log.Fatal("Runtime error in player code")
		}
	}
}

// Tells whether a contains x.
func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func createContextAndClient() (context.Context, *client.Client) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return ctx, cli
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func writeKeyAndMap(dllsDirectory string, gameMap []byte) {
	err := ioutil.WriteFile(path.Join(dllsDirectory, "key.txt"), key, 0666)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(dllsDirectory, "map.txt"), gameMap, 0666)
	if err != nil {
		panic(err)
	}
}

func setUpServeDirectory(outputDirectory string) {
	//Create folder to store player proto files
	protoDirectory := path.Join(outputDirectory, "proto")

	err := os.MkdirAll(protoDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}

	//Copy proto files to proto directory
	err = copy(path.Join(outputDirectory, "player_1.dlog"), path.Join(protoDirectory, "player_1.dlog"))
	err = copy(path.Join(outputDirectory, "player_2.dlog"), path.Join(protoDirectory, "player_2.dlog"))

	staticAsset, err := Asset("index.html")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(outputDirectory, "index.html"), staticAsset, 0666)
	if err != nil {
		panic(err)
	}

	staticAsset, err = Asset("libpack.js")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(outputDirectory, "libpack.js"), staticAsset, 0666)
	if err != nil {
		panic(err)
	}
}

func setupExecutionDirectories(currentPath string, player1Dlls string, player2Dlls string, gameMap []byte) (string, string) {
	//Create directory to store dlls for execution
	dllsDirectory := path.Join(currentPath, "dlls")
	err := os.MkdirAll(dllsDirectory, os.ModePerm)

	//copy only the required files to dll directory
	err = copy(path.Join(player1Dlls, "libplayer_1_code.so"), path.Join(dllsDirectory, "libplayer_1_code.so"))
	err = copy(path.Join(player2Dlls, "libplayer_2_code.so"), path.Join(dllsDirectory, "libplayer_2_code.so"))
	writeKeyAndMap(dllsDirectory, gameMap)

	if err != nil {
		panic(err)
	}

	//Create directory to store dlls for execution
	outputDirectory := path.Join(currentPath, "output_log")
	err = os.MkdirAll(outputDirectory, os.ModePerm)

	return dllsDirectory, outputDirectory
}

func setupCompileDirectories(currentPath string, playerSrcPath string, playerNumber string) (string, string) {
	//Create directory to store player code for compilation
	srcDirectory := path.Join(currentPath, "player"+playerNumber+"source")

	err := os.MkdirAll(srcDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}
	//Copy player1 code to source directory
	err = copy(playerSrcPath, path.Join(srcDirectory, "player_code.cpp"))
	if err != nil {
		panic(err)
	}

	//Create directory to store player1 dlls
	playerDirectory := path.Join(currentPath, "player1Dlls")
	err = os.MkdirAll(playerDirectory, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return srcDirectory, playerDirectory
}

func internetConnected() (ok bool) {
	_, err := http.Get("http://google.com/")
	if err != nil {
		return false
	}
	return true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
