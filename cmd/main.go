package cmd

import (
	"fmt"

	zap "go.uber.org/zap"
)

var Log *zap.Logger

func main() {
	fmt.Println("Hello World")
	Log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer Log.Sync()

}
