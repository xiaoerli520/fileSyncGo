package boot

import (
	"filesync/config"
	"bytes"
	"runtime"
	conf "github.com/spf13/viper"
	"sinago/logger"
	"sinago/utils"
	"strings"
	"sinago/smartPool"
	"strconv"
	"filesync/notification"
)

var (
	EmailList []string
	EmailObject string
	Logger *logger.Logger
	Env string
	Viper = conf.New()
	MailPool *smartPool.SmartPool
)

func InitConfig() *conf.Viper {

	var logLevel int
	confContent := utils.LoadFile(config.CONFIG_PATH)
	Viper.SetConfigType("yaml")
	Viper.ReadConfig(bytes.NewBuffer([]byte(confContent)))

	mailStrategiesTemp := Viper.GetStringMapString("mailStrategies")
	mailStrategies := make(map[int]int)
	for k,v := range mailStrategiesTemp {
		second,_ := strconv.Atoi(k)
		times,_  := strconv.Atoi(v)
		mailStrategies[second] = times
	}
	MailPool = smartPool.NewSmartPool(mailStrategies)
	MailPool.SetSatisfy(func(contents []string) {
		notification.Email(EmailObject, "fileSyncGo_Warning", EmailList, strings.Join(contents,"<br/>"))
	})
	MailPool.Start()

	EmailList  = Viper.GetStringSlice("NoticeEmail")
	EmailObject = Viper.GetString("AppName")

	Env = utils.LoadFile(Viper.GetString("envPath"))
	EnvArr := strings.Split(Env, "=")
	Env = EnvArr[1]
	if strings.Contains(Env,"prod") {
		logLevel = logger.INFO
	} else {
		logLevel = logger.INFO
	}
	Logger = logger.SetLogger(Viper.GetString("LogPath"), logLevel)

	runtime.GOMAXPROCS(runtime.NumCPU())
	return Viper
}

