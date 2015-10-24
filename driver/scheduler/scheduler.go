// Schedule tasks to run on available resources assigned by master.
package scheduler

import (
	"sync"

	"github.com/chrislusf/glow/driver/scheduler/market"
	"github.com/chrislusf/glow/resource"
)

type Scheduler struct {
	Leader                    string
	EventChan                 chan interface{}
	Market                    *market.Market
	option                    *SchedulerOption
	datasetShard2Location     map[string]resource.Location
	datasetShard2LocationLock sync.Mutex
	waitForAllInputs          *sync.Cond
}

type SchedulerOption struct {
	DataCenter     string
	Rack           string
	TaskMemoryMB   int
	DriverPort     int
	Module         string
	ExecutableFile string
}

func NewScheduler(leader string, option *SchedulerOption) *Scheduler {
	s := &Scheduler{
		Leader:                leader,
		EventChan:             make(chan interface{}),
		Market:                market.NewMarket(),
		datasetShard2Location: make(map[string]resource.Location),
		option:                option,
	}
	s.Market.SetScoreFunction(s.Score).SetFetchFunction(s.Fetch)
	s.waitForAllInputs = sync.NewCond(&s.datasetShard2LocationLock)
	return s
}
