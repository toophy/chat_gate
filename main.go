package main

import (
	. "github.com/toophy/chat_gate/master_thread"
	"github.com/toophy/toogo"
)

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 1000, 100, 10000)
	toogo.Run(main_thread)
}
