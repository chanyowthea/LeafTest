package internal

import (
	"fmt"
	"reflect"
	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func init() {
	// 向当前模块（game 模块）注册 Hello 消息的消息处理函数 handleHello
	handler(&msg.TosChat{}, handleTosChat)
	handler(&msg.LoginReq_0X0101{}, handleLoginReq)
}

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}
func handleTosChat(args []interface{}) {
	// 收到的 Hello 消息
	m := args[0].(*msg.TosChat)
	// 消息的发送者
	a := args[1].(gate.Agent)

	// 输出收到的消息的内容
	log.Debug("hello %v", m.Content)
	fmt.Println("TosChat ", m.Content)
	var err gate.Agent = nil
	if a != err {
		// fmt.Println(" != nil agent=", a)
	}

	//给发送者回应一个 Hello 消息
	a.WriteMsg(&msg.TocChat{
		Name:    m.Name,
		Content: m.Content + "[Server]",
	})
}

func handleLoginReq(args []interface{}) {
	m := args[0].(*msg.LoginReq_0X0101)
	a := args[1].(gate.Agent)

	log.Debug("LoginReq_0X0101 %v", m.AccountName)
	fmt.Println("LoginReq_0X0101 ", m.AccountName)
	if a == nil {
		fmt.Errorf("agent is empty! ")
		return
	}

	fmt.Println("LoginRes_0X0101")
	a.WriteMsg(&msg.LoginRes_0X0101{
		Result: true,
	})
}

//func handleHello(args []interface{}) {
//	// 收到的 Hello 消息
//	m := args[0].(*msg.Hello)
//	// 消息的发送者
//	a := args[1].(gate.Agent)
//
//	// 输出收到的消息的内容
//	log.Debug("hello %v", m.Name)
//	fmt.Println("hello %v", m.Name)
//	// 给发送者回应一个 Hello 消息
//	a.WriteMsg(&msg.Hello{
//		Name: "XXXXXXXXXXXXXXXXXXX",
//	})
//}
