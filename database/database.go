package database

import (
	"bwastartup/campaign"
	"bwastartup/transaction"
	"bwastartup/user"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func ConnectToDB() *gorm.DB {
// 	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
// 	dsn := "root@tcp(localhost:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
// 	// dsn := "root:a@tcp(0.0.0.0:5555)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	} else {
// 		fmt.Println(db, "connection succes")
// 	}
// 	return db
// }

func ConnectToDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=a dbname=bwa port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	} else {
		fmt.Println(db, "connection succes")
	}
	return db
}

func AutoMigrate(DB *gorm.DB) {
	DB.AutoMigrate(user.User{}, transaction.Transaction{}, campaign.Campaign{}, campaign.CampaignImage{})
}
