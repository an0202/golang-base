/**
 * @Author: jie.an
 * @Description:
 * @File:  aws_billing.go
 * @Version: 1.0.0
 * @Date: 2020/4/26 10:08 上午
 */

package csv

import (
	"strings"
)

type awsBilling struct {
	InvoiceID        string
	PayerAccountId   string
	LinkedAccountId  string
	RecordType       string
	RecordId         string
	ProductName      string
	RateId           string
	SubscriptionId   string
	PricingPlanId    string
	UsageType        string
	Operation        string
	AvailabilityZone string
	ReservedInstance string
	ItemDescription  string
	UsageStartDate   string
	UsageEndDate     string
	UsageQuantity    string
	BlendedRate      string
	BlendedCost      string
	UnBlendedRate    string
	UnBlendedCost    string
	ResourceId       string
}

type MSPBillings struct {
	Count       int
	MSPBillings map[string]MSPBilling
}

func NewMSPBillings() *MSPBillings {
	var ms MSPBillings
	ms.MSPBillings = make(map[string]MSPBilling)
	return &ms
}

func (ms *MSPBillings) ProcessMSPBillings(record []string) {
	m := NewMSPBilling()
	m.LinkedAccountId = record[2]
	m.UsageType = record[9]
	m.UsageStartDate = record[14]
	m.ResourceId = record[21]
	//update data if resourceId exist
	if _, ok := ms.MSPBillings[m.ResourceId]; ok {
		oldContent := ms.MSPBillings[m.ResourceId]
		oldContent.UsageStartDate = m.UsageStartDate
		oldContent.updateRunningDates()
		ms.MSPBillings[m.ResourceId] = oldContent
	} else {
		//add to map if not exist
		m.updateRunningDates()
		ms.MSPBillings[m.ResourceId] = *m
	}
}

type MSPBilling struct {
	LinkedAccountId string
	UsageType       string
	ResourceId      string
	UsageStartDate  string
	RunningDates    map[string]bool
	RunningDays     int
}

func NewMSPBilling() *MSPBilling {
	var m MSPBilling
	m.RunningDates = make(map[string]bool)
	return &m
}

func (m *MSPBilling) updateRunningDates() {
	//split string
	m.UsageStartDate = strings.Fields(m.UsageStartDate)[0]
	if _, ok := m.RunningDates[m.UsageStartDate]; !ok {
		m.RunningDates[m.UsageStartDate] = true
	}
	m.RunningDays = len(m.RunningDates)
}
