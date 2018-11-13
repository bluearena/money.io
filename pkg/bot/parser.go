package bot

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"

    "github.com/gromnsk/money.io/pkg/storage"
)

type parsedData struct {
    Amount   float32
    Currency string
    Type     storage.Type
    Title    string
}

type parsedMoney struct {
    Amount         float32
    Currency       string
    Type           storage.Type
    ReplacedString string
}

func parse(text string) (data parsedData, err error) {
    parsedMoney, err := parseMoney(text)
    if err != nil {
        return data, err
    }
    data = parsedData{
        Amount:   parsedMoney.Amount,
        Currency: parsedMoney.Currency,
        Type:     parsedMoney.Type,
        Title:    parsedMoney.ReplacedString,
    }

    return data, nil
}

func parseMoney(text string) (parsed parsedMoney, err error) {
    re, err := regexp.Compile(`-?(?P<currency>[$])?(?P<type>\+?)(?P<int>\d+)[.,]?(?P<frac>\d{1,2})?\s?(?P<currency>(руб(\w*))|(р\s)|(бат(\w*)))?`)
    if err != nil {
        return parsedMoney{}, err
    }

    if !re.MatchString(text) {
        return parsedMoney{}, fmt.Errorf("Bro, I can't find any money in your fucking shit \"%s\"", text)
    }

    result := re.FindStringSubmatch(text)
    currency := "RUB"
    if result[1] != "" {
        currency = parseCurrency(result[1])
    } else if result[5] != "" {
        currency = parseCurrency(result[5])
    }
    value, err := strconv.ParseFloat(result[3], 32)
    if err != nil {
        return parsedMoney{}, err
    }
    intAmount := float32(value)

    value, _ = strconv.ParseFloat(result[4], 32)
    fracAmount := float32(value) / 100

    amountType := storage.TYPE_EXPENSE
    if result[2] != "" {
        amountType = storage.TYPE_INCOME
    }

    amount := intAmount + fracAmount
    found := re.FindString(text)
    replaced := text
    if found != "" {
        replaced = strings.Replace(text, found, "", 1)
    }

    parsed = parsedMoney{
        Amount:         amount,
        Currency:       currency,
        Type:           amountType,
        ReplacedString: strings.TrimSpace(replaced),
    }

    return parsed, nil
}

func parseCurrency(text string) string {
    if text == "$" {
        return "USD"
    } else if text[:3] == "руб" {
        return "RUB"
    } else if text == "бат" {
        return "THB"
    }

    return "RUB"
}
