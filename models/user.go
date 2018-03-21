package models

type User struct {
	buss []FundBuss

	//汇总情况,这个不要汇总算了。
	total_gains float32 //总收益

	name string
	pwd  string
}

///操作

func (p *User) AddFundBuss(name string) {

}

//这个不应该被调用，不然历史数据不就没有了。
func (p *User) RemoveFundBuss(name string) {

}

func (p *User) AddFund(buss_name string, fund_info FundBaseInfo) {

}

func (p *User) RemoveFund(buss_name string, fund_info FundBaseInfo) {

}
