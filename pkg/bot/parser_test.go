package bot

import (
    "testing"
    "github.com/gromnsk/money.io/pkg/storage"
)

var parsedProvider = []struct {
    in string
    out parsedData
} {
    {"закинул баблишка $1370,23 чуваку", parsedData{Amount:1370.23, Currency:"USD",Type:storage.TYPE_EXPENSE,Title:"закинул баблишка чуваку"}},
    {"12500 ремонт машины", parsedData{Amount:12500, Currency:"RUB",Type:storage.TYPE_EXPENSE,Title:"ремонт машины"}},
    {"+40000 зарплата", parsedData{Amount:40000, Currency:"RUB",Type:storage.TYPE_INCOME,Title:"зарплата"}},
    {"-10000 вернул долг", parsedData{Amount:10000, Currency:"RUB",Type:storage.TYPE_EXPENSE,Title:"вернул долг"}},
    {"отдал на ипотеку 20000 бат", parsedData{Amount:20000, Currency:"THB",Type:storage.TYPE_EXPENSE,Title:"отдал на ипотеку"}},
    {"купил продуктов на $100 в 1000 островов", parsedData{Amount:100, Currency:"USD",Type:storage.TYPE_EXPENSE,Title:"купил продуктов на в 1000 островов"}},
}

func TestParse(t *testing.T) {
    for _, data := range parsedProvider {
        t.Run(data.in, func(t *testing.T) {
            parsed, err := parse(data.in)
            if err != nil {
                t.Fatal(err.Error())
            }
            if parsed.Amount != data.out.Amount {
                t.Errorf("Amount should be %d, received %d", data.out.Amount, parsed.Amount)
            }
            if parsed.Type != data.out.Type {
                t.Errorf("Type should be \"%s\", received \"%s\"", data.out.Type, parsed.Type)
            }
            if parsed.Currency != data.out.Currency {
                t.Errorf("Currency should be \"%s\", received \"%s\"", data.out.Currency, parsed.Currency)
            }
            if parsed.Title != data.out.Title {
                t.Errorf("Title should be \"%s\", received \"%s\"", data.out.Title, parsed.Title)
            }
        })
    }
}

func TestParseMoney(t *testing.T) {
    text := "закинул баблишка $1370,23 чуваку"
    parsed, err := parseMoney(text)
    if err != nil {
        t.Fatal(err.Error())
    }
    if parsed.Amount != 1370.23 {
        t.Errorf("Amount should be %d, received %d", 1370.23, parsed.Amount)
    }
    if parsed.Type != storage.TYPE_EXPENSE {
        t.Errorf("Type should be %s, received %s", storage.TYPE_EXPENSE, parsed.Type)
    }
    if parsed.Currency != "USD" {
        t.Errorf("Currency should be %s, received %s", "$", parsed.Currency)
    }
    if parsed.ReplacedString != "закинул баблишка чуваку" {
        t.Errorf("ReplacesString should be %s, received %s", "закинул баблишка чуваку", parsed.ReplacedString)
    }
}