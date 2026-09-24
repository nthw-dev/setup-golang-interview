// Package pointers — ข้อ 2: pointer, pointer receiver, value vs reference
// ⚠️ โค้ดบางส่วนมีบั๊กจงใจใส่ไว้ — หาและแก้ให้ test ผ่าน
package pointers

import "errors"

var (
	ErrInvalidAmount     = errors.New("invalid amount")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Swap สลับค่าที่ a และ b ชี้อยู่
func Swap(a, b *int) {
	// TODO
}

type Account struct {
	Owner   string
	Balance int
}

func NewAccount(owner string, balance int) *Account {
	return &Account{Owner: owner, Balance: balance}
}

// Deposit ฝากเงิน (amount <= 0 → ErrInvalidAmount)
// BUG: ฝากแล้วยอดเงินไม่เปลี่ยน ทำไม?
func (a Account) Deposit(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.Balance += amount
	return nil
}

// Withdraw ถอนเงิน
// amount <= 0 → ErrInvalidAmount, เงินไม่พอ → ErrInsufficientFunds
func (a *Account) Withdraw(amount int) error {
	// TODO
	return nil
}

// Transfer โอนเงิน — ถ้าไม่สำเร็จ ห้ามมีบัญชีไหนเปลี่ยน
func Transfer(from, to *Account, amount int) error {
	// TODO
	return nil
}

// ApplyBonus เพิ่ม bonus ให้ทุกบัญชีใน slice
// BUG: รันแล้วไม่มีบัญชีไหนได้ bonus
func ApplyBonus(accounts []Account, bonus int) {
	for _, a := range accounts {
		a.Balance += bonus
	}
}
