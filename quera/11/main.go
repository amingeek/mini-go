package main

type SavingsAccount struct {
	balance int
}

type CheckingAccount struct {
	balance int
}

type InvestmentAccount struct {
	balance int
}

func NewSavingsAccount() *SavingsAccount {
	// Code Here
	return &SavingsAccount{balance: 0}
}

func NewCheckingAccount() *CheckingAccount {
	// Code Here
	return &CheckingAccount{balance: 0}
}

func NewInvestmentAccount() *InvestmentAccount {
	// Code Here
	return &InvestmentAccount{balance: 0}
}

func (s *SavingsAccount) Deposit(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	s.balance += amount
	return "Success"
}

func (s *SavingsAccount) Withdraw(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	if s.balance < amount {
		return "Account balance is not enough"
	}
	s.balance -= amount
	return "Success"
}

func (s *SavingsAccount) CheckBalance() int {
	return s.balance
}

func (s *SavingsAccount) MonthlyInterest() int {
	return s.balance * 5 / 100 / 12
}

func (s *SavingsAccount) Transfer(receiver Account, amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}

	if s.balance < amount {
		return "Account balance is not enough"
	}

	switch receiver.(type) {
	case *SavingsAccount, *CheckingAccount, *InvestmentAccount:
	default:
		return "Invalid receiver account"
	}
	s.balance -= amount
	receiver.Deposit(amount)
	return "Success"
}

//---------------

func (s *CheckingAccount) Deposit(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	s.balance += amount
	return "Success"
}

func (s *CheckingAccount) Withdraw(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	if s.balance < amount {
		return "Account balance is not enough"
	}
	s.balance -= amount
	return "Success"
}

func (s *CheckingAccount) CheckBalance() int {
	return s.balance
}

func (s *CheckingAccount) MonthlyInterest() int {
	return s.balance * 1 / 100 / 12
}

func (s *CheckingAccount) Transfer(receiver Account, amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}

	if s.balance < amount {
		return "Account balance is not enough"
	}

	switch receiver.(type) {
	case *SavingsAccount, *CheckingAccount, *InvestmentAccount:
	default:
		return "Invalid receiver account"
	}
	s.balance -= amount
	receiver.Deposit(amount)
	return "Success"
}

//--------------

func (s *InvestmentAccount) Deposit(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	s.balance += amount
	return "Success"
}

func (s *InvestmentAccount) Withdraw(amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}
	if s.balance < amount {
		return "Account balance is not enough"
	}
	s.balance -= amount
	return "Success"
}

func (s *InvestmentAccount) CheckBalance() int {
	return s.balance
}

func (s *InvestmentAccount) MonthlyInterest() int {
	return s.balance * 2 / 100 / 12
}

func (s *InvestmentAccount) Transfer(receiver Account, amount int) string {
	if amount <= 0 {
		return "Amount cannot be negative"
	}

	if s.balance < amount {
		return "Account balance is not enough"
	}

	switch receiver.(type) {
	case *SavingsAccount, *CheckingAccount, *InvestmentAccount:
	default:
		return "Invalid receiver account"
	}
	s.balance -= amount
	receiver.Deposit(amount)
	return "Success"
}
