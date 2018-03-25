package models

import (
	"github.com/astaxie/beego"
)

type User struct {
	buss map[string]*FundBuss

	//汇总情况,这个不要汇总算了。
	total_gains float32 //总收益

	name string
	pwd  string

	plan FixInvestPlan
}

/*操作
主要包括：
1. 查询。当前结果。投资记录。投资盈利情况
2. 设置单价（股价，或者基金每股市值）
3. 下次定投资金
4. 添加：
	4.1 添加基金购买公司
	4.2 添加基金
5. 更新
	5.1 更新操作费用（这个在每次都有算了）
	5.2 更新公司信息
	5.3 更新基金信息
	5.4 更新基金购买费率

5. 单纯的操作记录
*/

func GetNewUser() *User {
	user := &User{}
	user.buss = make(map[string]*FundBuss)
	user.plan.invert_type = INVEST_WEEK
	user.plan.cost = 300 //默认300

	return user
}

func (p *User) AddFundBuss(name string) bool {
	beego.Info("AddFundBuss:", name)

	if p.buss == nil {
		return false
	}

	_, ok := p.buss[name]

	if !ok {
		p.buss[name] = GetNewFundBuss(name)
	}

	return true //也可以返回 !ok ，必需是空的。这里先都成功
}

//这个不应该被调用，不然历史数据不就没有了。
func (p *User) RemoveFundBuss(name string) {
	beego.Info("RemoveFundBuss:", name)

	if p.buss == nil {
		return
	}

	delete(p.buss, name)
}

func (p *User) AddFund(buss_name string, fund_info FundBaseInfo) bool {
	beego.Info("AddFund [", buss_name, "]", "[", fund_info, "]")

	if p.buss == nil {
		beego.Error("AddFund buss is nil")
		return false
	}

	v, ok := p.buss[buss_name]
	if !ok {
		beego.Error("AddFund [", buss_name, "] don't find fund buss")
		return false
	}

	v.AddFund(fund_info)
	v.SetFundInvestPlan(fund_info.id, p.plan)
	return true
}

func (p *User) RemoveFund(buss_name string, fund_id string) {
	beego.Info("RemoveFund [", buss_name, "]", "[", fund_id, "]")

	if p.buss == nil {
		beego.Error("RemoveFund buss is nil")
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return
	}

	v.RemoveFund(fund_id)
}

//先做简单的更新信息,只更新基本信息，基本信息包括基金id，注意，基金id更新后，也不能重复，比如360上市后，用了的360基金就更换id了
//还有就是费率，其他暂时不更新
func (p *User) UpdateFundId(buss_name string, old_fund_id string, new_fund_id string) {
	beego.Info("UpdateFundId [", buss_name, "]", "old:[", old_fund_id, "] new:[", new_fund_id, "]")

	if p.buss == nil {
		beego.Error("RemoveFund buss is nil")
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return
	}

	v.UpdateFundId(old_fund_id, new_fund_id)
}

//更新费率
func (p *User) UpdatePoudage(buss_name string, fund_id string, poudage Poundage) {
	beego.Info("UpdatePoudage:", buss_name, "poundage:", poudage)

	if p.buss == nil {
		beego.Error("RemoveFund buss is nil")
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return
	}

	v.UpdatePoudage(fund_id, poudage)
}

/*设置定投计划
type 定投周期
cost 定投额度，默认所有的都是这个。具体也可以分开
*/
func (p *User) SetInvestPlan(plan_type int, cost float32) {
	beego.Info("default SetInvestPlan type:", plan_type, "cost:", cost)

	p.plan.invert_type = plan_type
	p.plan.cost = cost
}

func (p *User) SetFundInvestPlan(buss_name, fund_id string, plan_type int, cost float32) bool {
	beego.Info("SetFundInvestPlan type:", plan_type, "cost:", cost, "buss:[", buss_name, "][", fund_id, "]")

	if p.buss == nil {
		beego.Error("RemoveFund buss is nil")
		return false
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return false
	}

	plan := FixInvestPlan{invert_type: plan_type, cost: cost}

	return v.SetFundInvestPlan(fund_id, plan)
}

//设置基金价格，一个一个设置好了。也可以支持批量设置
func (p *User) SetFundPrice(buss_name, fund_id string, price float32) bool {
	beego.Info("SetFundPrice buss[", buss_name, "][", fund_id, "] price:", price)
	if p.buss == nil {
		beego.Error("SetCurInvest buss is nil")
		return false
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return false
	}

	return v.SetFundPrice(fund_id, price)
}

//查询盈利情况
func (p *User) GetAllFundSummaryInfo(buss_name string) map[string]*FundQueryInfo {
	beego.Info("GetAllFundSummaryInfo buss[", buss_name, "]")

	//返回情况
	if p.buss == nil {
		beego.Error("SetCurInvest buss is nil")
		return nil
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return nil
	}

	return v.GetAllFundSummaryInfo()
}

func (p *User) GetAllFundInfo(buss_name string) map[string]*Fund {
	beego.Info("GetAllFundInfo buss[", buss_name, "]")

	//返回情况 所有信息了
	if p.buss == nil {
		beego.Error("SetCurInvest buss is nil")
		return nil
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return nil
	}

	return v.funds
}

//定投变更怎么计算，这个也是要搞懂，目前这个反正是比较简单
//投资，如果不按规定补仓，回导致的结果是want_value变少了。
func (p *User) Invest(buss_name string, fund_id string, value float32) {
	beego.Info("Invest buss[", buss_name, "][", fund_id, "] cost:", value)

	if p.buss == nil {
		beego.Error("Invest buss is nil")
		return
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return
	}

	v.Invest(fund_id, value)

}

//实际要转换成份额。为了方便统计，使用价值，在查询完成后计算
func (p *User) Gain(buss_name string, fund_id string, value float32) float32 {
	beego.Info("Gain buss[", buss_name, "][", fund_id, "] cost:", value)

	if p.buss == nil {
		beego.Error("Gain buss is nil")
		return 0
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return 0
	}

	return v.Gain(fund_id, value)
}

func (p *User) AdjustShare(buss_name string, fund_id string, share float32) {
	beego.Info("AdjustShare buss[", buss_name, "][", fund_id, "] share:", share)

	if p.buss == nil {
		beego.Error("AdjustShare buss is nil")
		return
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return
	}

	v.AdjustShare(fund_id, share)
}
