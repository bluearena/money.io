package bot

import (
    "fmt"
    "regexp"
    "time"

    "github.com/gromnsk/money.io/pkg/storage"
    "github.com/jinzhu/gorm"
    "gopkg.in/telegram-bot-api.v4"
)

type Command string

const COMMAND_START Command = "/start"
const COMMAND_BUDGET Command = "/budget"
const COMMAND_MONTHLY_REPORT Command = "/monthly_report"
const COMMAND_WEEKLY_REPORT Command = "/weekly_report"
const COMMAND_DAILY_REPORT Command = "/daily_report"
const COMMAND_STATISTICS Command = "/statistics"

func process(update tgbotapi.Update, record chan storage.Data, db *gorm.DB) (string, error) {
    command, err := defineCommand(update.Message.Text)
    if err != nil {
        return "", err
    }

    switch command {
    case string(COMMAND_START):
        data, err := processStart(update)
        if err != nil {
            return "", err
        }
        data.Type = storage.TYPE_USER
        record <- data
        return fmt.Sprintf("Hi! Here you can store your expense and income, " +
            "receive daily, weekly and monthly reports and ask some at any time"), nil
    case string(COMMAND_BUDGET):
        data, err := processBudget(update)
        if err != nil {
            return "", err
        }
        record <- data
        return fmt.Sprintf("%s money: %.2f for \"%s\" stored", data.Type, data.Amount, data.Title), nil
    case string(COMMAND_MONTHLY_REPORT):
        return processMonthlyReport(update, db)
    case string(COMMAND_WEEKLY_REPORT):
        return processWeeklyReport(update, db)
    case string(COMMAND_DAILY_REPORT):
        return processDailyReport(update, db)
    case string(COMMAND_STATISTICS):
        //return processStatisticsReport(update, db)
    }

    return "", fmt.Errorf("Can't define command or something happened")
}

func defineCommand(text string) (string, error) {
    re, err := regexp.Compile(`^/(\w+)`)
    if err != nil {
        return "", err
    }
    if !re.MatchString(text) {
        return string(COMMAND_BUDGET), nil
    }

    command := re.FindString(text)

    return command, nil
}

func processStart(update tgbotapi.Update) (storage.Data, error) {
    data := storage.Data{
        UserID:    update.Message.From.ID,
        FirstName: update.Message.From.FirstName,
        LastName:  update.Message.From.LastName,
        UserName:  update.Message.From.UserName,
    }

    return data, nil
}

func processBudget(update tgbotapi.Update) (storage.Data, error) {
    parsedData, err := parse(update.Message.Text)
    if err != nil {
        return storage.Data{}, err
    }
    data := storage.Data{
        UserID:    update.Message.From.ID,
        FirstName: update.Message.From.FirstName,
        LastName:  update.Message.From.LastName,
        UserName:  update.Message.From.UserName,
        Amount:    parsedData.Amount,
        Type:      parsedData.Type,
        Title:     parsedData.Title,
        Currency:  parsedData.Currency,
    }

    return data, nil
}

func processMonthlyReport(update tgbotapi.Update, db *gorm.DB) (string, error) {
    year, month, day := time.Now().Date()
    period := time.Date(year, month - 1, day, 0, 0, 0, 0, time.UTC)
    return processPeriodReport(update.Message.From.ID, period, db)
}

func processWeeklyReport(update tgbotapi.Update, db *gorm.DB) (string, error) {
    year, month, day := time.Now().Date()
    period := time.Date(year, month, day - 7, 0, 0, 0, 0, time.UTC)
    return processPeriodReport(update.Message.From.ID, period, db)
}

func processDailyReport(update tgbotapi.Update, db *gorm.DB) (string, error) {
    year, month, day := time.Now().Date()
    period := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
    return processPeriodReport(update.Message.From.ID, period, db)
}

func processPeriodReport(userId int, date time.Time, db *gorm.DB) (string, error) {
    incomes := []storage.Expense{}
    db.Where("user_id = ? AND date >= ?", userId, date).
        Find(&incomes)
    dataTable := ""
    var sum float32
    for _, money := range incomes {
        dataTable += fmt.Sprintf("%s | %10.2f | %s\n", money.Date.Format("2006-01-02"), money.Amount, money.Name)
        sum += money.Amount
    }

    dataTable += fmt.Sprintf("Total: %.2f", sum)

    return dataTable, nil
}
