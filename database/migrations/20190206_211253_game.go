package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Game_20190206_211253 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Game_20190206_211253{}
	m.Created = "20190206_211253"
	migration.Register("Game_20190206_211253", m)
}

// Run the migrations
func (m *Game_20190206_211253) Up() {
	m.Charset = "utf8"
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE game(`id` int(11) NOT NULL AUTO_INCREMENT,`type` longtext  NOT NULL,`url_list` longtext  NOT NULL,`url` longtext  NOT NULL,`name` longtext  NOT NULL,`final_price` int(11) DEFAULT NULL,`steam_appid` int(11) DEFAULT NULL,`required_age` int(11) DEFAULT NULL,`is_free` tinyint(1) NOT NULL,`dlc` longtext  NOT NULL,`detailed_description` longtext  NOT NULL,`about_the_game` longtext  NOT NULL,`short_description` longtext  NOT NULL,`supported_languages` longtext  NOT NULL,`reviews` longtext  NOT NULL,`header_image` longtext  NOT NULL,`website` longtext  NOT NULL,`view` int(11) DEFAULT NULL,`pc_requirements` longtext  NOT NULL,`mac_requirements` longtext  NOT NULL,`linux_requirements` longtext  NOT NULL,`developers` longtext  NOT NULL,`publishers` longtext  NOT NULL,`demos` longtext  NOT NULL,`price_overview` longtext  NOT NULL,`packages` longtext  NOT NULL,`package_groups` longtext  NOT NULL,`platforms` longtext  NOT NULL,`metacritic` longtext  NOT NULL,`categories` longtext  NOT NULL,`genres` longtext  NOT NULL,`screenshots` longtext  NOT NULL,`movies` longtext  NOT NULL,`recommendations` longtext  NOT NULL,`achievements` longtext  NOT NULL,`release_date` longtext  NOT NULL,`support_info` longtext  NOT NULL,`discount_percent` int(11) DEFAULT NULL,`score` int(11) DEFAULT NULL,`background` longtext  NOT NULL,`content_descriptors` longtext  NOT NULL,PRIMARY KEY (`id`)) CHARACTER SET=utf8;")
}

// Reverse the migrations
func (m *Game_20190206_211253) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `game`")
}
