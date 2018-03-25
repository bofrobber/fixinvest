package models

import (
	"time"
)

//用什么DB好呢，用Redis还是mongodb
//这次学习一下redis的使用吧。redis没有分表，直接查询吧

//手续费
type Poundage struct {
	buy  float32
	sell float32
}

//基金基本信息
type FundBaseInfo struct {
	id      string
	name    string
	poudage Poundage
	other   string
}

//每次的投资记录，可以分成更小的。cost和share
type FundValue struct {
	cost       float32 //投资花费
	share      float32 //份额
	value      float32 //价值。一次投资时，基本等于cost = cost - Poundage.buy
	price      float32 //1份价格
	want_value float32 //总投资的预期总价值
}

//更新价格，主要是更新价值,返回价值的变化
func (p *FundValue) SetPrice(price float32) float32 {
	old_value := p.value

	p.price = price
	p.value = price * p.share

	//更新投资价值变更？新增的价值（贬值就是负数了）
	diff_value := p.value - old_value
	return diff_value
}

//投资记录
type FundInverstRecrd struct {
	time     string
	value    FundValue
	new_cost float32 //新增投資
	gains    float32 //收益
}

//基金汇总
type FundSummary struct {
	value               FundValue //基金价值
	total_poundage_cost Poundage  //手续费总花费
	max_more_cost       float32   //最大单次补仓--一般是大跌产生
	max_cost            float32   //最大投资
}

//定投前分析，包括补仓。收割怎么记录

//定投结果
type FixInvertInfo struct {
	cost  float32 //投资花费
	share float32 //份额

	//持仓成本不单独存储。如果需要再说
}

const (
	INVEST_DAY int = iota
	INVEST_WEEK
	INVEST_MONTH
)

type FixInvestPlan struct {
	cost        float32
	invert_type int
}

//收益情况 -- 主要是收割或者割肉的表现。
type FundGainsInfo struct {
	captial float32 //本金
	gains   float32 //收益
}

//基金
type Fund struct {
	info FundBaseInfo

	start_time     string
	summary        FundSummary
	invest_records []FundInverstRecrd

	summary_gans  FundGainsInfo //割肉后的本金也可以再次入市
	gains_records []FundGainsInfo

	cash float32 //现金账号。目前没有。这个本来可以计算压力

	plan FixInvestPlan //定投计划
}

//基本查询情况，记录基金信息和资产概要信息
type FundQueryInfo struct {
	info             FundBaseInfo
	summary          FundSummary
	tmp_gains        float32 //临时收益情况，用来判断
	next_invest_cost float32 //根据定投计划，计算下次要投资的钱
	next_gains_value float32 //下次收割价值。只有next_invest_cost=0 才可能有收益
}

func GetNewFund() *Fund {

	return &Fund{start_time: time.Now().String()}
}

func (p *Fund) SetFundPrice(price float32) bool {
	p.summary.value.SetPrice(price)
	return true
}
