package models

//完成用户管理的几件基本事情：  -- 这个回头可以用一个默认用户，管理一下自己的信息就好了。以后再做多用户多个的支持
//
//预期的功能包括：
//注册，添加用户，并更新到数据库
//注销，删除用户
//查询用户，并登陆，在数据库中，需要记录用户密码
//登陆后，在内存中持有客户端的信息。
//用户的操作，需要同步更新到数据库

//目前先假定不存在访问的互斥问题。简单可以在外面加锁解决好了。
//用redis做数据库。这里实际可以不关心。
//beego目前应该是有redis的访问能力封装

//beego有这个session层
type UserSession struct {
	user  *User
	token string
}

type UserManager struct {
	users map[string]*UserSession
}

func (p *UserManager) Regist(user_name string, pwd string) {

}

func (p *UserManager) URegister(user_name string, pwd string) {
}

func (p *UserManager) Login(user_name string, pwd string) {

}

func (p *UserManager) GetUser(user_name string, token string) *User {
	return nil
}
