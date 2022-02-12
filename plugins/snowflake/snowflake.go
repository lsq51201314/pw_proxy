package snowflake

import (
	"time"

	"go.uber.org/zap"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(startTime string, machineID int64) bool {
	var st time.Time
	var err error
	if st, err = time.Parse("2006-01-02", startTime); err != nil {
		zap.L().Error("无法初始化SnowFlake时间。", zap.Error(err))
		return false
	}
	sf.Epoch = st.UnixNano() / 1000000
	if node, err = sf.NewNode(machineID); err != nil {
		zap.L().Error("无法初始化SnowFlake。", zap.Error(err))
		return false
	}
	return true
}

func GenID() int64 {
	return node.Generate().Int64()
}
