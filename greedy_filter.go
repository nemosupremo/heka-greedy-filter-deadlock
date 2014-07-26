package deadlock

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	. "github.com/channelmeter/heka/pipeline"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

type GreedyFilter struct {
	hostname string
	h        PluginHelper
	fr       FilterRunner
	msgLoop  uint
}

type GreedyFilterConfig struct {
}

func (gf *GreedyFilter) ConfigStruct() interface{} {
	return &GreedyFilterConfig{}
}

func (gf *GreedyFilter) Init(config interface{}) (err error) {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	return nil
}

func (gf *GreedyFilter) Flush() (err error) {
	// Ex. we may have 10,000 key/value pairs that we want to flush
	// to some output.
	for i := 0; i < 10000; i++ {
		fmt.Println("[Greedy] Requesting pack")
		pack := gf.h.PipelinePack(gf.msgLoop)
		pack.Message.SetUuid(uuid.NewRandom())
		pack.Message.SetTimestamp(time.Now().UnixNano())
		pack.Message.SetType("greedyfilter.output")
		pack.Message.SetLogger(gf.fr.Name())
		pack.Message.SetHostname(gf.hostname)
		pack.Message.SetPid(int32(os.Getpid()))
		pack.Message.SetSeverity(int32(6))
		pack.Message.SetPayload("nom")
		fmt.Println("[Greedy] injecting pack")
		gf.fr.Inject(pack)
	}
	return nil
}

func (gf *GreedyFilter) Run(fr FilterRunner, h PluginHelper) (err error) {
	var (
		pack   *PipelinePack
		ticker = fr.Ticker()
		inChan = fr.InChan()
	)

	ok := true
	gf.h = h
	gf.fr = fr
	gf.hostname, _ = os.Hostname()
	for ok {
		select {
		case pack, ok = <-inChan:
			if ok {
				fmt.Println("[Greedy] Accepting Message")
				gf.msgLoop = pack.MsgLoopCount
				// Do some work on the incoming pack
				// such as aggregate tham
				pack.Recycle()
			}
		case <-ticker:
			fmt.Println("[Greedy] Flushing Greedy")
			gf.Flush()
		}
	}
	return
}

func init() {
	RegisterPlugin("GreedyFilter", func() interface{} {
		return new(GreedyFilter)
	})
}
