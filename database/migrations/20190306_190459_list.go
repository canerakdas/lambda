package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type List_20190306_190459 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &List_20190306_190459{}
	m.Created = "20190306_190459"

	migration.Register("List_20190306_190459", m)
}

// Run the migrations
func (m *List_20190306_190459) Up() {
	m.Charset = "utf8"
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE list(`id` int(11) NOT NULL AUTO_INCREMENT,`appid` int(11) DEFAULT NULL,`name` varchar(128) NOT NULL,`provider` varchar(128) NOT NULL,`time` int(11) DEFAULT NULL,`status` tinyint(1) NOT NULL,PRIMARY KEY (`id`)) CHARACTER SET=utf8;")
}

// Reverse the migrations
func (m *List_20190306_190459) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `list`")
}
