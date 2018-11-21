package gateframework

import (
	"net"
	"MagicseaServerDemo/network"
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
	"reflect"
)

type NetType byte

const (
	TCP NetType = 0
	WEB_SOCKET NetType = 1
)

type Agent interface {
	WriteMsg(msg []byte)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	SetDead()
	GetNetType() NetType
}

type GFAgent struct {
	conn network.Conn
	gate       *Gate
	agentActor *actor.PID
	userData   interface{}
	dead       bool
	netType NetType
}

func (a *GFAgent) GetNetType() NetType {
	return a.netType
}

func (a *GFAgent) SetDead() {
	a.dead = true
}

func (a *GFAgent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		//log.Info("agent.read msg:", len(data))
		if err != nil {
			log.Println("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Println("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				log.Println("route message error: %v", err)
				break
			}
		}else {
			//todo:not safe protobuf
			//a.agentActor.Tell(&msgs.ReceviceClientMsg{data})


			//if err != nil {
			//	log.Error("ReceviceClientMsg message error: %v", err)
			//	break
			//}
		}

	}
}

func (a *GFAgent) OnClose() {
	//todo:not safe
	if a.agentActor != nil && !a.dead {
		//todo
		//a.agentActor.Tell(&msgs.ClientDisconnect{})
	}

}

func (a *GFAgent) WriteMsg(data []byte) {
	err := a.conn.WriteMsg(data)
	if err != nil {
		log.Println("write message %v error: %v", reflect.TypeOf(data), err)
	}

}

func (a *GFAgent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *GFAgent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}


func (a *GFAgent) Close() {
	a.conn.Close()
}

func (a *GFAgent) Destroy() {
	a.conn.Destroy()

}

func (a *GFAgent) UserData() interface{} {
	return a.userData
}

func (a *GFAgent) SetUserData(data interface{}) {
	a.userData = data
}