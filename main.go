package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {

	start := time.Now()
	//Get map file path from the -m option
	mapArg := flag.String("m", "", "Map")
	port := flag.String("p", "3000", "Port")
	// var port int
	// flag.IntVar(&port, "p", 3000, "Port")
	flag.Parse()
	var err error
	// gameMap, err := Asset("map.txt")
	if err != nil {
		panic(err)
	}
	if *mapArg != "" {
		gameMap, err = ioutil.ReadFile(*mapArg)
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
	compile(ctx, player1SrcDirectory, player1DllsDirectory, cli, &wg)

	//Setup compile directories for player2
	player2SrcDirectory, player2DllsDirectory := setupCompileDirectories(currentPath, player2SrcPath, "2")

	//Compile player2 code
	wg.Add(1)
	compile(ctx, player2SrcDirectory, player2DllsDirectory, cli, &wg)

	//Waiting for compilation to finish
	// wg.Wait()
	dllsDirectory, outputDirectory := setupExecutionDirectories(currentPath, player1DllsDirectory, player2DllsDirectory, gameMap)

	//Execute game
	execute(ctx, dllsDirectory, outputDirectory, cli)

	setUpServeDirectory(outputDirectory)

	elapsed := time.Since(start)

	log.Printf("Total time : %s\n", elapsed)

	//Serve output directory
	serve(outputDirectory, *port)
}

func checkAndPullCompilerImage(ctx context.Context, cli *client.Client) {
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
		log.Printf("Pulling compiler image... (Might take a few minutes for first pull)\n")
	}
	reader, err := cli.ImagePull(ctx, COMPILER_IMAGE, types.ImagePullOptions{})
	io.Copy(os.Stdout, reader)
	if err != nil {
		panic(err)
	}
}

func compile(ctx context.Context, inputDir string, outputDir string, cli *client.Client, wg *sync.WaitGroup) {

	log.Printf(inputDir)
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

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})

	if err != nil {
		panic(err)
	}
	handleContainerLogs(out, "COMPILE")
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
		log.Printf("Pulling runner image... (Might take a few minutes for first pull)\n")
	}
	reader, err := cli.ImagePull(ctx, RUNNER_IMAGE, types.ImagePullOptions{})

	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)

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
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		panic(err)
	}
	handleContainerLogs(out, "EXECUTE")

	elapsed := time.Since(start)
	log.Printf("Execution took %s\n", elapsed)
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
