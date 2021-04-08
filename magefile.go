//+build mage

package main

import (
	"context"
	"log"

	"github.com/magefile/mage/sh"
)

func BuildProxy(ctx context.Context) error {
	out, err := sh.Output("go", "build", "-o", "bin/proxy", "./proxy")
	log.Printf(out)
	return err
}

func RunProxy(ctx context.Context) error {
	out, err := sh.Output("go", "run", "./proxy")
	log.Printf(out)
	return err
}

func BuildOrigin(ctx context.Context) error {
	out, err := sh.Output("go", "build", "-o", "bin/origin", "./origin")
	log.Printf(out)
	return err
}

func RunOrigin(ctx context.Context) error {
	out, err := sh.Output("go", "run", "./origin")
	log.Printf(out)
	return err
}

func BuildScaler(ctx context.Context) error {
	out, err := sh.Output("go", "build", "-o", "bin/scaler", "./scaler")
	log.Printf(out)
	return err
}

func RunScaler(ctx context.Context) error {
	out, err := sh.Output("go", "run", "./scaler")
	log.Printf(out)
	return err
}
