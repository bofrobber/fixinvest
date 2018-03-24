package main

import (
	_ "github.com/bofrobber/fixinvest/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.SetLogger("file", `{"filename":".\\logs\\test.log"}`) //这个方法使用了beego中默认的logs对象BeeLogger

	/*日志属性比较多
	设置的例子如下所示：

	logs.SetLogger(logs.AdapterFile, `{"filename":"test.log"}`)
	1
	主要的参数如下说明：

	filename 保存的文件名
	maxlines 每个文件保存的最大行数，默认值 1000000
	maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
	daily 是否按照每天 logrotate，默认是 true
	maxdays 文件最多保存多少天，默认保存 7 天
	rotate 是否开启 logrotate，默认是 true
	level 日志保存的时候的级别，默认是 Trace 级别
	perm 日志文件权限
	*/

	beego.SetLevel(beego.LevelInformational)
	beego.SetLogFuncCall(true) //日志包括文件名和行号。默认可以不打开
	//beego.BeeLogger.DelLogger("console") //可以把console 关闭了

	beego.Run()
}
