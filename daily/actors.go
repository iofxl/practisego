package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oklog/run"
)

func main() {

	fmt.Println("\n\nExample Group.Add (Basic):\n")

	var g run.Group

	cancel := make(chan struct{})

	g.Add(

		func() error {
			select {
			case <-time.After(time.Second):
				fmt.Printf("The first actor had its time elapsed\n")
				return nil
			case <-cancel:
				fmt.Printf("The first actor was canceled\n")
				return nil

			}
		},

		func(err error) {
			fmt.Printf("The first actor was interrupted with: %v\n", err)
			close(cancel)
		},
	)

	g.Add(
		func() error {
			fmt.Printf("The second actor is returning immediately\n")
			return errors.New("immediate teardown")
		},
		func(err error) {
			fmt.Println("\nNote that this interrupt function is called, even though the corresponding execute function has already returned.")
			fmt.Printf("The second actor was interrupted with: %v\n", err)
		},
	)

	fmt.Printf("The group was terminated with: %v\n", g.Run())

	fmt.Println("\n\nExample Group.Add (Context):\n")

	var g1 run.Group

	ctx, cancelFunc := context.WithCancel(context.Background())

	{
		ctx, cancelFunc := context.WithCancel(ctx) // note: shadowed

		g1.Add(
			func() error {
				return func(ctx context.Context) error {
					return nil
				}(ctx)
			},

			func(err error) {
				cancelFunc()
			},
		)
	}
	go cancelFunc()
	fmt.Printf("The group was terminated with: %v\n", g1.Run())
}
