package accounting

import "github.com/fgrid/money"

type Customer struct {
	Name             string
	serviceAgreement *ServiceAgreement
	accounts         map[EntryType]*Account
}

func NewCustomer(name string, serviceAgreement *ServiceAgreement) *Customer {
	c := &Customer{
		Name:             name,
		serviceAgreement: serviceAgreement,
		accounts:         make(map[EntryType]*Account, 0),
	}
	for _, at := range EntryTypes {
		c.accounts[at] = NewAccount(money.NewCurrency("EUR"))
	}
	return c
}

func (c *Customer) BalanceFor(accountType EntryType) *money.Money {
	return c.accountFor(accountType).Balance()
}

func (c *Customer) AddEntry(entry *Entry, accountType EntryType) {
	c.accountFor(accountType).post(entry)
}

func (c *Customer) Entries(accountType EntryType) []*Entry {
	return c.accountFor(accountType).entries
}

func (c *Customer) accountFor(accountType EntryType) *Account {
	a, _ := c.accounts[accountType]
	return a
}
