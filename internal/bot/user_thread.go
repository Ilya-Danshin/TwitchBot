package bot

import (
	"fmt"
	"time"
)

type userThread struct {
	ErrorChan   chan error
	ChannelName string
	Modules     []string
}

func (t *userThread) Run(i int) {
	time.Sleep(time.Second * time.Duration(i))
	t.ErrorChan <- fmt.Errorf("test error %d\n", i)
	return
}
