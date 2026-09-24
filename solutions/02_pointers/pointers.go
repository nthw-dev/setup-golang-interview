// Package pointers — เฉลยข้อ 2: pointer, pointer receiver, value vs reference
package pointers

import "errors"

var (
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Swap สลับค่าที่ a และ b ชี้อยู่
func Swap(a, b *int) {
	*a, *b = *b, *a
}

type Account struct {
	Owner   string
	Balance int
}

// NewAccount คืน pointer — ตัวแปร local จะ "escape" ไป heap เอง (escape analysis)
func NewAccount(owner string, balance int) *Account {
	return &Account{Owner: owner, Balance: balance}
}

// Deposit ต้องใช้ pointer receiver เพราะต้องแก้ค่าใน struct ตัวจริง
// ถ้าเป็น value receiver (a Account) จะแก้แค่ "สำเนา" แล้วทิ้งไป
func (a *Account) Deposit(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.Balance += amount
	return nil
}

func (a *Account) Withdraw(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount > a.Balance {
		return ErrInsufficientFunds
	}
	a.Balance -= amount
	return nil
}

// Transfer โอนเงิน — ถ้าถอนไม่ผ่าน ต้องไม่มีบัญชีไหนเปลี่ยน
func Transfer(from, to *Account, amount int) error {
	if err := from.Withdraw(amount); err != nil {
		return err
	}
	// amount ถูก validate แล้วใน Withdraw จึง Deposit ไม่มีทาง error
	return to.Deposit(amount)
}

// ApplyBonus เพิ่ม bonus ให้ทุกบัญชีใน slice
// `for _, a := range accounts` ได้ a เป็น "สำเนา" → ต้องใช้ index แทน
func ApplyBonus(accounts []Account, bonus int) {
	for i := range accounts {
		accounts[i].Balance += bonus
	}
}
