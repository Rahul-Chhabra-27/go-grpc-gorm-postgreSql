package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	db , err := gorm.Open(postgres.New(postgres.Config{
		DSN:"user=postgres password=chhabrarahul dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Db Connected.....");
	db.AutoMigrate(&User{});
	return db;
}
