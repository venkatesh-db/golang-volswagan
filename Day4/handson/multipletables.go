
package main

import (
    "fmt"
    "log"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type User struct {
    Id     uint `gorm:"primaryKey"`
    Name   string
    Email  string
    Orders []Order `gorm:"foreignkey:UserID"`
}

type Order struct {
    Id       uint `gorm:"primaryKey"`
    ItemName string
    Amount   float64
    UserID   uint // forign key
}

func main() {

    db, err := gorm.Open(mysql.Open("root:jvt123@tcp(127.0.0.1:3306)/cars"), &gorm.Config{})

    if err != nil {
        panic("failed to connect database")
    }

    err = db.AutoMigrate(&User{}, &Order{}) // create one table

    if err != nil {
        panic(err)
    }

    fmt.Println("tables created")

    user := User{
        Name:  "somsa",
        Email: "somosa@pune.com",
        Orders: []Order{
            {ItemName: "roti", Amount: 200},
            {ItemName: "dosa", Amount: 160},
        },
    }

    if err := db.Create(&user).Error; err != nil {
        log.Fatal("insert error")
    }
    var users []User

    db.Preload("Orders").Find(&users)

    for _, u := range users {
        fmt.Println("pring users", u.Name, u.Email)
        for _, o := range u.Orders {
            fmt.Println("order", o.ItemName, o.Amount)
        }
    }

}

