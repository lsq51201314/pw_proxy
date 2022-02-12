package main

import (
	"fmt"
	"net"
	"pw_proxy/config"
	"pw_proxy/connlist"
	"pw_proxy/plugins/logger"
	"pw_proxy/plugins/snowflake"
	"pw_proxy/plugins/utils"

	"go.uber.org/zap"
)

func main() {
	logger.WriteText(logger.Log_Level_INFO, "正在载入运行配置项。")
	if err := config.Init(); err != nil {
		logger.WriteText(logger.Log_Level_INFO, "加载配置文件失败。")
		return
	}

	logger.WriteText(logger.Log_Level_INFO, "正在初始化日志系统。")
	if err := logger.Init(config.Configs.LogCfg); err != nil {
		logger.WriteError(logger.Log_Level_INFO, "初始化日志失败。", err)
		return
	}
	defer zap.L().Sync()

	logger.WriteText(logger.Log_Level_INFO, "正在初始化雪花算法。")
	if ok := snowflake.Init(config.Configs.Snowflake.StartTime, config.Configs.Snowflake.MachineID); !ok {
		return
	}

	logger.WriteText(logger.Log_Level_INFO, "正在启动代理服务器。")
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Configs.Localcfg.Port))
	if err != nil {
		logger.WriteText(logger.Log_Level_INFO, fmt.Sprintf("代理服务器运行错误：%v", err))
		return
	}
	defer ln.Close()

	logger.WriteText(logger.Log_Level_INFO, "代理服务器正在运行。")
	for {
		conn, lOk := ln.Accept()
		if lOk != nil {
			logger.WriteText(logger.Log_Level_INFO, fmt.Sprintf("监听客户端连接失败：%v", lOk))
		} else {
			toCon, nOk := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Configs.RemoteCfg.Host, config.Configs.RemoteCfg.Port))
			if nOk != nil {
				logger.WriteText(logger.Log_Level_INFO, fmt.Sprintf("创建新的代理失败：%v", nOk))
				continue
			}

			item := new(connlist.ConnInfo)
			id := snowflake.GenID()

			connlist.ConnList.LoadOrStore(id, item)

			go clientProcess(id, conn, toCon)
			go serverProcess(id, toCon, conn)

			zap.L().Info("通道创建", zap.Int64("id", id), zap.String("request", conn.RemoteAddr().String()), zap.String("response", toCon.RemoteAddr().String()))
		}
	}
}

func clientProcess(id int64, r, w net.Conn) {
	defer r.Close()
	defer w.Close()

	var buffer = make([]byte, 8192)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			logger.WriteText(logger.Log_Level_INFO, fmt.Sprintf("接收客户端请求失败：%v", err))
			break
		}

		data := buffer[:n]
		go request(id, data)

		n, err = w.Write(buffer[:n])
		if err != nil {
			logger.WriteText(logger.Log_Level_INFO, fmt.Sprintf("转发客户端请求失败：%v", err))
			break
		}
	}
	connlist.ConnList.Delete(id)
}

func serverProcess(id int64, r, w net.Conn) {
	defer r.Close()
	defer w.Close()

	var buffer = make([]byte, 8192)
	for {
		n, err := r.Read(buffer)
		if err != nil {
			logger.WriteText(logger.Log_Level_INFO, fmt.Sprintf("接收服务端返回失败：%v", err))
			break
		}

		data := buffer[:n]
		go response(id, data)

		n, err = w.Write(buffer[:n])
		if err != nil {
			logger.WriteText(logger.Log_Level_INFO, fmt.Sprintf("转发服务端返回失败：%v", err))
			break
		}
	}
	connlist.ConnList.Delete(id)
}

func request(id int64, data []byte) {
	if conn, ok := connlist.ConnList.Load(id); ok {
		item := conn.(*connlist.ConnInfo)
		if !item.GetLogin() {
			t := int(data[0])
			switch t {
			case 0x02: //第一步：用户登录
				nLen := utils.CutBytes(data, 2, 1)
				name := utils.CutBytes(data, 3, int(nLen[0]))
				pLen := utils.CutBytes(data, 3+int(nLen[0]), 1)
				passwd := utils.CutBytes(data, 4+int(nLen[0]), int(pLen[0]))
				item.SetUser(string(name), passwd)
			}
		} else {
			e := item.GetDec(data)
			t := int(e[0])
			switch t {
			case 0x03: //第三步：用户登录
				item.SetEnc(data[3:19])
			default:
				zap.L().Info("请求数据", zap.Int64("id", id), zap.String("hex", utils.BytesToHex(e)))
			}
		}
	}
}

func response(id int64, data []byte) {
	if conn, ok := connlist.ConnList.Load(id); ok {
		item := conn.(*connlist.ConnInfo)
		if !item.GetLogin() {
			t := int(data[0])
			switch t {
			case 0x03: //第二步：请求密钥
				item.SetDec(data[3:19])
				item.SetLogin(true)
			}
		} else {
			zap.L().Info("返回数据", zap.Int64("id", id), zap.String("hex", utils.BytesToHex(item.GetEnc(data))))
		}
	}
}
