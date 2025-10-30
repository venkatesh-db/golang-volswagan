

package main

import (
    "fmt"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

type Product struct {
    ID    uint `gorm:"primaryKey;default:auto_random()"`
    Code  string
    Price uint
}

func main() {
    db, err := gorm.Open(mysql.Open("root:jvt123@tcp(127.0.0.1:3306)/sakila"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    db.AutoMigrate(&Product{})

    insertProduct := &Product{Code: "D42", Price: 100}

    db.Create(insertProduct)
    fmt.Printf("insert ID: %d, Code: %s, Price: %d\n",
        insertProduct.ID, insertProduct.Code, insertProduct.Price)

    readProduct := &Product{}
    db.First(&readProduct, "code = ?", "D42") // find product with code D42

    fmt.Printf("read ID: %d, Code: %s, Price: %d\n",
        readProduct.ID, readProduct.Code, readProduct.Price)
}

