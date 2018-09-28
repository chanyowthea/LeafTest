package internal

import (
	"fmt"
	"server/mysql"
)

type Account struct {
	ID          uint   `gorm:"primary_key;AUTO_INCREMENT"`
	AccountName string `gorm:"not null;unique"`
	Password    string
}

func getAccountByAccountID(accountName string) *Account {

	var account Account

	// gorm的database
	db := mysql.MysqlDB()
	// db.DropTableIfExists(&Account{})
	if rs := db.HasTable(&Account{}); !rs {
		db.AutoMigrate(&Account{})
	}

	err := db.Where("account_name = ?", accountName).Limit(1).Find(&account).Error
	if nil != err {
		fmt.Println(err)
		return nil
	}
	fmt.Println("password:", account.Password)
	return &account
}

// 创建账号名存储在数据库
func creatAccountByAccountIDAndPassword(accountName string, password string) *Account {
	db := mysql.MysqlDB()
	// db.DropTableIfExists(&Account{})
	if rs := db.HasTable(&Account{}); !rs {
		db.AutoMigrate(&Account{})
	}

	var account = Account{AccountName: accountName, Password: password}
	err := db.Create(&account).Error
	if nil != err {
		return nil
	}

	return &account
}
