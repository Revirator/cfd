package external

import (
	"encoding/json"
	"fmt"
	"log"
)

type EdgarEntry struct {
	Facts FinancialFacts `json:"facts"`
}

type FinancialFacts struct {
	Entity     EntityInformation `json:"dei"`
	Principles Principles        `json:"us-gaap"`
}

type EntityInformation struct {
	EntityCommonStockSharesOutstanding *Metric
}

type Principles struct {
	Cash                                  *Metric // TODO: check why some companies are missing this
	CashAndCashEquivalentsAtCarryingValue *Metric
	CommonStockSharesOutstanding          *Metric // TODO: might not be up to date?
	CostsAndExpenses                      *Metric
	EarningsPerShareDiluted               *Metric
	LongTermDebt                          *Metric
	NetIncomeLoss                         *Metric
	PaymentsOfDividends                   *Metric
	PaymentsOfDividendsCommonStock        *Metric
	Revenues                              *Metric
	ShortTermInvestments                  *Metric
}

type Metric struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Units       Units  `json:"units"`
}

type Units struct {
	PrimaryEntries   []FinancialDataEntry `json:"usd"`
	SecondaryEntries []FinancialDataEntry `json:"shares"`
	TertiaryEntries  []FinancialDataEntry `json:"usd/shares"`
}

type FinancialDataEntry struct {
	Value float64 `json:"val"`
	Frame string  `json:"frame"`
	Form  string  `json:"form"`
}

func (entry FinancialDataEntry) IsQuarterlyReport() bool {
	return entry.Form == "10-Q" || entry.Form == "10-K" || entry.Form == "10-Q/A"
}

func (entry FinancialDataEntry) IsAnnualReport() bool {
	return entry.Form == "10-K"
}

const (
	EDGAR_HOST             = "www.sec.gov"
	EDGAR_COMPANY_DATA_URL = "https://data.sec.gov/api/xbrl/companyfacts/CIK%s.json"
)

func GetFinancialFactsForCompanyGivenCIK(cik string) *FinancialFacts {
	request := prepareRequest("GET", fmt.Sprintf(EDGAR_COMPANY_DATA_URL, cik))
	request.Header.Add("Host", EDGAR_HOST)
	body := sendRequestAndGetBody(request)

	data := EdgarEntry{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err)
		return nil
	}

	return &data.Facts
}
