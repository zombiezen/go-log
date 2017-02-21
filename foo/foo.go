package main

import (
	"context"
	"os"

	"zombiezen.com/go/log"
	"zombiezen.com/go/log/logutil"
)

var myLog = log.DefaultLogger()

func main() {
	ctx := context.Background()
	logutil.Log(ctx, myLog, "Hello during package init!")
	initLog()
	logutil.Log(ctx, myLog, "Hello after init!")
}

func initLog() {
	stderrLog := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile, nil)
	log.SetDefaultLogger(stderrLog)
}
