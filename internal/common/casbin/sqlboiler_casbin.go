package casbinadapter

import (
	"context"
	"database/sql"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"saas/internal/common/orm"
)

type SQLBoilerCasbinAdapter struct {
	db boil.Executor
}

type CasbinRule struct {
	ID    int
	Ptype string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

func NewSQLBoilerCasbinAdapter() *SQLBoilerCasbinAdapter {
	return &SQLBoilerCasbinAdapter{
		db: boil.GetDB(),
	}
}

// LoadPolicy 从存储中加载所有策略规则
func (ad *SQLBoilerCasbinAdapter) LoadPolicy(model model.Model) error {
	rules, err := orm.CasbinRules().AllG()
	if err != nil {
		zap.L().Error("casbin LoadPolicy失败", zap.Error(err))
		return err
	}

	for _, rule := range rules {
		if err := persist.LoadPolicyLine(ormCasbinRulesToStrings(rule), model); err != nil {
			return err
		}
	}
	return nil
}

// SavePolicy 将所有策略规则保存到存储中
func (ad *SQLBoilerCasbinAdapter) SavePolicy(model model.Model) error {
	tx, err := boil.BeginTx(context.Background(), &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  false,
	})

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1.清空数据库
	rows, err := orm.CasbinRules().DeleteAll(tx)
	if err != nil {
		return err
	}
	zap.L().Info("成功从psql清除数据", zap.Int64("rows", rows))

	// 2.从casbin中加载策略到数据库
	// 遍历 model，插入所有策略

	for _, astMap := range model {
		for ptype, ast := range astMap {
			if len(ast.Policy) == 0 {
				continue
			}
			for _, rule := range ast.Policy {
				vals := make([]string, 6)
				for i := 0; i < 6; i++ {
					if i < len(rule) {
						vals[i] = rule[i]
					} else {
						vals[i] = ""
					}
				}
				casbinRule := &CasbinRule{
					Ptype: ptype,
					V0:    vals[0],
					V1:    vals[1],
					V2:    vals[2],
					V3:    vals[3],
					V4:    vals[4],
					V5:    vals[5],
				}
				if err := domainCasbinRuleToOrm(casbinRule).Insert(tx, boil.Infer()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// AddPolicy adds a policy rule to the storage.
func (ad *SQLBoilerCasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	casbinRule := &CasbinRule{
		Ptype: ptype,
	}
	if len(rule) > 0 {
		casbinRule.V0 = rule[0]
	}
	if len(rule) > 1 {
		casbinRule.V1 = rule[1]
	}
	if len(rule) > 2 {
		casbinRule.V2 = rule[2]
	}
	if len(rule) > 3 {
		casbinRule.V3 = rule[3]
	}
	if len(rule) > 4 {
		casbinRule.V4 = rule[4]
	}
	if len(rule) > 5 {
		casbinRule.V5 = rule[5]
	}
	return domainCasbinRuleToOrm(casbinRule).Insert(ad.db, boil.Infer())
}

// RemovePolicy removes a policy rule from the storage.
func (ad *SQLBoilerCasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	vals := make([]string, 6)
	for i := 0; i < 6; i++ {
		if i < len(rule) {
			vals[i] = rule[i]
		} else {
			vals[i] = ""
		}
	}
	casbinRule := &CasbinRule{
		Ptype: ptype,
		V0:    vals[0],
		V1:    vals[1],
		V2:    vals[2],
		V3:    vals[3],
		V4:    vals[4],
		V5:    vals[5],
	}

	ormCasbinRule := domainCasbinRuleToOrm(casbinRule)

	mods := make([]qm.QueryMod, 0, 6)
	mods = append(mods, orm.CasbinRuleWhere.Ptype.EQ(ormCasbinRule.Ptype))
	if ormCasbinRule.V0.Valid {
		mods = append(mods, orm.CasbinRuleWhere.V0.EQ(ormCasbinRule.V0))
	} else {
		mods = append(mods, orm.CasbinRuleWhere.V0.IsNull())
	}
	if ormCasbinRule.V1.Valid {
		mods = append(mods, orm.CasbinRuleWhere.V1.EQ(ormCasbinRule.V1))
	} else {
		mods = append(mods, orm.CasbinRuleWhere.V1.IsNull())
	}
	if ormCasbinRule.V2.Valid {
		mods = append(mods, orm.CasbinRuleWhere.V2.EQ(ormCasbinRule.V2))
	} else {
		mods = append(mods, orm.CasbinRuleWhere.V2.IsNull())
	}
	if ormCasbinRule.V3.Valid {
		mods = append(mods, orm.CasbinRuleWhere.V3.EQ(ormCasbinRule.V3))
	} else {
		mods = append(mods, orm.CasbinRuleWhere.V3.IsNull())
	}
	if ormCasbinRule.V4.Valid {
		mods = append(mods, orm.CasbinRuleWhere.V4.EQ(ormCasbinRule.V4))
	} else {
		mods = append(mods, orm.CasbinRuleWhere.V4.IsNull())
	}
	if ormCasbinRule.V5.Valid {
		mods = append(mods, orm.CasbinRuleWhere.V5.EQ(ormCasbinRule.V5))
	} else {
		mods = append(mods, orm.CasbinRuleWhere.V5.IsNull())
	}

	_, err := orm.CasbinRules(mods...).DeleteAll(ad.db)
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (ad *SQLBoilerCasbinAdapter) RemoveFilteredPolicy(
	sec string, ptype string, fieldIndex int, fieldValues ...string,
) error {
	return errors.New("not implemented")
}
