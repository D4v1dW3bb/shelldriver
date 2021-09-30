package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/progrium/qtalk-go/mux"
	"github.com/progrium/qtalk-go/mux/frame"
	"github.com/progrium/shelldriver/bridge"
)

const Version = "0.1.0"

func init() {
	runtime.LockOSThread()
}

func main() {

	flagDebug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	if *flagDebug {
		fmt.Fprintf(os.Stderr, "shellbridge %s\n", Version)
		frame.Debug = os.Stderr
	}

	sess, err := mux.DialIO(os.Stdout, os.Stdin)
	if err != nil {

		log.Fatal(err)
	}

	srv := bridge.NewServer()

	go srv.Respond(sess)
	bridge.Main()
	//select {}
	// for i := 0; i < 5; i++ {
	// 	time.Sleep(1 * time.Millisecond)
	// }

	// sess.Wait()
}
