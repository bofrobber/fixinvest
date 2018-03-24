package models

import "github.com/astaxie/beego"

type FundBuss struct {
	name  string           //什么地方
	funds map[string]*Fund //key是基金id

	//汇总情况
	total_gains float32 //总收益
}

func GetNewFundBuss(name string) *FundBuss {
	bus := &FundBuss{}
	bus.name = name
	bus.funds = make(map[string]*Fund)
	bus.total_gains = 0 //可以不设置用默认值

	return bus
}

func (p *FundBuss) AddFund(info FundBaseInfo) {
	beego.Info("FundBuss.AddFund", info)
	if p.funds == nil {
		return
	}

	_, ok := p.funds[info.id]
	if !ok {
		p.funds[info.id] = GetNewFund()
	}

}

func (p *FundBuss) RemoveFund(fund_id string) {
	if p.funds == nil {
		return
	}

	delete(p.funds, fund_id)
}

func (p *FundBuss) UpdateFundId(old_fund_id string, new_fund_id string) int {
	if p.funds == nil {
		beego.Error("UpdateFund funds is nil")
		return -11
	}

	v, ok := p.funds[old_fund_id]
	if !ok {
		return -12
	}

	_, ok_new := p.funds[new_fund_id]
	if ok_new {
		return -13
	}

	p.funds[new_fund_id] = v

	delete(p.funds, old_fund_id)
	return 0
}

func (p *FundBuss) UpdatePoudage(fund_id string, poudage Poundage) int {
	if p.funds == nil {
		beego.Error("UpdateFund funds is nil")
		return -11
	}

	v, ok := p.funds[fund_id]
	if !ok {
		return -12
	}

	v.info.poudage = poudage

	return 0
}

func (p *FundBuss) SetFundInvestPlan(fund_id string, plan FixInvestPlan) bool {
	beego.Info("FundBuss.SetFundInvestPlan", plan)
	if p.funds == nil {
		return false
	}

	v, ok := p.funds[fund_id]
	if !ok {
		return false
	}

	v.plan = plan

	return true
}

func (p *FundBuss) SetFundPrice(fund_id string, price float32) bool {
	beego.Info("FundBuss.SetFundPrice["+fund_id+"] price", price)
	if p.funds == nil {
		return false
	}

	v, ok := p.funds[fund_id]
	if !ok {
		return false
	}

	return v.SetFundPrice(price)
}

//查询盈利情况
func (p *FundBuss) GetAllFundSummaryInfo() map[string]*FundQueryInfo {
	ret := make(map[string]*FundQueryInfo)

	for k, v := range p.funds {
		info := &FundQueryInfo{info: v.info, summary: v.summary}
		ret[k] = info
	}

	return ret
}

//投资
func (p *FundBuss) Invest(fund_id string, value float32) {
	beego.Info("FundBuss.Invest["+fund_id+"] value", value)
	if p.funds == nil {
		return
	}

	v, ok := p.funds[fund_id]
	if !ok {
		return
	}

	//TODO，需要做投资记录，更新最新的投资情况,手续费怎么计算
	v.summary.value.cost += value
	//v.summary.value.share += （cost - 费用）/price
	//更新want_value
	//更新最大投资
}

//割肉或者收益
//func(0)
