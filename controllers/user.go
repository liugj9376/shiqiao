package controllers

import (
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"regexp"
	"shiqiao/models"
	"strconv"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowReg() {
	this.TplName = "register.html"
}

func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")

	temp, _ := base64.StdEncoding.DecodeString(userName)

	if string(temp) == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = string(temp)
		this.Data["checked"] = "checked"
	}

	this.TplName = "login.html"
}

func (this *UserController) HandleReg() {
	//1.获取数据
	name := this.GetString("user_name")
	pwd := this.GetString("pwd")
	cpwd := this.GetString("cpwd")
	email := this.GetString("email")
	//allow:=this.GetString("allow")
	//2.数据校验
	if name == "" || pwd == "" || cpwd == "" || email == "" {
		this.Data["errmsg"] = "数据不完整,请重新注册~"
		this.TplName = "register.html"
		return
	}

	if pwd != cpwd {
		this.Data["errmsg"] = "两次输入密码不一致,请重新输入"
		this.TplName = "register.html"
		return
	}

	reg, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(email)
	if res == "" {
		this.Data["errmsg"] = "邮箱格式不正确"
		this.TplName = "register.html"
		return
	}

	//3.数据处理
	orm := orm.NewOrm()
	var user models.User
	user.Name = name
	user.PassWord = pwd
	user.Email = email
	_, err := orm.Insert(&user)

	if err != nil {
		this.Data["errmsg"] = "注册失败,请更换数据注册"
		this.TplName = "register.html"
		return
	}

	//4.发送邮件  ohoufzcinvcgdidc
	emailConfig := `{"username":"2669423290@qq.com","password":"ohoufzcinvcgdidc","host":"smtp.qq.com","port":587}`
	emailConn := utils.NewEMail(emailConfig)
	emailConn.From = "2669423290@qq.com"
	emailConn.To = []string{email}
	emailConn.Subject = "天天生鲜用户注册"
	emailConn.Text = "http://192.168.0.104:8080/active?id=" + strconv.Itoa(user.Id)
	emailConn.Send()

	//5.返回视图
	this.Ctx.WriteString("注册成功，请去相应邮箱激活用户！")
}

func (this *UserController) ActiveUser() {
	//1.获取参数
	id, err := this.GetInt("id")

	if err != nil {
		this.Data["errmsg"] = "要激活的邮箱不存在"
		this.TplName = "register.html"
		return
	}

	//2.查询数据
	o := orm.NewOrm()
	var user models.User
	user.Id = id
	err = o.Read(&user)
	if err != nil {
		this.Data["errmsg"] = "激活失败"
		this.TplName = "register.html"
		return
	}
	user.Active = true
	//修改数据
	o.Update(&user)

	this.Redirect("/login", 302)
}

func (this *UserController) HandleLogin() {
	//1.获取参数
	name := this.GetString("username")
	pwd := this.GetString("pwd")

	//2.数据校验
	if name == "" || pwd == "" {
		this.Data["errmsg"] = "登录信息不能为空"
		this.TplName = "login.html"
		return
	}

	//3.数据查询
	o := orm.NewOrm()
	var user models.User
	user.Name = name
	//user.PassWord=pwd
	//user.Active=true

	err := o.Read(&user, "name")
	if err != nil {
		this.Data["errmsg"] = "用户名或密码错误，请重新输入！"
		this.TplName = "login.html"
		return
	}

	if user.PassWord != pwd {
		this.Data["errmsg"] = "用户名或密码错误，请重新输入！"
		this.TplName = "login.html"
		return
	}

	if user.Active != true {
		this.Data["errmsg"] = "用户没有激活邮箱,请先激活邮箱！"
		this.TplName = "login.html"
		return
	}

	remember := this.GetString("remember")
	if remember == "on" {
		temp := base64.StdEncoding.EncodeToString([]byte(name))
		this.Ctx.SetCookie("userName", temp, 24*3600*30)
	} else {
		this.Ctx.SetCookie("userName", name, -1)
	}
	this.Ctx.WriteString("登录成功")
}
