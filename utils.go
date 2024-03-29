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
	"github.com/pebbe/zmq4"
)

type SocketSet struct {
	Zmq4PubSocket, Zmq4SubSocket   *zmq4.Socket
	Zmq4PullSocket, Zmq4PushSocket *zmq4.Socket
}

func CreateSocketSet() *SocketSet {
	return &SocketSet{}
}

func (socketSet *SocketSet) Close() {
	if socketSet == nil {
		return
	}

	if socketSet.Zmq4PubSocket != nil {
		_ = socketSet.Zmq4PubSocket.SetLinger(0)
		_ = socketSet.Zmq4PubSocket.Close()
	}

	if socketSet.Zmq4SubSocket != nil {
		_ = socketSet.Zmq4SubSocket.SetLinger(0)
		_ = socketSet.Zmq4SubSocket.Close()
	}

	if socketSet.Zmq4PullSocket != nil {
		_ = socketSet.Zmq4PullSocket.SetLinger(0)
		_ = socketSet.Zmq4PullSocket.Close()
	}

	if socketSet.Zmq4PushSocket != nil {
		_ = socketSet.Zmq4PushSocket.SetLinger(0)
		_ = socketSet.Zmq4PushSocket.Close()
	}
}
