package main

import (
	"server/gamedata"
	"fmt"
	"server/conf"
	"server/game"
	"server/gate"
	"server/login"
	// "server/mysql"

	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"github.com/name5566/leaf/log"
	"golang.garena.com/dts/gameserver/gcore/math"
)

// 需要配置的地方：
// 环境变量
// db1, err := gorm.Open("mysql", "mike:123456@tcp(localhost:3306)/poker?parseTime=true")
// conf.Server.LogPath = "C:/Go/bin/log/"
// data, err := ioutil.ReadFile("C:/Go/bin/conf/server.json")
// err = rf.Read("C:/Go/bin/gamedata/" + fn)

func main() {
	// mysql.OpenDB()
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	gamedata.LoadTables()
	testData := gamedata.GetDataByID(2)
	log.Debug(testData.Name)
	log.Debug("1111 log")

	_SceneConfigInstance := new(game.SceneConfig).Create()
	data := new(game.SceneData)
	data.Create("PVP_Paradise",1,"")
	data.Load(_SceneConfigInstance)
	var util = new(game.AutoPathUtil)
	util.Create(data.GetAutoPath())
	// 1268.772,-9.057119,-559.3461
	rs, path := util.FindAutoPath(&math.Vector3{1268.772,-9.057119,-559.3461,},&math.Vector3{1267.772,-9.057119,-559.3461,},)
	fmt.Println("rs ", rs, path)

	fmt.Println("start to run server. ")
	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}

func InitDBTable() {

}
