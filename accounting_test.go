package accounting_test

import (
	"testing"
	"time"

	"github.com/fgrid/money"
        "github.com/openyard/accounting"
)

var (
	Jan4th1999, _  = time.Parse("2006-01-02", "1999-01-04")
	Oct1st1999, _  = time.Parse("2006-01-02", "1999-10-01")
	Oct5th1999, _  = time.Parse("2006-01-02", "1999-10-05")
	Oct15th1999, _ = time.Parse("2006-01-02", "1999-10-15")
	Nov1st1999, _  = time.Parse("2006-01-02", "1999-11-01")
	Dec1st1999, _  = time.Parse("2006-01-02", "1999-12-01")
	Dec5th1999, _  = time.Parse("2006-01-02", "1999-12-05")
	Dec15th1999, _ = time.Parse("2006-01-02", "1999-12-15")
	Jan4th2000, _  = time.Parse("2006-01-02", "2000-01-04")

	standard = accounting.NewServiceAgreement(10)
	poor     = accounting.NewServiceAgreement(10)
)

func init() {
	standard.AddPostingRule(accounting.EventTypeUsage, accounting.NewMultiplyByRatePR(accounting.EntryTypeBaseUsage), Oct1st1999)
	standard.AddPostingRule(accounting.EventTypeServiceCall, accounting.NewAmountFormulaPR(accounting.EntryTypeService, 0.5, money.EUR(100)), Oct1st1999)
	standard.AddPostingRule(accounting.EventTypeServiceCall, accounting.NewAmountFormulaPR(accounting.EntryTypeService, 0.5, money.EUR(150)), Dec1st1999)
	standard.AddPostingRule(accounting.EventTypeTax, accounting.NewAmountFormulaPR(accounting.EntryTypeTax, 0.055, money.EUR(0)), Oct1st1999)

	poor.AddPostingRule(accounting.EventTypeUsage, accounting.NewPoorCapPR(accounting.EntryTypeBaseUsage, 5, accounting.NewQuantity(50, accounting.UnitKWH)), Oct1st1999)
	poor.AddPostingRule(accounting.EventTypeServiceCall, accounting.NewAmountFormulaPR(accounting.EntryTypeService, 0, money.EUR(10)), Oct1st1999)
}

func TestUsage(t *testing.T) {
	acm := accounting.NewCustomer("Acme Coffee Makers", standard)

	event := accounting.NewUsageEvent(accounting.NewQuantity(50, accounting.UnitKWH), Oct1st1999, Oct1st1999, acm)
	event.Process()

	resultingEntry := acm.Entries(accounting.EntryTypeBaseUsage)[0]
	if "EUR 5.00" != resultingEntry.Amount.String() {
		t.Errorf("unexpected amount in entry: %s", resultingEntry.Amount.String())
	}
}

func TestService(t *testing.T) {
	acm := accounting.NewCustomer("Acme Coffee Makers", standard)

	event := accounting.NewMonetaryEvent(money.EUR(400), accounting.EventTypeServiceCall, Oct5th1999, Oct5th1999, acm)
	event.Process()

	resultingEntry := acm.Entries(accounting.EntryTypeService)[0]
	if "EUR 3.00" != resultingEntry.Amount.String() {
		t.Errorf("unexpected amount in entry: %s", resultingEntry.Amount.String())
	}
}

func TestLaterService(t *testing.T) {
	acm := accounting.NewCustomer("Acme Coffee Makers", standard)

	event := accounting.NewMonetaryEvent(money.EUR(400), accounting.EventTypeServiceCall, Dec5th1999, Dec15th1999, acm)
	event.Process()

	resultingEntry := acm.Entries(accounting.EntryTypeService)[0]
	if "EUR 3.50" != resultingEntry.Amount.String() {
		t.Errorf("unexpected amount in entry: %s", resultingEntry.Amount.String())
	}
}

func TestLowPayUsage(t *testing.T) {
	reggie := accounting.NewCustomer("Reginald Perrin", poor)

	event := accounting.NewUsageEvent(accounting.NewQuantity(50, accounting.UnitKWH), Oct1st1999, Oct1st1999, reggie)
	event.Process()
	event2 := accounting.NewUsageEvent(accounting.NewQuantity(51, accounting.UnitKWH), Nov1st1999, Nov1st1999, reggie)
	event2.Process()

	resultingEntry1 := reggie.Entries(accounting.EntryTypeBaseUsage)[0]
	if "EUR 2.50" != resultingEntry1.Amount.String() {
		t.Errorf("unexpected amount in entry: %s", resultingEntry1.Amount.String())
	}

	resultingEntry2 := reggie.Entries(accounting.EntryTypeBaseUsage)[1]
	if "EUR 5.10" != resultingEntry2.Amount.String() {
		t.Errorf("unexpected amount in entry: %s", resultingEntry2.Amount.String())
	}
}

func TestSecondaryUsage(t *testing.T) {
	acm := accounting.NewCustomer("Acme Coffee Makers", standard)

	event := accounting.NewUsageEvent(accounting.NewQuantity(50, accounting.UnitKWH), Oct1st1999, Oct1st1999, acm)
	event.Process()

	usageEntry := acm.Entries(accounting.EntryTypeBaseUsage)[0]
	if "EUR 5.00" != usageEntry.Amount.String() {
		t.Errorf("unexpected amount in entry: %s", usageEntry.Amount.String())
	}
	taxEntry := acm.Entries(accounting.EntryTypeTax)[0]
	if "EUR 0.27" != taxEntry.Amount.String() {
		t.Errorf("unexpected amount in entry: %s", taxEntry.Amount.String())
	}
}

func TestBalanceUsingTransactions(t *testing.T) {
	revenue := accounting.NewAccount(money.NewCurrency("EUR"))
	deferred := accounting.NewAccount(money.NewCurrency("EUR"))
	receivables := accounting.NewAccount(money.NewCurrency("EUR"))
	revenue.Withdraw(money.EUR(500), receivables, Jan4th1999)
	revenue.Withdraw(money.EUR(200), deferred, Jan4th1999)

	if "EUR 5.00" != receivables.Balance().String() {
		t.Errorf("unexpected balance for receivables: %s", receivables.Balance().String())
	}
	if "EUR 2.00" != deferred.Balance().String() {
		t.Errorf("unexpected balance for deferred: %s", deferred.Balance().String())
	}
	if "EUR -7.00" != revenue.Balance().String() {
		t.Errorf("unexpected balance for revenue: %s", revenue.Balance().String())
	}
}

func TestMultiLeggedTransactions(t *testing.T) {
	revenue := accounting.NewAccount(money.NewCurrency("EUR"))
	deferred := accounting.NewAccount(money.NewCurrency("EUR"))
	receivables := accounting.NewAccount(money.NewCurrency("EUR"))

	multi := accounting.NewMultiLeggedTransaction(Jan4th2000)
	multi.Add(money.EUR(700).Debit(), revenue)
	multi.Add(money.EUR(500), receivables)
	multi.Add(money.EUR(200), deferred)
	multi.Post()

	if "EUR 5.00" != receivables.Balance().String() {
		t.Errorf("unexpected balance for receivables: %s", receivables.Balance().String())
	}
	if "EUR 2.00" != deferred.Balance().String() {
		t.Errorf("unexpected balance for deferred: %s", deferred.Balance().String())
	}
	if "EUR -7.00" != revenue.Balance().String() {
		t.Errorf("unexpected balance for revenue: %s", revenue.Balance().String())
	}
}

func TestAdjustment(t *testing.T) {
	acm := accounting.NewCustomer("Acme Coffee Makers", standard)

	event := accounting.NewUsageEvent(accounting.NewQuantity(50, accounting.UnitKWH), Oct1st1999, Oct1st1999, acm)
	adjustment := accounting.NewUsageAdjustment(accounting.NewQuantity(70, accounting.UnitKWH), Oct1st1999, Oct15th1999, event)

	eventList := accounting.NewEventList()
	eventList.Add(event)
	eventList.Add(adjustment)
	eventList.Process()

	m1 := acm.BalanceFor(accounting.EntryTypeBaseUsage)
	if "EUR 7.00" != m1.String() {
		t.Errorf("unexpected amount in entry: %s", m1.String())
	}
	m2 := acm.BalanceFor(accounting.EntryTypeTax)
	if "EUR 0.38" != m2.String() {
		t.Errorf("unexpected amount in entry: %s", m2.String())
	}
}
