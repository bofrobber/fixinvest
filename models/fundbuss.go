package models

import (
	"time"

	"github.com/astaxie/beego"
)

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

		//当前价值和
		info.tmp_gains = v.summary.value.price*v.summary.value.share - v.summary.value.want_value

		//收益超过了定投资金，产生收益了。
		if info.tmp_gains > v.plan.cost {
			info.next_gains_value = info.tmp_gains - v.plan.cost
			info.next_invest_cost = 0
		} else {
			info.next_gains_value = 0
			info.next_invest_cost = v.plan.cost - info.tmp_gains
		}

		ret[k] = info
	}

	return ret
}

//投资，买
func (p *FundBuss) Invest(fund_id string, cost float32) {
	beego.Info("FundBuss.Invest["+fund_id+"] value", cost)
	if p.funds == nil {
		return
	}

	v, ok := p.funds[fund_id]
	if !ok {
		return
	}

	//TODO，需要做投资记录，更新最新的投资情况,手续费怎么计算
	v.summary.value.cost += cost
	//原来的价值。这次的费用也算的预期价值，本钱都是要挣回来
	old_value := v.summary.value.price * v.summary.value.share
	//更新want_value，等于当前市值，或者加上当前手续费？
	v.summary.value.want_value = old_value + cost

	pound := cost * v.info.poudage.buy
	v.summary.total_poundage_cost.buy += pound
	v.summary.value.share += (cost - pound) / v.summary.value.price

	//更新最大投资
	if cost > v.summary.max_cost {
		v.summary.max_cost = cost
	}

	//记录投资
	record := FundInverstRecrd{}
	record.value = v.summary.value //
	record.new_cost = cost
	record.time = time.Now().String()
	//record.gains = 0  没有收益没有记录
	v.invest_records = append(v.invest_records, record)

	//正常情况，want_value = 定投计划*N期
}

//割肉或者收益，卖。实际按照操作规则，只有盈利会进行卖，越亏越买
//cost 是交易的总费用，包括盈利部分和手续费
//返回值，应该收割的份数，这样方便实际操作。实际买的时候就是按份数进行。
func (p *FundBuss) Gain(fund_id string, cost float32) float32 {
	beego.Info("FundBuss.Gain["+fund_id+"] value：", cost)
	if p.funds == nil {
		return 0
	}

	v, ok := p.funds[fund_id]
	if !ok {
		return 0
	}

	//收益和割肉怎么区分？ 临时计算该次是否有收益。如果有收益就不计算本金。如果割肉，就按本金算？
	info := FundGainsInfo{}
	cur_value := v.summary.value.price * v.summary.value.share
	tmp_gains := cur_value - v.summary.value.want_value

	v.summary.value.want_value = cur_value - cost //预期的钱肯定是少了。

	//减少份额，手续费怎么算？
	//正式的计算都是减少份额，收到的钱中，抽取份额
	//cost = share * price * (1 - pound)
	share := cost / (v.summary.value.price * (1 - v.info.poudage.sell))
	//费用
	pound := cost/(1-v.info.poudage.sell) - cost

	v.summary.value.share -= share
	v.summary.total_poundage_cost.sell += pound

	if tmp_gains > 0 {
		//挣钱
		//这个通常是收割
		//盈利范围内的将钱算到本金和盈利。
		if tmp_gains > cost { //最正常场景，会留下下一期的投资金额，本金也没有动的
			info.gains = cost
			info.captial = 0 //还没有动本金。预期金额高了。 正常场景是收割到预留下个投资周期的钱

		} else if tmp_gains == cost { //收割所有寸头，本金实际没有动
			info.gains = tmp_gains
			info.captial = 0
		} else {
			//多收割了，要算到本金里面
			info.gains = tmp_gains
			info.captial = cost - tmp_gains

			v.summary.value.cost -= info.captial //已经有本金出逃了。这个不算。不然下期就要算亏钱了
		}

	} else {
		//亏钱 还跑。肯定拿本金。正常需要补仓的。
		info.gains = 0
		info.captial = cost

		v.summary.value.cost -= info.captial //实际就是拿走本金，投资收益比更差？
	}

	v.gains_records = append(v.gains_records, info)

	return share
}

//如何解决纠正问题，纠正问题在于购买的时候没有进行交易。所以，有事后交易成功问题。
//纠正，就更新份额。 买入和卖出都是涉及份额，卖出不需要纠正份额。卖出纠正的是盈利情况
func (p *FundBuss) AdjustShare(fund_id string, share float32) {
	if p.funds == nil {
		return
	}

	v, ok := p.funds[fund_id]
	if !ok {
		return
	}

	v.summary.value.share = share //直接更新份额。投资的金额实际没有变化。

}

//卖出后，纠正收益情况，由于真实卖出的时候，实际是按份数操作。所以，份额不会错误。
//如果卖出价格符合预期，want_value就是正确。如果不一致，实际会影响到价值。但是，want_value是没有变化的。
//这个 不好操作，暂时不更新。
//卖出时按份额操作。所以，假设没有手续费，后面的price涨回上次购买预期，实际效果是一样的。不需要纠正。实际影响了手续费。
//这个暂时先不管

// func (p *FundBuss) AdjustGans(fund_id string, gans float32) {
// 	if p.funds == nil {
// 		return
// 	}

// 	v, ok := p.funds[fund_id]
// 	if !ok {
// 		return
// 	}

// }
