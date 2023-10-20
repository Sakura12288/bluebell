package main

import (
	"bluebellproject/dao/mysql"
	"bluebellproject/dao/redis"
	"bluebellproject/logger"
	"bluebellproject/routes"
	"bluebellproject/setting"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	//1. 加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("setting.Init() failed err : %v", err)
		return
	}
	//2. 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("logger.Init() failed err : %v", err)
		return
	}
	zap.L().Debug("init logger success")
	defer zap.L().Sync()
	//3. 初始化MySQL链接
	if err := mysql.Init(setting.Conf.MySQLConfig); err != nil {
		zap.L().Fatal("mysql.Init() failed err : %v", zap.Error(err))
		return
	}
	zap.L().Debug("init mysql success")
	defer mysql.Close()
	//4. 初始化Redis链接
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		zap.L().Fatal("redis.Init() failed err : %v", zap.Error(err))
		return
	}
	zap.L().Debug("init redis success")
	defer redis.Close()
	//5. 注册路由
	router := routes.Setup()
	//6. 启动服务(优雅关机)
	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%d",
			setting.Conf.AppConfig.Host,
			setting.Conf.AppConfig.Port),
		Handler: router,
	}
	zap.L().Info("监听地址端口为", zap.Int("port", setting.Conf.AppConfig.Port))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen:", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown server ....")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server shutdown:", zap.Error(err))
	}
	zap.L().Info("success exit")

}
