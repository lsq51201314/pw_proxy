package connlist

import (
	"pw_proxy/plugins/rpenc"
	"pw_proxy/plugins/utils"
	"sync"

	"go.uber.org/zap"
)

var ConnList sync.Map

type ConnInfo struct {
	IsLogin bool
	User    struct {
		Name   string
		Passwd string
	}
	Msg struct {
		Dec *rpenc.RPEnc
		Enc *rpenc.RPEnc
	}
}

func (e *ConnInfo) SetLogin(ok bool) {
	e.IsLogin = ok
	zap.L().Info("设置登录", zap.Bool("isLogin", e.IsLogin))
}

func (e *ConnInfo) GetLogin() bool {
	return e.IsLogin
}

func (e *ConnInfo) SetUser(name string, passwd []byte) {
	e.User.Name = name
	e.User.Passwd = utils.BytesToHex(passwd)
	zap.L().Info("设置用户", zap.String("username", e.User.Name), zap.String("password", e.User.Passwd))
}

func (e *ConnInfo) SetDec(passwd []byte) {
	p := utils.MergeBytes(utils.HexToBytes(e.User.Passwd), passwd)
	d := utils.GetBytesHmac(p, []byte(e.User.Name))
	e.Msg.Dec = new(rpenc.RPEnc)
	e.Msg.Dec.Init(d)
	zap.L().Info("请求解密", zap.String("key", utils.BytesToHex(d)))
}

func (e *ConnInfo) SetEnc(passwd []byte) {
	p := utils.MergeBytes(utils.HexToBytes(e.User.Passwd), passwd)
	d := utils.GetBytesHmac(p, []byte(e.User.Name))
	e.Msg.Enc = new(rpenc.RPEnc)
	e.Msg.Enc.Init(d)
	zap.L().Info("返回解密", zap.String("key", utils.BytesToHex(d)))
}

func (e *ConnInfo) GetDec(data []byte) []byte {
	return e.Msg.Dec.Get(data)
}

func (e *ConnInfo) GetEnc(data []byte) []byte {
	return e.Msg.Enc.Get(data)
}
