package models

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
	cost            float32 //投资花费
	share           float32 //份额
	value           float32 //价值。一次投资时，基本等于cost = cost - Poundage.buy
	one_share_value float32 //1份价格
	want_value      float32 //总投资的预期总价值
}

//投资记录
type FundInverstRecrd struct {
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
	invert_type float32
}

//收益情况 -- 主要是收割或者割肉的表现。
type FundGainsInfo struct {
	captial float32 //本金
	gains   float32 //本金
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
}

func GetNewFund() *Fund {

	return &Fund{}
}
