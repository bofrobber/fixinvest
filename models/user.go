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

	return user
}

func (p *User) AddFundBuss(name string) bool {
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
	return true
}

func (p *User) RemoveFund(buss_name string, fund_id string) {
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
func (p *User) UpdateCost(buss_name string, fund_id string, poudage Poundage) {
	if p.buss == nil {
		beego.Error("RemoveFund buss is nil")
	}

	v, ok := p.buss[buss_name]
	if !ok {
		return
	}

	v.UpdatePoudage(fund_id, poudage)
}

//设置定投计划

//查询预计购买金额

//购买
