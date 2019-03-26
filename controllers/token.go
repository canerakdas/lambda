package controllers

import (
	"github.com/astaxie/beego"
	"github.com/juusechec/jwt-beego"
	"lambda/models"
	"time"
)

//  ListController operations for List
type TokenController struct {
	beego.Controller
}

// URLMapping ...
func (c *TokenController) URLMapping() {
	c.Mapping("Post", c.Post)
}

func (c *TokenController) Post() {
	username := c.Ctx.Input.Query("username")
	password := c.Ctx.Input.Query("password")

	tokenString := ""
	user, err := models.GetUserByEmail(username)
	if err == nil {
		if user.Confirmation == "Active" {
			if username == user.Email && password == user.Password {
				et := jwtbeego.EasyToken{
					Username: username,
					Expires:  time.Now().Unix() + 3600,
				}
				tokenString, _ = et.GetToken()
				c.Data["json"] = tokenString
				c.ServeJSON()
			} else {
				c.Ctx.Output.SetStatus(401)
				c.Data["json"] = "Access Denied"
				c.ServeJSON()
			}
		} else {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = "Confirmation failed"
			c.ServeJSON()
		}
	} else {
		c.Ctx.Output.SetStatus(401)
		c.Data["json"] = "Went wrong"
		c.ServeJSON()
	}
	return
}
