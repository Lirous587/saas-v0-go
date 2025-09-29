package casbinadapter

import (
	"saas/internal/common/orm"
	"strings"
)

func ormCasbinRulesToStrings(rule *orm.CasbinRule) string {
	vals := []string{rule.Ptype}
	// 依次追加所有非空字段
	if rule.V0.Valid {
		vals = append(vals, rule.V0.String)
	}
	if rule.V1.Valid {
		vals = append(vals, rule.V1.String)
	}
	if rule.V2.Valid {
		vals = append(vals, rule.V2.String)
	}
	if rule.V3.Valid {
		vals = append(vals, rule.V3.String)
	}
	if rule.V4.Valid {
		vals = append(vals, rule.V4.String)
	}
	if rule.V5.Valid {
		vals = append(vals, rule.V5.String)
	}
	return strings.Join(vals, ", ")
}

// nolint
func ormCasbinRuleToDomain(orm *orm.CasbinRule) *CasbinRule {
	casbinRule := &CasbinRule{
		ID:    orm.ID,
		Ptype: orm.Ptype,
	}

	if orm.V0.Valid {
		casbinRule.V0 = orm.V0.String
	}
	if orm.V1.Valid {
		casbinRule.V1 = orm.V1.String
	}
	if orm.V2.Valid {
		casbinRule.V2 = orm.V2.String
	}
	if orm.V3.Valid {
		casbinRule.V3 = orm.V3.String
	}
	if orm.V4.Valid {
		casbinRule.V4 = orm.V4.String
	}
	if orm.V5.Valid {
		casbinRule.V5 = orm.V5.String
	}

	return casbinRule
}

func domainCasbinRuleToOrm(domain *CasbinRule) *orm.CasbinRule {
	ormCasbinRule := &orm.CasbinRule{
		ID:    domain.ID,
		Ptype: domain.Ptype,
	}

	if domain.V0 != "" {
		ormCasbinRule.V0.String = domain.V0
		ormCasbinRule.V0.Valid = true
	}
	if domain.V1 != "" {
		ormCasbinRule.V1.String = domain.V1
		ormCasbinRule.V1.Valid = true
	}
	if domain.V2 != "" {
		ormCasbinRule.V2.String = domain.V2
		ormCasbinRule.V2.Valid = true
	}
	if domain.V3 != "" {
		ormCasbinRule.V3.String = domain.V3
		ormCasbinRule.V3.Valid = true
	}
	if domain.V4 != "" {
		ormCasbinRule.V4.String = domain.V4
		ormCasbinRule.V4.Valid = true
	}
	if domain.V5 != "" {
		ormCasbinRule.V5.String = domain.V5
		ormCasbinRule.V5.Valid = true
	}

	return ormCasbinRule
}
