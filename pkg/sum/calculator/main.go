package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/efficientgo/examples/pkg/sum"
	"github.com/pkg/errors"
)

var (
	sumFlags = flag.NewFlagSet("sum", flag.ExitOnError)
	sumFile  = sumFlags.String("file", "", "File with numbers delimited by a new line")
)

func main() {
	if err := runMain(os.Args[1:]); err != nil {
		// Use %+v for github.com/pkg/errors error to print with stack.
		log.Fatalf("Error: %+v", errors.Wrapf(err, "%s", flag.Arg(0)))
	}
}

func runMain(args []string) (err error) {
	if len(args) < 1 {
		return errors.New("missing a sub-command; available: 'sum'")
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "sum":
		if err := sumFlags.Parse(args); err != nil {
			return err
		}

		if *sumFile == "" {
			return errors.New("missing -file flag")
		}

		s, err := sum.Sum(*sumFile)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	default:
		return errors.Errorf("unknown command: %v", subcommand)
	}
}
