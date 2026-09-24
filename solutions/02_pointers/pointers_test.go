package pointers

import (
	"errors"
	"testing"
)

func TestSwap(t *testing.T) {
	a, b := 1, 2
	Swap(&a, &b)
	if a != 2 || b != 1 {
		t.Fatalf("after Swap a=%d b=%d; want a=2 b=1", a, b)
	}
}

func TestDeposit(t *testing.T) {
	acc := Account{Owner: "tah", Balance: 100}
	if err := acc.Deposit(50); err != nil {
		t.Fatal(err)
	}
	if acc.Balance != 150 {
		t.Fatalf("Balance = %d; want 150 (hint: value receiver?)", acc.Balance)
	}
	if err := acc.Deposit(0); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("Deposit(0) error = %v; want ErrInvalidAmount", err)
	}
}

func TestWithdraw(t *testing.T) {
	acc := NewAccount("tah", 100)
	if err := acc.Withdraw(30); err != nil {
		t.Fatal(err)
	}
	if acc.Balance != 70 {
		t.Fatalf("Balance = %d; want 70", acc.Balance)
	}
	if err := acc.Withdraw(1000); !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("Withdraw(1000) error = %v; want ErrInsufficientFunds", err)
	}
	if err := acc.Withdraw(-5); !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("Withdraw(-5) error = %v; want ErrInvalidAmount", err)
	}
	if acc.Balance != 70 {
		t.Errorf("Balance changed after failed withdraw: %d", acc.Balance)
	}
}

func TestTransfer(t *testing.T) {
	from := NewAccount("a", 100)
	to := NewAccount("b", 0)
	if err := Transfer(from, to, 40); err != nil {
		t.Fatal(err)
	}
	if from.Balance != 60 || to.Balance != 40 {
		t.Fatalf("from=%d to=%d; want 60, 40", from.Balance, to.Balance)
	}
	if err := Transfer(from, to, 999); !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("Transfer(999) error = %v; want ErrInsufficientFunds", err)
	}
	if from.Balance != 60 || to.Balance != 40 {
		t.Errorf("balances changed after failed transfer: from=%d to=%d", from.Balance, to.Balance)
	}
}

func TestApplyBonus(t *testing.T) {
	accs := []Account{{"a", 10}, {"b", 20}}
	ApplyBonus(accs, 5)
	if accs[0].Balance != 15 || accs[1].Balance != 25 {
		t.Fatalf("got %+v; want balances 15 and 25 (hint: range copies)", accs)
	}
}
