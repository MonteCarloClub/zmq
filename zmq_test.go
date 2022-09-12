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
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	messagePrefix         = "msg_"
	messageCount          = 10
	pubEndpoint           = "tcp://*:5555"
	subEndpointExpression = "tcp://127.0.0.*:5555"
	pullAndPushEndpoint   = "tcp://127.0.0.1:5555"
)

// require: index >= 0
func getMessage(index int) string {
	if index < 0 {
		return ""
	}
	return messagePrefix + strconv.Itoa(index)
}

func getIndex(message string) int {
	if len(message) <= len(messagePrefix) {
		return -1
	}
	index, err := strconv.Atoi(message[len(messagePrefix):])
	if err != nil || index < 0 {
		return -1
	}
	return index
}

func TestPubAndSub(t *testing.T) {
	pubSocketSet := CreateSocketSet()
	err := pubSocketSet.SetPubSocket(pubEndpoint)
	assert.NotNil(t, pubSocketSet.Zmq4PubSocket)
	assert.Nil(t, err)

	subSocketSets := make([]*SocketSet, 0)
	for i := 0; i < messageCount; i++ {
		subSocketSet := CreateSocketSet()
		// "tcp://127.0.0.1:5555", "tcp://127.0.0.2:5555",..., "tcp://127.0.0.10:5555"
		err = subSocketSet.SetSubSocket(strings.ReplaceAll(subEndpointExpression, "*", strconv.Itoa(i+1)),
			messagePrefix)
		assert.NotNil(t, subSocketSet.Zmq4SubSocket)
		assert.Nil(t, err)
		subSocketSets = append(subSocketSets, subSocketSet)
	}

	receivedMessages := make(chan string, messageCount*2)
	subErrs := make(chan error, messageCount*2)
	for _, subSocketSet := range subSocketSets {
		// test Sub
		//go func(subSocketSet *SocketSet) {
		//	for {
		//		message, err := subSocketSet.Sub()
		//		if err == nil {
		//			receivedMessages <- message
		//		} else {
		//			subErrs <- err
		//		}
		//	}
		//}(subSocketSet)
		// test SubToChan
		err = subSocketSet.SubToChan(receivedMessages, subErrs)
		assert.Nil(t, err)
	}

	pubErrs := make(chan error, messageCount*2)
	go func(pubSocketSet *SocketSet) {
		for i := 0; i < messageCount; i++ {
			err := pubSocketSet.Pub(getMessage(i))
			if err != nil {
				pubErrs <- err
			}
		}
	}(pubSocketSet)

	time.Sleep(1 * time.Second)

	require.Equal(t, messageCount, len(receivedMessages))
	receivedIndexes := make(map[int]struct{})
	for i := 0; i < messageCount; i++ {
		receivedMessage := <-receivedMessages
		index := getIndex(receivedMessage)
		require.LessOrEqual(t, 0, index)
		require.Greater(t, messageCount, index)
		if _, ok := receivedIndexes[index]; !ok {
			receivedIndexes[index] = struct{}{}
		}
	}
	assert.Equal(t, messageCount, len(receivedIndexes))
	assert.Equal(t, 0, len(pubErrs))
	assert.Equal(t, 0, len(subErrs))
}

func TestPullAndPush(t *testing.T) {
	pullSocketSet, pushSocketSet := CreateSocketSet(), CreateSocketSet()

	err := pushSocketSet.SetPushSocket(pullAndPushEndpoint)
	assert.NotNil(t, pushSocketSet.Zmq4PushSocket)
	assert.Nil(t, err)

	err = pullSocketSet.SetPullSocket(pullAndPushEndpoint)
	assert.NotNil(t, pullSocketSet.Zmq4PullSocket)
	assert.Nil(t, err)

	receivedMessages := make(chan string, messageCount*2)
	pullErrs := make(chan error, messageCount*2)
	// test Pull
	//go func(pullSocketSet *SocketSet) {
	//	for {
	//		message, err := pullSocketSet.Pull()
	//		if err == nil {
	//			receivedMessages <- message
	//		} else {
	//			pullErrs <- err
	//		}
	//	}
	//}(pullSocketSet)
	// test PullToChan
	err = pullSocketSet.PullToChan(receivedMessages, pullErrs)
	assert.Nil(t, err)

	pushErrs := make(chan error, messageCount*2)
	go func(pushSocketSet *SocketSet) {
		for i := 0; i < messageCount; i++ {
			err := pushSocketSet.Push(getMessage(i))
			if err != nil {
				pushErrs <- err
			}
		}
	}(pushSocketSet)

	time.Sleep(1 * time.Second)

	require.Equal(t, messageCount, len(receivedMessages))
	receivedIndexes := make(map[int]struct{})
	for i := 0; i < messageCount; i++ {
		receivedMessage := <-receivedMessages
		index := getIndex(receivedMessage)
		require.LessOrEqual(t, 0, index)
		require.Greater(t, messageCount, index)
		if _, ok := receivedIndexes[index]; !ok {
			receivedIndexes[index] = struct{}{}
		}
	}
	assert.Equal(t, messageCount, len(receivedIndexes))
	assert.Equal(t, 0, len(pullErrs))
	assert.Equal(t, 0, len(pushErrs))
}
