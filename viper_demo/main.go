package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Port        int    `mapstructure:"port"` //注意：viper反序列化时，tag必须用 "mapstructure"， 不管配置文件是json, yaml, ....
	Version     string `mapstructure:"version"`
	MySQLConfig `mapstructure:"mysql"`
}

type MySQLConfig struct {
	Host   string `mapstructure:"host"`
	Dbname string `mapstructure:"dbname"`
	Port   int    `mapstructure:"port"`
}

func main() {
	//设置默认值
	viper.SetDefault("fileDir", "./")

	viper.SetConfigFile("./config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")        // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")          // 如果配置文件的名称中没有扩展名，则需要配置此项

	//依次查找多个配置文件的配置项
	viper.AddConfigPath("/etc/appname/")  // 查找配置文件所在的路径(指定多个配置文件的目录）
	viper.AddConfigPath("$HOME/.appname") // 多次调用以添加多个搜索路径
	viper.AddConfigPath(".")              // 还可以在工作目录中查找配置

	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.WriteConfig()                            // 将当前配置写入“viper.AddConfigPath()”和“viper.SetConfigName”设置的预定义路径(WriteConfig writes the current configuration to a file.)
	viper.SafeWriteConfig()                        //SafeWriteConfig writes current configuration to file only if the file does not exist
	viper.WriteConfigAs("/path/to/my/.config")     //将当前的viper配置写入给定的文件路径。将覆盖给定的文件(如果它存在的话)。
	viper.SafeWriteConfigAs("/path/to/my/.config") // 因为该配置文件写入过，所以会报错
	viper.SafeWriteConfigAs("/path/to/my/.other_config")

	/*
		需要重新启动服务器以使配置生效的日子已经一去不复返了，
		viper驱动的应用程序可以在运行时读取配置文件的更新，而不会错过任何消息。
	*/
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件发生变更之后会调用的回调函数
		fmt.Println("Config file changed:", e.Name)
	})

	fmt.Println(viper.Get("filedir"), viper.Get("mysql.dbname"))

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Println("Unmarshal failed")
	} else {
		fmt.Printf("c:%#v\n", c)
	}
	//r := gin.Default()
	//r.GET("/version", func(c *gin.Context) {
	//	c.String(http.StatusOK, viper.GetString("version"))
	//})
	//r.Run()
}
