package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/juusechec/jwt-beego"
	"lambda/models"
	"math/rand"
	"strconv"
)

//  UserController operations for User
type UserController struct {
	beego.Controller
}

// URLMapping ...
func (c *UserController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Post ...
// @Title Post
// @Description create User
// @Param	body		body 	models.User	true		"body for User content"
// @Success 201 {int} models.User
// @Failure 403 body is empty
// @router / [post]
func (c *UserController) Post() {
	var v models.User
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	v.Confirmation = RandStringBytes(5)

	if _, err := models.AddUser(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		SendConfirmationEmail(v.Email, v.Confirmation, &v)
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UserController) GetOne() {
	token := c.Ctx.Input.Query("token")
	et := jwtbeego.EasyToken{}
	validation, email, _ := et.ValidateToken(token)

	//idStr := c.Ctx.Input.Param(":id")
	//id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetUserByEmail(email)

	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		if v.Email == email && validation {
			c.Data["json"] = v
		} else {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = "Access Denied"
		}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get User
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.User
// @Failure 403
// @router / [get]
func (c *UserController) GetAll() {
}

// Put ...
// @Title Put
// @Description update the User
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.User	true		"body for User content"
// @Success 200 {object} models.User
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UserController) Put() {
	token := c.Ctx.Input.Query("token")
	et := jwtbeego.EasyToken{}
	validation, email, _ := et.ValidateToken(token)

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.User{Id: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateUserById(&v); err == nil {
		if v.Email == email && validation {
			c.Data["json"] = "OK"
		} else {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = "Access Denied"
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the User
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UserController) Delete() {
	token := c.Ctx.Input.Query("token")
	et := jwtbeego.EasyToken{}
	validation, email, _ := et.ValidateToken(token)

	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.User{Id: id}

	if err := models.DeleteUser(id); err == nil {
		if v.Email == email && validation {
			c.Data["json"] = "OK"
		} else {
			c.Ctx.Output.SetStatus(401)
			c.Data["json"] = "Access Denied"
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
