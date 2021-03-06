// +build datasources

package sources

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/OWASP/Amass/amass/core"
)

func TestCrtsh(t *testing.T) {
	config := &core.Config{}
	config.AddDomain("letsencrypt.owasp-amass.com")
	buf := new(strings.Builder)
	config.Log = log.New(buf, "", log.Lmicroseconds)

	out := make(chan *core.Request)
	bus := core.NewEventBus()
	bus.Subscribe(core.NewNameTopic, func(req *core.Request) {
		out <- req
	})
	defer bus.Stop()

	srv := NewCrtsh(config, bus)
	srv.Start()
	defer srv.Stop()

	expected := 100
	results := make(map[string]int)
	done := time.After(time.Second * 10)

loop:
	for {
		select {
		case req := <-out:
			results[req.Name]++
		case <-done:
			break loop
		}
	}

	if expected != len(results) {
		t.Errorf("Found %d names, expected %d instead", len(results), expected)
	}
}
