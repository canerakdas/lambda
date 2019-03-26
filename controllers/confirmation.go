package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"lambda/conf"
	"lambda/mail"
	"lambda/models"
	"log"
	"net/smtp"
	"strings"
)

//  ListController operations for List
type ConfirmationController struct {
	beego.Controller
}

// URLMapping ...
func (c *ConfirmationController) URLMapping() {
	c.Mapping("Post", c.Post)
}

func (c *ConfirmationController) Post() {
	email := c.Ctx.Input.Query("email")
	key := c.Ctx.Input.Query("key")
	user, err := models.GetUserByEmail(email)

	if key == "" {
		SendConfirmationEmail(email, user.Confirmation, user)
	} else {
		if user.Confirmation == key {
			user.Confirmation = "Active"
			models.UpdateUserById(user)
		}
	}

	if err != nil {
		log.Fatal(err)
	}
}

func SendConfirmationEmail(email string, key string, user *models.User) {

	templateData := struct {
		Name string
		URL  string
	}{
		Name: user.Realname,
		URL:  strings.Join([]string{"http://localhost:3000/confirmation/?email=", email, "&key=", key}, ""),
	}

	var auth = smtp.PlainAuth("", static.Sender.Email, static.Sender.Password, static.Sender.Host)

	r := mail.NewRequest([]string{email}, "Welcome to lambda", "Lambd", auth)
	err := r.ParseTemplate("./mail/template/welcome.html", templateData)
	if err := r.ParseTemplate("./mail/template/welcome.html", templateData); err == nil {
		ok, _ := r.SendEmail()
		fmt.Println(ok)
	}
	if err != nil {
		fmt.Println(err)
	}
}
