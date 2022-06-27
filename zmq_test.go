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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestZmq(t *testing.T) {
	endpoint := "tcp://127.0.0.1:5555"

	pushSoc, err := CreatePushSocket(endpoint)
	assert.NotNil(t, pushSoc)
	assert.Nil(t, err)

	pullSoc, err := CreatePullSocket(endpoint)
	assert.NotNil(t, pullSoc)
	assert.Nil(t, err)

	receivedMessages := make(chan string, 20)
	go Pull(pullSoc, receivedMessages)

	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("msg_%v", i)
		err := Push(pushSoc, message)
		assert.Nil(t, err)
	}

	time.Sleep(1 * time.Second)

	assert.Equal(t, 10, len(receivedMessages))
}