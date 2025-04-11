package payment_orders

import (
	"encoding/xml"
	"time"
)

type PaymentOrder struct {
	DebtorIban           string
	DebtorName           string
	CreditorIban         string
	CreditorName         string
	Amount               string
	Currency             string
	IdempotencyUniqueKey string
	Status               string
	CreatedAt            *time.Time
}

type PaymentFile struct {
	XMLName        xml.Name `xml:"Document"`
	Text           string   `xml:",chardata"`
	Xmlns          string   `xml:"xmlns,attr"`
	Xsi            string   `xml:"xsi,attr"`
	SchemaLocation string   `xml:"schemaLocation,attr"`
	GrpHdr         struct {
		Text    string `xml:",chardata"`
		MsgId   string `xml:"MsgId"`
		CreDtTm string `xml:"CreDtTm"`
	} `xml:"GrpHdr"`
	Cdtr struct {
		Text     string `xml:",chardata"`
		Nm       string `xml:"Nm"`
		CdtrAcct struct {
			Text string `xml:",chardata"`
			ID   struct {
				Text string `xml:",chardata"`
				IBAN string `xml:"IBAN"`
			} `xml:"Id"`
		} `xml:"CdtrAcct"`
	} `xml:"Cdtr"`
	Dbtr struct {
		Text     string `xml:",chardata"`
		Nm       string `xml:"Nm"`
		CdtrAcct struct {
			Text string `xml:",chardata"`
			ID   struct {
				Text string `xml:",chardata"`
				IBAN string `xml:"IBAN"`
			} `xml:"Id"`
		} `xml:"CdtrAcct"`
	} `xml:"Dbtr"`
	Amt struct {
		Text string `xml:",chardata"`
		Ccy  string `xml:"Ccy,attr"`
	} `xml:"Amt"`
}
