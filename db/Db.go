package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Db struct {
	db *gorm.DB
}


func CreateDb() *Db{
	connect := viper.Get("mysql.connect")
	 db,err := gorm.Open("mysql", connect)
	if err!= nil {
		fmt.Println(err)
	}
	 return &Db{db:db}
}

func GetConfig(){
	path, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(path+"/../bitcoin/config/")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
	}
}
