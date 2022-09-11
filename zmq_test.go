/*
Copyright (c) 2022 Zhang Zhanpeng <zhangregister@outlook.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package zmq

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	endpoint = "tcp://127.0.0.1:5555"
)

// require: index >= 0
func getMessage(index int) string {
	if index < 0 {
		return ""
	}
	return fmt.Sprintf("msg_%v", index)
}

func getIndex(message string) int {
	msgAndIndex := strings.Split(message, "_")
	if len(msgAndIndex) < 2 {
		return -1
	}
	index, err := strconv.Atoi(msgAndIndex[1])
	if err != nil {
		return -1
	}
	return index
}

func TestPullAndPush(t *testing.T) {
	pullSocketSet, pushSocketSet := &SocketSet{}, &SocketSet{}

	err := pushSocketSet.SetPushSocket(endpoint)
	assert.NotNil(t, pushSocketSet.Zmq4PushSocket)
	assert.Nil(t, err)

	err = pullSocketSet.SetPullSocket(endpoint)
	assert.NotNil(t, pullSocketSet.Zmq4PullSocket)
	assert.Nil(t, err)

	receivedMessages := make(chan string, 20)
	pullErrs := make(chan error, 20)
	go func(pullSocketSet *SocketSet) {
		for {
			message, err := pullSocketSet.Pull()
			if err == nil {
				receivedMessages <- message
			} else {
				pullErrs <- err
			}
		}
	}(pullSocketSet)

	pushErrs := make(chan error, 10)
	go func(pushSocketSet *SocketSet) {
		for i := 0; i < 10; i++ {
			err := pushSocketSet.Push(getMessage(i))
			if err != nil {
				pushErrs <- err
			}
		}
	}(pushSocketSet)

	time.Sleep(1 * time.Second)

	require.Equal(t, 10, len(receivedMessages))
	receivedIndexes := make(map[int]struct{})
	for i := 0; i < 10; i++ {
		receivedMessage := <-receivedMessages
		index := getIndex(receivedMessage)
		require.LessOrEqual(t, 0, index)
		require.Greater(t, 10, index)
		if _, ok := receivedIndexes[index]; !ok {
			receivedIndexes[index] = struct{}{}
		}
	}
	assert.Equal(t, 10, len(receivedIndexes))
	assert.Equal(t, 0, len(pullErrs))
	assert.Equal(t, 0, len(pushErrs))
}
