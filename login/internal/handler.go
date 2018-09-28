package internal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"server/game"
	"server/msg"

	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handleMsg(&msg.LoginReq_0X0101{}, handleAuth)
}

func handleAuth(args []interface{}) {
	m := args[0].(*msg.LoginReq_0X0101)
	a := args[1].(gate.Agent)
	fmt.Printf("客户端%s请求登陆\n", m.AccountName)
	log.Debug("客户端%s请求登陆\n", m.AccountName)
	// fmt.Printf("服务器返回%s的请求登陆结果\n", m.AccountName)
	// a.WriteMsg(&msg.LoginRes_0X0101{Result: true})

	if len(m.AccountName) < 2 || len(m.AccountName) > 12 {
		a.WriteMsg(&msg.LoginFaild{Code: msg.LoginFaild_AccIDInvalid})
		return
	}

	// Account是账户名
	// 想查找是否存在该账户名
	account := getAccountByAccountID(m.AccountName)

	// 将密码存储在数据库中
	data := []byte(m.Password)
	var hash = md5.Sum(data)
	password := hex.EncodeToString(hash[:])

	if nil == account {
		//not having this account,creat account
		newAccount := creatAccountByAccountIDAndPassword(m.AccountName, password)
		if nil != newAccount {
			// game.ChanRPC.Go("CreatePlayer", a, newAccount.ID)
			game.ChanRPC.Go("UserLogin", a, newAccount.ID)
		} else {
			log.Debug("create account error ", m.AccountName)
			a.WriteMsg(&msg.LoginFaild{Code: msg.LoginFaild_InnerError})
		}
	} else {
		// match password
		if password == account.Password {
			game.ChanRPC.Go("UserLogin", a, account.ID)
		} else {
			a.WriteMsg(&msg.LoginFaild{Code: msg.LoginFaild_AccountOrPasswardNotMatch})
		}
	}

}
