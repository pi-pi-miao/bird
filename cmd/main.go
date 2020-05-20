package main

import (
	"bird/version"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"time"
)

func main(){
	bird := cli.NewApp()
	bird.Name = "bird"
	bird.Compiled = time.Now()
	bird.Version = version.Version
	bird.Commands = []cli.Command{
		start,
	}
	if err := bird.Run(os.Args);err != nil {
		fmt.Println("bird start failed",err)
	}
}

