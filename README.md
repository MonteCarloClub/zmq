# zmq
ZeroMQ out of the box

## Pull Repository, Install Dependencies and Test
```bash
git clone https://github.com/MonteCarloClub/zmq.git && cd zmq
apt install pkg-config libzmq3-dev
go test -run TestPubAndSub github.com/MonteCarloClub/zmq
go test -run TestPullAndPush github.com/MonteCarloClub/zmq
```

## Import This Go Module
```go
import "github.com/MonteCarloClub/zmq"
```
