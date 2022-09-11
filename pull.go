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

func (socketSet *SocketSet) SetPullSocket(endpoint string) error {
	if socketSet == nil {
		log.Error("invalid socket set", utils.NilPtrDeref, utils.NilPtrDerefErr)
		return utils.NilPtrDerefErr
	}

	ctx, err := zmq4.NewContext()
	if err != nil {
		log.Error("fail to create zmq pull context", "err", err)
		return err
	}

	soc, err := ctx.NewSocket(zmq4.PULL)
	if err != nil {
		log.Error("fail to create zmq pull socket", "err", err)
		return err
	}

	err = soc.Connect(endpoint)
	if err != nil {
		log.Error("fail to create outgoing connection from zmq pull socket", "endpoint", endpoint, "err", err)
		return err
	}

	socketSet.Zmq4PullSocket = soc
	return nil
}

func (socketSet *SocketSet) Pull() (string, error) {
	if socketSet == nil || socketSet.Zmq4PullSocket == nil {
		log.Error("invalid pull socket", utils.NilPtrDeref, utils.NilPtrDerefErr)
		return "", utils.NilPtrDerefErr
	}

	message, err := socketSet.Zmq4PullSocket.Recv(0)
	if err != nil {
		log.Error("fail to receive message from zmq subscribe socket", "err", err)
		return "", err
	}
	return message, nil
}
