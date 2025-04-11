package payment_orders

import (
	"encoding/xml"
	"fmt"
	"os"
)

type BankFileSender struct {
	outputLocation string
}

func NewBankFileSender(outputLocation string) *BankFileSender {
	return &BankFileSender{
		outputLocation: outputLocation,
	}
}

func (g BankFileSender) SendFileToBank(po *PaymentOrder) error {
	file := PaymentFile{
		GrpHdr: struct {
			Text    string `xml:",chardata"`
			MsgId   string `xml:"MsgId"`
			CreDtTm string `xml:"CreDtTm"`
		}{MsgId: "Message-ID", CreDtTm: po.CreatedAt.String()},
		Amt: struct {
			Text string `xml:",chardata"`
			Ccy  string `xml:"Ccy,attr"`
		}{Text: po.Amount, Ccy: po.Currency},
		Dbtr: struct {
			Text     string `xml:",chardata"`
			Nm       string `xml:"Nm"`
			CdtrAcct struct {
				Text string `xml:",chardata"`
				ID   struct {
					Text string `xml:",chardata"`
					IBAN string `xml:"IBAN"`
				} `xml:"Id"`
			} `xml:"CdtrAcct"`
		}{Nm: po.DebtorName, CdtrAcct: struct {
			Text string `xml:",chardata"`
			ID   struct {
				Text string `xml:",chardata"`
				IBAN string `xml:"IBAN"`
			} `xml:"Id"`
		}{ID: struct {
			Text string `xml:",chardata"`
			IBAN string `xml:"IBAN"`
		}{IBAN: po.DebtorIban}}},
		Cdtr: struct {
			Text     string `xml:",chardata"`
			Nm       string `xml:"Nm"`
			CdtrAcct struct {
				Text string `xml:",chardata"`
				ID   struct {
					Text string `xml:",chardata"`
					IBAN string `xml:"IBAN"`
				} `xml:"Id"`
			} `xml:"CdtrAcct"`
		}{Nm: po.CreditorName, CdtrAcct: struct {
			Text string `xml:",chardata"`
			ID   struct {
				Text string `xml:",chardata"`
				IBAN string `xml:"IBAN"`
			} `xml:"Id"`
		}{ID: struct {
			Text string `xml:",chardata"`
			IBAN string `xml:"IBAN"`
		}{IBAN: po.CreditorIban}}},
	}
	b, err := xml.Marshal(file)
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%s/%s.request.xml", g.outputLocation, po.IdempotencyUniqueKey)
	fmt.Println("fileName:", fileName)
	err = os.WriteFile(fileName, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
