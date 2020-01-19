package main

import (
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/wire"
	"time"
)

func main() {
	wire.GetLogContainerInstance().WriteLog("this is a test!", constant.LOG_INFO)
	time.Sleep(time.Second * 1)
}
