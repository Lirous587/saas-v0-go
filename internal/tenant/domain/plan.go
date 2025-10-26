package domain

import "time"

type PlanType string

const PlanFreeType PlanType = "free"
const PlanCaringType PlanType = "caring"
const PlanProfessionalType PlanType = "professional"

type PlanStatus string

const PlanActiveStatus PlanStatus = "active"
const PlanInactiveStatus PlanStatus = "inactive"

type PlanBillingCycle string

const PlanMonthlyBillingCycle PlanBillingCycle = "active"
const PlanYearlyBillingCycle PlanBillingCycle = "yearly"
const PlanLifetimeBillingCycle PlanBillingCycle = "lifetime"

type Plan struct {
	TenantID     int64
	PlanType     PlanType
	StartTime    time.Time
	EndTime      time.Time
	Status       PlanStatus
	BillingCycle PlanBillingCycle
	CanUpgrade   bool
}

// type PlanLimits struct {
// 	ApiCalls int // API调用次数限制
// 	Plates   int // 板块数限制
// }
