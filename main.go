package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func printServices() {

	client := NewClient() 
	services, err := client.GetServices()
	if err != nil {
		panic(err)
	}

	for _, service := range services {
		fmt.Println(service.Name)
	}
}

func main() {

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name: "services",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					printServices()
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
