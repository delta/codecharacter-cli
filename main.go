package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"golang.org/x/net/context"
)

func main() {

	start := time.Now()
	//Get map file path from the -m option
	mapFile := flag.String("m", "", "Map")
	port := flag.String("p", "3000", "Port")
	// var port int
	// flag.IntVar(&port, "p", 3000, "Port")
	flag.Parse()
	var err error
	if *mapFile != "" {
		gameMap, err = ioutil.ReadFile(*mapFile)
		if err != nil {
			panic(err)
		}
	}

	//Create new client and context
	ctx, cli := createContextAndClient()

	player1SrcPath := flag.Arg(0)
	player2SrcPath := flag.Arg(1)

	//Get current working directory
	currentPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	var wg sync.WaitGroup

	//Setup compile directories for player1
	player1SrcDirectory, player1DllsDirectory := setupCompileDirectories(currentPath, player1SrcPath, "1")

	checkAndPullCompilerImage(ctx, cli)

	//Compile player2 code
	wg.Add(1)
	go compile(ctx, player1SrcDirectory, player1DllsDirectory, cli, &wg)

	//Setup compile directories for player2
	player2SrcDirectory, player2DllsDirectory := setupCompileDirectories(currentPath, player2SrcPath, "2")

	//Compile player2 code
	wg.Add(1)
	go compile(ctx, player2SrcDirectory, player2DllsDirectory, cli, &wg)

	//Waiting for compilation to finish
	wg.Wait()
	dllsDirectory, outputDirectory := setupExecutionDirectories(currentPath, player1DllsDirectory, player2DllsDirectory)

	//Execute game
	execute(ctx, dllsDirectory, outputDirectory, cli)

	setUpServeDirectory(outputDirectory)

	elapsed := time.Since(start)

	log.Printf("Total time : %s\n", elapsed)

	//Serve output directory
	serve(outputDirectory, *port)
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

func checkAndPullCompilerImage(ctx context.Context, cli *client.Client) {
	start := time.Now()
	//List all images
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Printf("Make sure docker is installed on your machine\n")
		panic(err)
	}
	compilerImageExists := false
	for _, image := range images {
		if len(image.RepoTags) > 0 && image.RepoTags[0] == COMPILER_IMAGE {
			compilerImageExists = true
		}
	}

	if !compilerImageExists {
		log.Printf("Pulling compiler image...\n")
		reader, err := cli.ImagePull(ctx, COMPILER_IMAGE, types.ImagePullOptions{})
		io.Copy(os.Stdout, reader)
		if err != nil {
			panic(err)
		}
		elapsed := time.Since(start)
		log.Printf("Compiler image pull took %s\n", elapsed)
	}
	if err != nil {
		panic(err)
	}
}

func compile(ctx context.Context, inputDir string, outputDir string, cli *client.Client, wg *sync.WaitGroup) {

	defer wg.Done()

	log.Printf("Compiling...\n")
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: COMPILER_IMAGE,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: outputDir,
				Target: "/root/output_libs",
			},
			{
				Type:   mount.TypeBind,
				Source: inputDir,
				Target: "/root/codecharacter/src/player_code/src",
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		panic(err)
	}

}

func setupExecutionDirectories(currentPath string, player1Dlls string, player2Dlls string) (string, string) {
	//Create directory to store dlls for execution
	dllsDirectory := path.Join(currentPath, "dlls")
	err := os.MkdirAll(dllsDirectory, os.ModePerm)

	//copy only the required files to dll directory
	err = copy(path.Join(player1Dlls, "libplayer_1_code.so"), path.Join(dllsDirectory, "libplayer_1_code.so"))
	err = copy(path.Join(player2Dlls, "libplayer_2_code.so"), path.Join(dllsDirectory, "libplayer_2_code.so"))
	writeKeyAndMap(dllsDirectory)

	if err != nil {
		panic(err)
	}

	//Create directory to store dlls for execution
	outputDirectory := path.Join(currentPath, "output_log")
	err = os.MkdirAll(outputDirectory, os.ModePerm)

	return dllsDirectory, outputDirectory
}

func execute(ctx context.Context, inputDir string, outputDir string, cli *client.Client) {

	start := time.Now()
	//List all images
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}
	runnerImageExists := false
	for _, image := range images {
		if len(image.RepoTags) > 0 && image.RepoTags[0] == RUNNER_IMAGE {
			runnerImageExists = true
		}
	}

	if !runnerImageExists {
		log.Printf("Pulling runner image...\n")
		reader, err := cli.ImagePull(ctx, RUNNER_IMAGE, types.ImagePullOptions{})
		io.Copy(os.Stdout, reader)
		elapsed := time.Since(start)
		log.Printf("Runner image pull took %s\n", elapsed)
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}

	log.Printf("Executing...\n")
	//Create runner container
	start = time.Now()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: RUNNER_IMAGE,
		Tty:   true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: inputDir,
				Target: "/root/input_libs",
			},
			{
				Type:   mount.TypeBind,
				Source: outputDir,
				Target: "/root/output_log",
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	elapsed := time.Since(start)
	log.Printf("Execution took %s\n", elapsed)
}

func writeKeyAndMap(dllsDirectory string) {
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

func serve(outputDirectory string, port string) {
	http.Handle("/", http.FileServer(http.Dir(outputDirectory)))
	clipboard.WriteAll("http://localhost:" + port)
	log.Printf("Serving at http://localhost:" + port + " (Link copied to clipboard)")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("If port is busy, specify a different one using the -p flag")
		log.Fatal("ListenAndServe: ", err)
	}
}
