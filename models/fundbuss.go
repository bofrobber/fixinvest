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
