package single_flight

import (
	"fmt"
	"sync"
)

type (
	WriteFunc func() error

	ReadFunc func() (interface{}, error)

	RWSingleFlight interface {
		Do(key string, wf WriteFunc, rf ReadFunc) (interface{}, error)
	}
	call struct {
		wg    sync.WaitGroup
		rf    ReadFunc
		wfErr error
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
		fmt.Println("bb")
		if c.wfErr != nil {
			return nil, c.wfErr
		} else {
			return c.rf()
		}
	}
	g.makeCall(c, key, wf, rf)
	if c.wfErr != nil {
		return nil, c.wfErr
	} else {
		return c.rf()
	}
}

func (g *flightGroup) createCall(key string) (c *call, done bool) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		fmt.Println("aaa")
		g.lock.Unlock()
		c.wg.Wait()
		return c, true
	}
	fmt.Println(key)
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

	c.wfErr = wf()
	if c.wfErr != nil {
		return
	}
	c.rf = rf

}
