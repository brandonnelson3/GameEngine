package main

import (
	"fmt"

	"github.com/brandonnelson3/GameEngine/messagebus"
)

func logger(m *messagebus.Message) {
	fmt.Println(m.System + ": " + m.Data.(string))
}

func init() {
	messagebus.RegisterType("log", logger)
}
