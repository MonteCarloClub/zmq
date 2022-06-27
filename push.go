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

	"github.com/KofClubs/log"
	zmq "github.com/pebbe/zmq4"
)

func CreatePushSocket(endpoint string) (*zmq.Socket, error) {
	ctx, err := zmq.NewContext()
	if err != nil {
		log.Error("fail to create zmq push context", "err", err)
		return nil, err
	}

	soc, err := ctx.NewSocket(zmq.PUSH)
	if err != nil {
		log.Error("fail to create zmq push socket", "err", err)
		return nil, err
	}

	err = soc.Bind(endpoint)
	if err != nil {
		log.Error("fail to accept incoming connections on zmq push socket", "endpoint", endpoint, "err", err)
		return nil, err
	}

	return soc, nil
}

func Push(pushSoc *zmq.Socket, message string) error {
	size, err := pushSoc.Send(message, 0)
	if err != nil {
		log.Error("fail to send message on zmq push socket", "message", message, "size", size, "err", err)
		return err
	}
	if size <= 0 {
		err = fmt.Errorf("send size non-positive: %v", size)
		log.Error("fail to send message on zmq push socket", "message", message, "size", size, "err", err)
		return err
	}
	return nil
}
