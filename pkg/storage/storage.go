package storage

import (
    "time"

    "github.com/gromnsk/money.io/pkg/logger"
    "github.com/jinzhu/gorm"
)

type Data struct {
    UserID    int
    UserName  string
    FirstName string
    LastName  string
    Title     string
    Amount    float32
    Type      Type
    Currency  string
}

type Type string

const TYPE_INCOME Type = "income"
const TYPE_EXPENSE Type = "expense"
const TYPE_USER Type = "user"

func Configure(db *gorm.DB) {
    db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&User{}, &Income{}, &Expense{})
}

func AsyncSave(db *gorm.DB, logger logger.Logger, record chan Data) {
    var data Data
    var income Income
    var expense Expense
    var user User
    for {
        data = <-record
        switch data.Type {
        case TYPE_INCOME:
            income = Income{
                UserID:   data.UserID,
                Amount:   data.Amount,
                Currency: data.Currency,
                Name:     data.Title,
                Date:     time.Now(),
            }
            db.Save(&income)
        case TYPE_EXPENSE:
            expense = Expense{
                UserID:   data.UserID,
                Amount:   data.Amount,
                Currency: data.Currency,
                Name:     data.Title,
                Date:     time.Now(),
            }
            db.Save(&expense)
        case TYPE_USER:
            user = User{
                UserID:    data.UserID,
                UserName:  data.UserName,
                FirstName: data.FirstName,
                LastName:  data.LastName,
            }
            db.Save(&user)
        }
    }
}
