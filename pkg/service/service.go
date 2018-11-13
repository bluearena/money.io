package service

import (
    "fmt"

    "github.com/gromnsk/money.io/pkg/app"
    "github.com/gromnsk/money.io/pkg/assets"
    "github.com/gromnsk/money.io/pkg/bot"
    "github.com/gromnsk/money.io/pkg/config"
    "github.com/gromnsk/money.io/pkg/handlers/v1"
    pkglog "github.com/gromnsk/money.io/pkg/logger"
    stdlog "github.com/gromnsk/money.io/pkg/logger/standard"
    "github.com/gromnsk/money.io/pkg/storage"
    "github.com/gromnsk/money.io/pkg/system"
    "github.com/gromnsk/money.io/pkg/utils/startup"
    "github.com/gromnsk/money.io/pkg/version"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/takama/bit"
    "github.com/takama/envconfig"
)

func Run(cfg *config.Config) {
    logger := stdlog.New(&pkglog.Config{
        Level: cfg.LogLevel,
        Time:  true,
        UTC:   true,
    })

    logger.Info("Version: ", version.RELEASE)
    logger.Warnf("%s log level is used", pkglog.LevelDebug.String())

    var dbConf config.DBConfig
    err := envconfig.Process("mysql_service", &dbConf)
    if err != nil {
        logger.Fatalf("Couldn't get database config for DB name, user, password: %s", err)
    }
    logger.Infof("Database configuration: %v %v %v %v\n", dbConf.Host, dbConf.Port, dbConf.Database, dbConf.Username)

    // Create router
    router := bit.NewRouter()

    // Create handlers
    handlers := v1.New(logger, cfg)

    // Asset manager
    assetManager := assets.NewManagerGenerated()

    // Create app
    a := app.NewApp(logger, handlers, router)

    // Setup default routes
    a.SetupDefaultRoutes()

    // Setup blueprints routes
    a.SetupBlueprints(assetManager)

    // Setup app routes
    a.SetupAppRoutes()

    // Listen and serve handlers
    go router.Listen(fmt.Sprintf("%s:%d", cfg.LocalHost, cfg.LocalPort))

    record := make(chan storage.Data, 1000)

    db, err := startup.DB(dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.Database, logger, true)
    if err != nil {
        logger.Fatal(err)
        panic(err.Error())
    }
    defer db.Close()
    storage.Configure(db)

    go bot.Start(logger, cfg, record, db)

    go storage.AsyncSave(db, logger, record)

    logger.Infof("Service %s is listening on %s:%d", config.SERVICENAME, cfg.LocalHost, cfg.LocalPort)

    // Wait signals
    signals := system.NewSignals()
    if err := signals.Wait(logger, new(system.Handling)); err != nil {
        logger.Fatal(err)
    }
}
