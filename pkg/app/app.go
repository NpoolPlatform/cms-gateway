package app

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/xerrors"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	"github.com/NpoolPlatform/go-service-framework/pkg/consul"
	"github.com/NpoolPlatform/go-service-framework/pkg/envconf"
	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	msgsrv "github.com/NpoolPlatform/go-service-framework/pkg/rabbitmq/server"
	"github.com/NpoolPlatform/go-service-framework/pkg/redis"
	"github.com/NpoolPlatform/go-service-framework/pkg/version"

	banner "github.com/common-nighthawk/go-figure"
	cli "github.com/urfave/cli/v2"
)

type app struct {
	app   *cli.App
	Redis *redis.Client
}

var myApp = app{}

func Init(
	serviceName, description, usageText, argsUsage string,
	configPath string,
	authors []*cli.Author,
	commands []*cli.Command,
	deps ...string,
) error {
	banner.NewColorFigure(serviceName, "", "green", true).Print()
	ver, err := version.GetVersion()
	if err != nil {
		panic(xerrors.Errorf("Fail to get version: %v", err))
	}

	app := &cli.App{
		Name:        serviceName,
		Version:     ver,
		Description: description,
		ArgsUsage:   argsUsage,
		Usage:       usageText,
		Flags:       nil,
		Commands:    commands,
	}

	myApp.app = app

	err = envconf.Init()
	if err != nil {
		panic(xerrors.Errorf("Fail to init environment config: %v", err))
	}

	err = consul.Init()
	if err != nil {
		panic(xerrors.Errorf("Fail to create consul client: %v", err))
	}

	serviceName = strings.ReplaceAll(serviceName, " ", "")

	err = config.Init(configPath, serviceName, deps...)
	if err != nil {
		panic(xerrors.Errorf("Fail to create configuration: %v", err))
	}

	logDir := config.GetStringValueWithNameSpace("", config.KeyLogDir)
	err = os.MkdirAll(logDir, 0755) //nolint
	if err != nil {
		panic(xerrors.Errorf("Fail to create log dir %v: %v", logDir, err))
	}

	err = logger.Init(logger.DebugLevel, fmt.Sprintf("%v/%v.log", logDir, serviceName))
	if err != nil {
		panic(xerrors.Errorf("Fail to init logger: %v", err))
	}

	myApp.Redis, err = redis.Init()
	if err != nil {
		return xerrors.Errorf("fail to init redis client: %v", err)
	}
	logger.Sugar().Infof("success to create redis client")

	err = msgsrv.Init()
	if err != nil {
		msgsrv.Deinit()
		return xerrors.Errorf("fail to init rabbitmq server: %v", err)
	}
	logger.Sugar().Infof("success to create rabbitmq server")

	return nil
}

func Run(args []string) error {
	return myApp.app.Run(args)
}

// func MySQLConn() *sql.DB {
// 	return myApp.MySQLConn
// }

func Redis() *redis.Client {
	return myApp.Redis
}
