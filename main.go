package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/config"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/routes"
)

// Go Web开发通用脚手架模板

func main() {

	//1.加载配置
	if err := config.Init(); err != nil {
		fmt.Printf("Config Init Failed ,err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(config.App); err != nil {
		fmt.Printf("Logger Init Failed ,err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	//3.初始化MySQL连接
	if err := mysql.Init(config.App); err != nil {
		fmt.Printf("Mysql Init Failed ,err:%v\n", err)
		return
	}
	defer mysql.Close()
	//4.初始化Redis连接
	if err := redis.Init(config.App); err != nil {
		fmt.Printf("Redis Init Failed ,err:%v\n", err)
		return
	}
	defer redis.Close()
	//5.注册路由
	r := routes.Setup()
	//6.启动服务(优雅关机)
	/*err := r.Run(fmt.Sprintf("%s:%s", viper.GetString("host"), viper.GetString("port")))
	if err != nil {
		fmt.Printf("Routes Init Failed ,err:%v\n", err)
		return
	}*/
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", viper.GetString("app.host"), viper.GetString("app.port")),
		Handler: r,
	}
	go func() {
		//开启一个goroutine启动服务
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error(fmt.Sprintf("listen:%s\n", err))
			return
		}
	}()
	//等待终端信号来优雅关闭服务器,为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info(fmt.Sprintf("Shutdown Server ..."))
	//创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//5秒内优雅关闭服务
	if err := server.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown ", zap.Error(err))
		return
	}
	zap.L().Info("Server Exiting")
	return
}
