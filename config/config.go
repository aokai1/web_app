package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var App = new(AppConfig)
var app1 *AppConfig

type AppConfig struct {
	Name  string `mapstructure:"name"`
	Model string `mapstructure:"model"`
	Host  string `mapstructure:"Host"`
	Port  int    `mapstructure:"port"`
	Log   LogConfig
	Mysql MysqlConfig
	Redis RedisConfig
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init() (err error) {

	viper.SetConfigFile("config.yaml") //设置配置文件
	//viper.SetConfigType("yaml")   //设置远程文件类型
	//viper.SetConfigName("config") //设置文件名称
	viper.AddConfigPath(".")   //设置工作目录为配置目录
	err = viper.ReadInConfig() //加入配置
	if err != nil {
		return err
	}
	//把读取到的配置信息反序列化到App 变量中
	if err = viper.Unmarshal(App); err != nil {
		fmt.Printf("viper.Unmarshal failed,err :%v\n", err)
		return
	}
	viper.WatchConfig() //自动检测配置项是否被更改
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置信息被修改了")
		if err = viper.Unmarshal(App); err != nil {
			fmt.Printf("viper.Unmarshal failed,err :%v\n", err)
			return
		}
	})
	return
}
