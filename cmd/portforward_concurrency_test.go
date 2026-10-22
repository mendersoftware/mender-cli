// Copyright 2026 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package cmd

import (
	"strconv"
	"sync"
	"testing"

	"github.com/mendersoftware/go-lib-micro/ws"
)

func TestConnectionChannelsConcurrentAccess(t *testing.T) {
	channels := newConnectionChannels()
	var wg sync.WaitGroup

	for i := 0; i < 200; i++ {
		id := strconv.Itoa(i)
		wg.Add(3)
		go func() {
			defer wg.Done()
			channels.store(id, make(chan *ws.ProtoMsg, 1))
		}()
		go func() {
			defer wg.Done()
			channels.load(id)
		}()
		go func() {
			defer wg.Done()
			channels.delete(id)
		}()
	}

	wg.Wait()
}
