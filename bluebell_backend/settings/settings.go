package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var WebConf = new(WebConfig)

type WebConfig struct {
	Name       string `mapstructure:"name"`
	Mode       string `mapstructure:"mode"`
	Version    string `mapstructure:"version"`
	Port       int    `mapstructure:"port"`
	*LogConfig `mapstructure:"log"`
	*DBs       `mapstructure:"dbs"`
	Caches     []*RedisConfig `mapstructure:"caches"`
	//*LogConfig `mapstructure:"log"`
	//*DBConfig `mapstructure:"db"`
	//*RedisConfig `mapstructure:"redis"`

}

// moduleConfig could be in a module specific package

type DBs struct {
	LoginDB DBConfig `mapstructure:"login_db"`
	ShopDB  DBConfig `mapstructure:"shop_db"`
}

type DBConfig struct {
	Type         string `mapstructure:"type"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

func Init() error {
	//viper.SetDefault("ContentDir", "content")
	viper.SetConfigFile("./conf/config.yaml")
	viper.AddConfigPath(".") // 设置配置文件和可执行二进制文件在同一个目录
	//viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	//viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	//viper.AddConfigPath(".") // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			//fmt.Errorf("Fatal error config file: %w \n", err)
			panic("config file not found; ignore error if desired")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("config file was found but another error was produced:%w", err))
		}
	}

	// unmarshal config
	viper.Unmarshal(&WebConf)

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)

		viper.Unmarshal(&WebConf)
	})
	viper.WatchConfig()

	return nil
}
