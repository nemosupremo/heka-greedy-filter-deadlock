package deadlock

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	. "github.com/channelmeter/heka/pipeline"
	"os"
	"time"
)

type EagerInput struct {
	ok bool
}

type EagerInputConfig struct {
}

func (t *EagerInput) EagerInputConfig() interface{} {
	return new(EagerInputConfig)
}

func (t *EagerInput) Init(config interface{}) error {
	t.ok = true
	return nil
}

// Basically this plugin simulates a large number of events
// from some outside world input
func (t *EagerInput) Run(ir InputRunner, h PluginHelper) error {
	var (
		pack   *PipelinePack
		inChan = ir.InChan()
	)
	hostname, _ := os.Hostname()
	for t.ok {
		pack = <-inChan
		pack.Decoded = true
		pack.Message.SetUuid(uuid.NewRandom())
		pack.Message.SetTimestamp(time.Now().UnixNano())
		pack.Message.SetType("greedyfilter.input")
		pack.Message.SetLogger(ir.Name())
		pack.Message.SetHostname(hostname)
		pack.Message.SetPid(int32(os.Getpid()))
		pack.Message.SetSeverity(int32(6))
		pack.Message.SetPayload("nom")
		fmt.Println("[Eager] Sending Pack")
		ir.Inject(pack)
	}
	return nil
}

func (t *EagerInput) Stop() {
	t.ok = false
}

func init() {
	RegisterPlugin("EagerInput", func() interface{} {
		return new(EagerInput)
	})
}
