
package main

import (
    "fmt"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Rich struct {
    Id      uint `gorm:"primaryKey;default:auto_random()"`
    Respect string
    Spend   int
}

func main() {

    db, err := gorm.Open(mysql.Open("root:jvt123@tcp(127.0.0.1:3306)/cars"), &gorm.Config{})

    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&Rich{}) // create one table

    Richies := Rich{Respect: "no value", Spend: 10000}

    db.Create(&Richies) // insert data

    reder := Rich{}

    db.First(&reder, "Respect=?", "no value") // read data

    fmt.Println(reder)

}