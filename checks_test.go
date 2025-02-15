// Copyright 2017 by the contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package healthcheck

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const testURL = "jsonplaceholder.typicode.com"

func TestTCPDialCheck(t *testing.T) {
	assert.NoError(t, TCPDialCheck(testURL+":80", 5*time.Second)())
	assert.Error(t, TCPDialCheck(testURL+":25327", 5*time.Second)())
}

func TestHTTPGetCheck(t *testing.T) {
	assert.NoError(t, HTTPGetCheck("https://"+testURL+"/posts/1", 5*time.Second)())
	assert.Error(t, HTTPGetCheck("http://x"+testURL+"/posts/1", 5*time.Second)(), "redirect should fail")
	assert.Error(t, HTTPGetCheck("https://"+testURL+"/nonexistent", 5*time.Second)(), "404 should fail")
}
func TestHTTPGetCheckExtended(t *testing.T) {

	assert.NoError(t, HTTPGetCheckExtended("https://"+testURL+"/posts/1", 5*time.Second, []int{200, 201, 301, 437, 502})())
	assert.Error(t, HTTPGetCheckExtended("http://x"+testURL+"/posts/1", 5*time.Second, []int{200, 201, 301, 437, 502})(), "redirect should fail")
	assert.Error(t, HTTPGetCheckExtended("https://"+testURL+"/nonexistent", 5*time.Second, []int{200, 201, 301, 437, 502})(), "404 should fail")
	assert.Error(t, HTTPGetCheckExtended("https://"+testURL+"/posts/1", 5*time.Second, []int{201, 301, 437, 502})(), "200 not in list should fail")
}

func TestDatabasePingCheck(t *testing.T) {
	assert.Error(t, DatabasePingCheck(nil, 1*time.Second)(), "nil DB should fail")

	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	assert.NoError(t, DatabasePingCheck(db, 1*time.Second)(), "ping should succeed")
}

func TestDNSResolveCheck(t *testing.T) {
	assert.NoError(t, DNSResolveCheck(testURL, 5*time.Second)())
	assert.Error(t, DNSResolveCheck("nonexistent."+testURL, 5*time.Second)())
}

func TestGoroutineCountCheck(t *testing.T) {
	assert.NoError(t, GoroutineCountCheck(1000)())
	assert.Error(t, GoroutineCountCheck(0)())
}

func TestGCMaxPauseCheck(t *testing.T) {
	runtime.GC()
	assert.NoError(t, GCMaxPauseCheck(1*time.Second)())
	assert.Error(t, GCMaxPauseCheck(0)())
}
