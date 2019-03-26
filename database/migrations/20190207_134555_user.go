package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type User_20190207_134555 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &User_20190207_134555{}
	m.Created = "20190207_134555"
	migration.Register("User_20190207_134555", m)
}

// Run the migrations
func (m *User_20190207_134555) Up() {
	m.Charset = "utf8"
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE user(`id` int(11) NOT NULL AUTO_INCREMENT,`steamid` longtext  NOT NULL,`communityvisibilitystate` int(11) DEFAULT NULL,`profilestate` int(11) DEFAULT NULL,`personaname` longtext  NOT NULL,`lastlogoff` int(11) DEFAULT NULL,`profileurl` longtext  NOT NULL,`avatar` longtext  NOT NULL,`avatarmedium` longtext  NOT NULL,`avatarfull` longtext  NOT NULL,`personastate` int(11) DEFAULT NULL,`realname` longtext  NOT NULL,`primaryclanid` longtext  NOT NULL,`timecreated` int(11) DEFAULT NULL,`personastateflags` int(11) DEFAULT NULL,`loccountrycode` longtext  NOT NULL,`locstatecode` longtext  NOT NULL,`loccityid` int(11) DEFAULT NULL,`email` VARCHAR(100) NOT NULL,`confirmation` VARCHAR(10) NOT NULL,`password` longtext  NOT NULL,`notifications` longtext  NOT NULL,PRIMARY KEY (`id`),UNIQUE KEY (email)) CHARACTER SET=utf8;")
}

// Reverse the migrations
func (m *User_20190207_134555) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `user`")
}
