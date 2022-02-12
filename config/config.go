package config

var Configs = new(AppCfg)

type AppCfg struct {
	*Localcfg  `mapstructure:"local"`
	*RemoteCfg `mapstructure:"remote"`
	*LogCfg    `mapstructure:"log"`
	*Snowflake `mapstructure:"snowflake"`
}

type Localcfg struct {
	Port uint16 `mapstructure:"port"`
}

type RemoteCfg struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

type LogCfg struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"maxsize"`
	MaxAge     int    `mapstructure:"maxage"`
	MaxBackups int    `mapstructure:"maxbackups"`
}

type Snowflake struct {
	StartTime string `mapstructure:"starttime"`
	MachineID int64  `mapstructure:"machineid"`
}
