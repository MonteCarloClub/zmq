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
	"github.com/MonteCarloClub/log"
	"github.com/MonteCarloClub/utils"
	"github.com/pebbe/zmq4"
)

func (socketSet *SocketSet) SetSubSocket(endpoint string, filter string) error {
	if socketSet == nil {
		log.Error("invalid socket set", utils.NilPtrDeref, utils.NilPtrDerefErr)
		return utils.NilPtrDerefErr
	}

	ctx, err := zmq4.NewContext()
	if err != nil {
		log.Error("fail to create zmq sub context", "err", err)
		return err
	}

	soc, err := ctx.NewSocket(zmq4.SUB)
	if err != nil {
		log.Error("fail to create zmq sub socket", "err", err)
		return err
	}

	err = soc.Connect(endpoint)
	if err != nil {
		log.Error("fail to create outgoing connection from zmq sub socket", "endpoint", endpoint, "err", err)
		return err
	}

	err = soc.SetSubscribe(filter)
	if err != nil {
		log.Error("fail to establish message filter", "filter", filter, "err", err)
		return err
	}

	socketSet.Zmq4SubSocket = soc
	return nil
}

func (socketSet *SocketSet) Sub() (string, error) {
	if socketSet == nil || socketSet.Zmq4SubSocket == nil {
		log.Error("invalid sub socket", utils.NilPtrDeref, utils.NilPtrDerefErr)
		return "", utils.NilPtrDerefErr
	}

	message, err := socketSet.Zmq4SubSocket.Recv(0)
	if err != nil {
		log.Error("fail to receive message from zmq sub socket", "err", err)
		return "", err
	}
	return message, nil
}

func (socketSet *SocketSet) SubToChan(messages chan string, errs chan error) error {
	if socketSet == nil || socketSet.Zmq4SubSocket == nil {
		log.Error("invalid sub socket", utils.NilPtrDeref, utils.NilPtrDerefErr)
		return utils.NilPtrDerefErr
	}

	go func(messages chan string, errs chan error) {
		for {
			message, err := socketSet.Zmq4SubSocket.Recv(0)
			if err != nil {
				log.Warn("fail to receive message from zmq sub socket", "err", err)
				errs <- err
				continue
			}
			messages <- message
		}
	}(messages, errs)
	return nil
}
