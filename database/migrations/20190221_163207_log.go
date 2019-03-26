package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Log_20190221_163207 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Log_20190221_163207{}
	m.Created = "20190221_163207"

	migration.Register("Log_20190221_163207", m)
}

// Run the migrations
func (m *Log_20190221_163207) Up() {
	m.Charset = "utf8"
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE log(`id` int(11) NOT NULL AUTO_INCREMENT,`time` varchar(128) NOT NULL,`external_id` varchar(128) NOT NULL,`description` varchar(128) NOT NULL,PRIMARY KEY (`id`)) CHARACTER SET=utf8;")
}

// Reverse the migrations
func (m *Log_20190221_163207) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `log`")
}
