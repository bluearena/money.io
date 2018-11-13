package startup

import (
    "fmt"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/gromnsk/money.io/pkg/logger"
)

// DB makes connection with DB, initializes gorm DB level.
func DB(host, port, user, password, name string, log logger.Logger, debug bool) (*gorm.DB, error) {
    dataSource := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", user, password, host, port, name,
    )

    db, err := gorm.Open("mysql", dataSource)
    if err != nil {
        return nil, err
    }

    fmt.Printf("DB: %+v\n", db)
    db.SingularTable(true)
    //db.SetLogger(log)
    if debug {
        db.LogMode(debug)
    }

    return db, nil
}
