package storage

import "time"

// User is represent telegram user for app
type User struct {
    ID        uint      `gorm:"primary_key" json:"id"`
    UserID    int       `json:"userID"`
    UserName  string    `json:"userName"`
    FirstName string    `json:"firstName"`
    LastName  string    `json:"lastName"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// Income is represent amount of income for some operations
type Income struct {
    ID        uint      `gorm:"primary_key" json:"id"`
    UserID    int       `json:"userID"`
    Amount    float32   `json:"amount"`
    Currency  string    `json:"currency"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
    Date      time.Time `json:"date"`
}

//Expense is represent amount of expenses for some operations
type Expense struct {
    ID        uint      `gorm:"primary_key" json:"id"`
    UserID    int       `json:"userID"`
    Amount    float32   `json:"amount"`
    Currency  string    `json:"currency"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
    Date      time.Time `json:"date"`
}
