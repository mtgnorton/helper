package single_flight

import (
	"sync"
)

type (
	WriteFunc func() (interface{}, error)

	ReadFunc func(interface{}, error) (interface{}, error)

	RWSingleFlight interface {
		Do(key string, wf WriteFunc, rf ReadFunc) (interface{}, error)
	}
	call struct {
		wg      sync.WaitGroup
		rf      ReadFunc
		wfValue interface{}
		wfErr   error
	}

	flightGroup struct {
		calls map[string]*call
		lock  sync.Mutex
	}
)

// NewRWSingleFlight returns a RWSingleFlight.
func NewRWSingleFlight() RWSingleFlight {
	return &flightGroup{
		calls: make(map[string]*call),
		lock:  sync.Mutex{},
	}
}
func (g *flightGroup) Do(key string, wf WriteFunc, rf ReadFunc) (interface{}, error) {
	c, done := g.createCall(key)
	if done {
		return c.rf(c.wfValue, c.wfErr)

	}
	g.makeCall(c, key, wf, rf)
	return c.rf(c.wfValue, c.wfErr)
}

func (g *flightGroup) createCall(key string) (c *call, done bool) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c, true
	}
	c = new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()
	return c, false
}

func (g *flightGroup) makeCall(c *call, key string, wf WriteFunc, rf ReadFunc) {
	defer func() {
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()
		c.wg.Done()
	}()

	c.wfValue, c.wfErr = wf()
	c.rf = rf

}
