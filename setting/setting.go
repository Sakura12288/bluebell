package setting

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

//结构体映射

var Conf = new(AllConfig)

type AllConfig struct {
	*AppConfig   `mapstructure:"app"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Database    string `mapstructure:"database"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
}
type RedisConfig struct {
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	Password    string        `mapstructure:"password"`
	MaxOpenConn int           `mapstructure:"max_open_conn"`
	MaxIdleConn int           `mapstructure:"max_idle_conn"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
	DB          int           `mapstructure:"db"`
}

func Init() (err error) {
	viper.SetConfigFile("config.yaml")
	//viper.SetConfigName("config")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed err : %v", err)
		return
	}
	if err = viper.Unmarshal(&Conf); err != nil {
		fmt.Printf("viper.Unmarshal(&Conf) failed err : %v", err)
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("系统配置已改变")
		if err = viper.Unmarshal(&Conf); err != nil {
			fmt.Printf("viper.Unmarshal(&Conf) again failed err : %v", err)
			return
		}
	})
	return
}
