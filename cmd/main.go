package main

import (
	zap "go.uber.org/zap"
)

var Log *zap.Logger

func main() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer Log.Sync()
}
