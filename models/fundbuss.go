package models

type FundBuss struct{
	name string //什么地方
	funds map[string]Fund //key是基金id

	//汇总情况
	total_gains float32 //总收益
}




