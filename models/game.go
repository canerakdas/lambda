package models

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Game struct {
	Id                  int64  `orm:"auto"`
	Type                string `orm:"type(longtext)"`
	Name                string `orm:"type(longtext)"`
	SteamAppid          int
	RequiredAge         int
	FinalPrice          int
	IsFree              bool
	Dlc                 string `orm:"type(longtext)"`
	DetailedDescription string `orm:"type(longtext)"`
	AboutTheGame        string `orm:"type(longtext)"`
	ShortDescription    string `orm:"type(longtext)"`
	SupportedLanguages  string `orm:"type(longtext)"`
	Reviews             string `orm:"type(longtext)"`
	HeaderImage         string `orm:"type(longtext)"`
	Website             string `orm:"type(longtext)"`
	Url                 string `orm:"type(longtext)"`
	UrlList             string `orm:"type(longtext)"`
	View                int
	Score               int
	DiscountPercent     int
	PcRequirements      string `orm:"type(longtext)"`
	MacRequirements     string `orm:"type(longtext)"`
	LinuxRequirements   string `orm:"type(longtext)"`
	Developers          string `orm:"type(longtext)"`
	Publishers          string `orm:"type(longtext)"`
	Demos               string `orm:"type(longtext)"`
	PriceOverview       string `orm:"type(longtext)"`
	Packages            string `orm:"type(longtext)"`
	PackageGroups       string `orm:"type(longtext)"`
	Platforms           string `orm:"type(longtext)"`
	Metacritic          string `orm:"type(longtext)"`
	Categories          string `orm:"type(longtext)"`
	Genres              string `orm:"type(longtext)"`
	Screenshots         string `orm:"type(longtext)"`
	Movies              string `orm:"type(longtext)"`
	Recommendations     string `orm:"type(longtext)"`
	Achievements        string `orm:"type(longtext)"`
	ReleaseDate         string `orm:"type(longtext)"`
	SupportInfo         string `orm:"type(longtext)"`
	Background          string `orm:"type(longtext)"`
	ContentDescriptors  string `orm:"type(longtext)"`
}

func init() {
	orm.RegisterModel(new(Game))
}

// AddGame insert a new Game into database and returns
// last inserted Id on success.
func AddGame(m *Game) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetGameById retrieves Game by Id. Returns error if
// Id doesn't exist
func GetGameById(id int64) (v *Game, err error) {
	o := orm.NewOrm()
	v = &Game{Id: id}
	if err = o.QueryTable(new(Game)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGame retrieves all Game matches certain condition. Returns empty list if
// no records existwindow.location.search
func GetAllGame(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64, priceRange string, count bool) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Game))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(strings.Join([]string{k, "__icontains"}, ""), v)
	}

	// price range
	if priceRange != "" {
		price, _ := strconv.Atoi(priceRange)
		qs = qs.Filter("final_price__lte", price).Exclude("final_price__lte", price-999)
	}

	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Game
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		// count genres
		if count != false {
			c, _ := qs.Count()
			var count []interface{}
			count = append(count, c)
			return count, nil
		}
		return ml, nil
	}
	return nil, err

}

// UpdateGame updates Game by Id and returns error if
// the record to be updated doesn't exist
func UpdateGameById(m *Game) (err error) {
	o := orm.NewOrm()
	v := Game{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		if _, err = o.Update(m); err == nil {
		}
	}
	return
}

// DeleteGame deletes Game by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGame(id int64) (err error) {
	o := orm.NewOrm()
	v := Game{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Game{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
