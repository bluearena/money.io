package bot

import (
    "fmt"

    "github.com/gromnsk/money.io/pkg/config"
    "github.com/gromnsk/money.io/pkg/logger"
    "github.com/gromnsk/money.io/pkg/storage"
    "github.com/jinzhu/gorm"
    "gopkg.in/telegram-bot-api.v4"
)

func Start(log logger.Logger, cfg *config.Config, record chan storage.Data, db *gorm.DB) {
    bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
    if err != nil {
        log.Error(err)
    }

    bot.Debug = false

    log.Infof("Authorized on account %s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }
        fmt.Println(update.UpdateID)
        processUpdate(bot, update, record, db)
    }
}

func processUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, record chan storage.Data, db *gorm.DB) {
    response, err := process(update, record, db)
    if err != nil {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
        bot.Send(msg)
        return
    }

    msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
    msg.ReplyToMessageID = update.Message.MessageID
    bot.Send(msg)
}
