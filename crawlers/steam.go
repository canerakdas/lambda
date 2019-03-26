package crawlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lambda/conf"
	"lambda/mail"
	"lambda/models"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type GameModel struct {
	Success bool `json:"success"`
	Data    models.Game
}

type GameMarshal struct {
	Success bool `json:"success"`
	Data    struct {
		Type                string `json:"type"`
		Name                string `json:"name"`
		SteamAppid          int    `json:"steam_appid"`
		RequiredAge         int    `json:"required_age"`
		IsFree              bool   `json:"is_free"`
		Dlc                 []int  `json:"dlc"`
		DetailedDescription string `json:"detailed_description"`
		AboutTheGame        string `json:"about_the_game"`
		ShortDescription    string `json:"short_description"`
		SupportedLanguages  string `json:"supported_languages"`
		Reviews             string `json:"reviews"`
		HeaderImage         string `json:"header_image"`
		Website             string `json:"website"`
		View                int    `json:"view"`
		PcRequirements      struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"pc_requirements"`
		MacRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"mac_requirements"`
		LinuxRequirements struct {
			Minimum     string `json:"minimum"`
			Recommended string `json:"recommended"`
		} `json:"linux_requirements"`
		Developers []string `json:"developers"`
		Publishers []string `json:"publishers"`
		Demos      []struct {
			Appid       int    `json:"appid"`
			Description string `json:"description"`
		} `json:"demos"`
		PriceOverview struct {
			Currency         string `json:"currency"`
			Initial          int    `json:"initial"`
			Final            int    `json:"final"`
			DiscountPercent  int    `json:"discount_percent"`
			InitialFormatted string `json:"initial_formatted"`
			FinalFormatted   string `json:"final_formatted"`
			Provider         string `json:"provider"`
		} `json:"price_overview"`
		Packages      []int `json:"packages"`
		PackageGroups []struct {
			Name                    string `json:"name"`
			Title                   string `json:"title"`
			Description             string `json:"description"`
			SelectionText           string `json:"selection_text"`
			SaveText                string `json:"save_text"`
			DisplayType             int    `json:"display_type"`
			IsRecurringSubscription string `json:"is_recurring_subscription"`
			Subs                    []struct {
				Packageid                int    `json:"packageid"`
				PercentSavingsText       string `json:"percent_savings_text"`
				PercentSavings           int    `json:"percent_savings"`
				OptionText               string `json:"option_text"`
				OptionDescription        string `json:"option_description"`
				CanGetFreeLicense        string `json:"can_get_free_license"`
				IsFreeLicense            bool   `json:"is_free_license"`
				PriceInCentsWithDiscount int    `json:"price_in_cents_with_discount"`
			} `json:"subs"`
		} `json:"package_groups"`
		Platforms struct {
			Windows bool `json:"windows"`
			Mac     bool `json:"mac"`
			Linux   bool `json:"linux"`
		} `json:"platforms"`
		Metacritic struct {
			Score int    `json:"score"`
			URL   string `json:"url"`
		} `json:"metacritic"`
		Categories []struct {
			ID          int    `json:"id"`
			Description string `json:"description"`
		} `json:"categories"`
		Genres []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Url         string `json:"url"`
		} `json:"genres"`
		Screenshots []struct {
			ID            int    `json:"id"`
			PathThumbnail string `json:"path_thumbnail"`
			PathFull      string `json:"path_full"`
		} `json:"screenshots"`
		Movies []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
			Webm      struct {
				Num480 string `json:"480"`
				Max    string `json:"max"`
			} `json:"webm"`
			Highlight bool `json:"highlight"`
		} `json:"movies"`
		Recommendations struct {
			Total int `json:"total"`
		} `json:"recommendations"`
		Achievements struct {
			Total       int `json:"total"`
			Highlighted []struct {
				Name string `json:"name"`
				Path string `json:"path"`
			} `json:"highlighted"`
		} `json:"achievements"`
		ReleaseDate struct {
			ComingSoon bool   `json:"coming_soon"`
			Date       string `json:"date"`
		} `json:"release_date"`
		SupportInfo struct {
			URL   string `json:"url"`
			Email string `json:"email"`
		} `json:"support_info"`
		Background         string `json:"background"`
		ContentDescriptors struct {
			Ids   []interface{} `json:"ids"`
			Notes interface{}   `json:"notes"`
		} `json:"content_descriptors"`
	} `json:"data"`
}

type Screenshot struct {
	ID            int    `json:"id"`
	PathThumbnail string `json:"path_thumbnail"`
	PathFull      string `json:"path_full"`
}

type PriceOverview struct {
	Currency         string `json:"currency"`
	Initial          int    `json:"initial"`
	Final            int    `json:"final"`
	DiscountPercent  int    `json:"discount_percent"`
	InitialFormatted string `json:"initial_formatted"`
	FinalFormatted   string `json:"final_formatted"`
	Provider         string `json:"provider"`
}

type Metacritic struct {
	Score int    `json:"score"`
	URL   string `json:"url"`
}

type GameList struct {
	Applist struct {
		Apps []struct {
			Appid    int    `json:"appid"`
			Name     string `json:"name"`
			Provider string `json:"provider"`
			Time     string `json:"time"`
			Status   bool   `json:"crawled"`
		} `json:"apps"`
	} `json:"applist"`
}

type Folder struct {
	pwd       string
	static    string
	identify  string
	separator string
	id        string
	ext       string
}

func (i Folder) path() string {
	if i.separator != "" {
		return strings.Join([]string{"", i.static, i.identify, i.separator}, "/")
	} else {
		return strings.Join([]string{"", i.static, i.identify}, "/")
	}
}

func (i Folder) fullPath() string {
	if i.separator != "" {
		return strings.Join([]string{i.pwd, i.static, i.identify, i.separator}, "/")
	} else {
		return strings.Join([]string{i.pwd, i.static, i.identify}, "/")
	}
}

func (i Folder) filePath() string {
	return strings.Join([]string{"", i.static, i.identify, i.separator, strings.Join([]string{i.id, i.ext}, "")}, "/")
}

func (i Folder) fullFilePath() string {
	return strings.Join([]string{i.pwd, i.static, i.identify, i.separator, strings.Join([]string{i.id, i.ext}, "")}, "/")
}

func GetSteamGames(id string) {
	var gameMarshal map[string]GameMarshal
	var gameModel GameModel

	urlPath := "http://store.steampowered.com/api/appdetails?appids={{id}}&cc=us&l=en"
	url := strings.Replace(urlPath, "{{id}}", id, 1)

	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
		os.Exit(1)
	}

	t, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(t, &gameMarshal)

	for _, v := range gameMarshal {
		if v.Data.Type == "game" {
			var Screenshots []Screenshot
			var id = strconv.Itoa(v.Data.SteamAppid)

			gameModel.Data.Id = int64(v.Data.SteamAppid)
			gameModel.Data.Type = v.Data.Type
			gameModel.Data.Name = v.Data.Name
			gameModel.Data.SteamAppid = v.Data.SteamAppid
			gameModel.Data.RequiredAge = v.Data.RequiredAge
			gameModel.Data.IsFree = v.Data.IsFree
			gameModel.Data.DetailedDescription = v.Data.DetailedDescription
			gameModel.Data.AboutTheGame = v.Data.AboutTheGame
			gameModel.Data.ShortDescription = v.Data.ShortDescription
			gameModel.Data.SupportedLanguages = v.Data.SupportedLanguages
			gameModel.Data.Reviews = v.Data.Reviews
			gameModel.Data.HeaderImage = v.Data.HeaderImage
			gameModel.Data.Website = v.Data.Website
			gameModel.Data.Background = v.Data.Background
			gameModel.Data.View = 0

			var keys = map[string]string{
				",":      "_",
				"'":      "",
				":":      "",
				"%27":    "",
				"?":      "",
				"*":      "",
				"&#199;": "o",
				"&#246;": "o",
				"&#214;": "o",
				"&#252;": "u",
				"&#220;": "u",
				"&#231;": "c",
				"&#174;": "�",
				"&amp;":  "-",
				"&nbsp;": "-",
				" ":      "-",
				";":      "-",
				"%20":    "-",
				"/":      "-",
				".":      "",
				"(":      "_",
				")":      "_",
				"<":      "_",
				">":      "_",
				"\"":     "_",
				"\\":     "_",
				"{":      "_",
				"}":      "_",
				"%":      "_",
				"&":      "_",
				"+":      "_",
				"//":     "_",
				"__":     "_",
				"[^\\w]": "-",

				"[":     "-",
				"]":     "-",
				"^":     "-",
				"~":     "-",
				"|":     "-",
				"#":     "-",
				"-----": "-",
				"----":  "-",
				"---":   "-",
				"--":    "-"}
			var formattedName = v.Data.Name
			for key, value := range keys {
				formattedName = strings.Replace(formattedName, key, value, -1)
			}

			gameModel.Data.Url = strings.Join([]string{"/game/", id, "/", strings.ToLower(formattedName), "/"}, "")
			gameModel.Data.UrlList = strings.Join([]string{
				`{"steam":"https://store.steampowered.com/app/`,
				id, `"}`}, "")

			gameModel.Data.Dlc = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v.Data.Dlc)), ","), "[]")

			PcRequirements, _ := json.Marshal(v.Data.PcRequirements)
			gameModel.Data.PcRequirements = string(PcRequirements)

			MacRequirements, _ := json.Marshal(v.Data.MacRequirements)
			gameModel.Data.MacRequirements = string(MacRequirements)

			gameModel.Data.Developers = strings.Trim(strings.Join(v.Data.Developers, ","), "[]")
			gameModel.Data.Publishers = strings.Trim(strings.Join(v.Data.Publishers, ","), "[]")

			v.Data.PriceOverview.Provider = "Steam"
			PriceOverview, _ := json.Marshal(v.Data.PriceOverview)
			gameModel.Data.PriceOverview = strings.Join([]string{`[`, string(PriceOverview), `]`}, "")

			Platforms, _ := json.Marshal(v.Data.Platforms)
			gameModel.Data.Platforms = string(Platforms)

			Metacritic, _ := json.Marshal(v.Data.Metacritic)
			gameModel.Data.Metacritic = string(Metacritic)

			gameModel.Data.Score = ((v.Data.Metacritic.Score * 20) / 100) + ((v.Data.PriceOverview.DiscountPercent * 20) / 100)
			gameModel.Data.DiscountPercent = v.Data.PriceOverview.DiscountPercent
			gameModel.Data.FinalPrice = v.Data.PriceOverview.Final

			for i, k := range v.Data.Genres {
				formattedDescription := strings.Replace(k.Description, " ", "_", -1)
				v.Data.Genres[i].Url = strings.Join([]string{"/catalog/", formattedDescription, "/"}, "")
			}

			Genres, _ := json.Marshal(v.Data.Genres)
			gameModel.Data.Genres = string(Genres)

			ScreenshotJson, _ := json.Marshal(v.Data.Screenshots)
			gameModel.Data.Screenshots = string(ScreenshotJson)

			Movies, _ := json.Marshal(v.Data.Movies)
			gameModel.Data.Movies = string(Movies)

			Recommendations, _ := json.Marshal(v.Data.Recommendations)
			gameModel.Data.Recommendations = string(Recommendations)

			Achievements, _ := json.Marshal(v.Data.Achievements)
			gameModel.Data.Achievements = string(Achievements)

			ReleaseDate, _ := json.Marshal(v.Data.ReleaseDate)
			gameModel.Data.ReleaseDate = string(ReleaseDate)

			SupportInfo, _ := json.Marshal(v.Data.SupportInfo)
			gameModel.Data.SupportInfo = string(SupportInfo)

			ContentDescriptors, _ := json.Marshal(v.Data.ContentDescriptors)
			gameModel.Data.ContentDescriptors = string(ContentDescriptors)

			LinuxRequirements, _ := json.Marshal(v.Data.LinuxRequirements)
			gameModel.Data.LinuxRequirements = string(LinuxRequirements)

			Demos, _ := json.Marshal(v.Data.Demos)
			gameModel.Data.Demos = string(Demos)

			Packages, _ := json.Marshal(v.Data.Packages)
			gameModel.Data.Packages = string(Packages)

			PackageGroups, _ := json.Marshal(v.Data.PackageGroups)
			gameModel.Data.PackageGroups = string(PackageGroups)

			Categories, _ := json.Marshal(v.Data.Categories)
			gameModel.Data.Categories = string(Categories)

			_, err := models.AddGame(&gameModel.Data)
			// Save Header Image
			pwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Fatal(err)
			}

			var staticDirectory = Folder{pwd, "static/image", id, "", id, ""}
			var headerDirectory = Folder{pwd, "static/image", id, "header", id, ".jpg"}
			var backgroundDirectory = Folder{pwd, "static/image", id, "background", id, ".jpg"}
			var screenshotsDirectory = Folder{pwd, "static/image", id, "screenshots", id, ".jpg"}
			os.Mkdir(staticDirectory.fullPath(), os.ModePerm)
			os.Mkdir(headerDirectory.fullPath(), os.ModePerm)
			os.Mkdir(backgroundDirectory.fullPath(), os.ModePerm)
			os.Mkdir(screenshotsDirectory.fullPath(), os.ModePerm)


			Image(v.Data.HeaderImage, headerDirectory.fullFilePath())
			Image(v.Data.Background, backgroundDirectory.fullFilePath())

			// Save Screenshots
			for i, k := range v.Data.Screenshots {
				screenshotsDirectory.id = strconv.Itoa(i)
				Image(k.PathThumbnail, screenshotsDirectory.fullFilePath())

				Screenshots = append(Screenshots, Screenshot{
					PathThumbnail: screenshotsDirectory.filePath(),
					PathFull:      screenshotsDirectory.filePath(),
					ID:            i,
				})
			}

			if err != nil {
				var log models.Log
				log.Description = "Game cant added"
				log.ExternalId = id
				log.Time = time.Now().String()
				models.AddLog(&log)
			} else {
				ScreenshotJson, _ := json.Marshal(Screenshots)
				gameModel.Data.Screenshots = string(ScreenshotJson)

				gameModel.Data.HeaderImage = headerDirectory.filePath()
				gameModel.Data.Background = backgroundDirectory.filePath()

				err = models.UpdateGameById(&gameModel.Data)
				if err != nil {
					var log models.Log
					log.Description = "Game added but cant updated"
					log.ExternalId = id
					log.Time = time.Now().String()
					models.AddLog(&log)
				} else {
					var base64Id int64
					base64Id, err = strconv.ParseInt(id, 10, 64)
					if err == nil{
						list, er := models.GetListById(base64Id)
						if er == nil{
							list.Status = true
							list.Time = time.Now().Unix()
							models.UpdateListById(list)
						}else{
							fmt.Println(er)
						}

					}
				}
			}
		} else {
			base64Id, _ := strconv.ParseInt(id, 10, 64)
			list, _ := models.GetListById(base64Id)
			list.Status = true
			list.Time = time.Now().Unix()
			models.UpdateListById(list)
		}
	}

	if err == nil {
	} else {
		fmt.Printf("server not responding %s", err.Error())
		os.Exit(1)
	}

	defer r.Body.Close()
}

func GetSteamGame(id string) (game map[string]GameMarshal) {
	urlPath := "http://store.steampowered.com/api/appdetails?appids={{id}}&cc=us&l=en"
	url := strings.Replace(urlPath, "{{id}}", id, 1)
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
		os.Exit(1)
	}

	t, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(t, &game)
	defer r.Body.Close()
	return game
}

func UpdateSteamGame(game *models.Game, total int) {
	steamGameInformation := GetSteamGame(strconv.Itoa(game.SteamAppid))
	var initialPrice PriceOverview
	var priceIndex int

	for _, steamRaw := range steamGameInformation {
		// New game data
		var steamPrice = steamRaw.Data.PriceOverview
		// Old game data
		var price []PriceOverview
		json.Unmarshal([]byte(game.PriceOverview), &price)

		for i, v := range price {
			if v.Provider == "Steam" {
				initialPrice = v
				priceIndex = i
			}
		}

		if initialPrice.Final != 0 && steamPrice.Final != 0 {
			if initialPrice.Final != steamPrice.Final {
				steamPrice.Provider = "Steam"
				price[priceIndex] = steamPrice

				PriceOverview, _ := json.Marshal(price)
				game.PriceOverview = strings.Join([]string{string(PriceOverview)}, "")

				var meta Metacritic
				json.Unmarshal([]byte(game.Metacritic), meta)

				var appId = strconv.Itoa(game.SteamAppid)
				var q = map[string]string{"notifications": strings.Join([]string{`\"`, appId, `\"`}, "")}
				var r, _ = models.GetAllUser(q, []string{}, []string{}, []string{}, 0, 0)

				for i := range r {
					notificationsRaw := []byte(r[i].(models.User).Notifications)
					var notifications map[string]interface{}
					json.Unmarshal(notificationsRaw, &notifications)
					if int(notifications[appId].(float64)*10) < steamPrice.DiscountPercent {
						var email = r[i].(models.User).Email

						templateData := struct {
							Name        string
							URL         string
							Discount    int
							Description string
							Image       string
						}{
							Name:        game.Name,
							URL:         game.Url,
							Discount:    game.DiscountPercent,
							Description: game.ShortDescription,
							Image:       game.HeaderImage,
						}

						var auth = smtp.PlainAuth("", static.Sender.Email, static.Sender.Password, static.Sender.Host)

						var subject = strings.Join([]string{game.Name, " down to the ", strconv.Itoa(steamPrice.DiscountPercent)}, "")
						r := mail.NewRequest([]string{email}, subject, "Lambd", auth)
						err := r.ParseTemplate("./mail/template/discount.html", templateData)
						if err := r.ParseTemplate("./mail/template/discount.html", templateData); err == nil {
							ok, _ := r.SendEmail()
							fmt.Println(ok)
						}
						if err != nil {
							fmt.Println(err)
						}
					}
				}

				game.Score = ((meta.Score * 20) / 100) + ((steamPrice.DiscountPercent * 20) / 100)
				game.DiscountPercent = steamPrice.DiscountPercent
				game.FinalPrice = steamPrice.Final
				game.View = total
				models.UpdateGameById(game)
			} else {
				game.View = total
				models.UpdateGameById(game)
			}
		} else {
			game.View = total
			models.UpdateGameById(game)
		}
	}
}

func UpdateSteamGameList() {
	var gameList GameList
	var gameModel models.List
	url := "http://api.steampowered.com/ISteamApps/GetAppList/v0002/?key=0DF75B7B51BAF08BB992D5327748BFB8&format=json"
	r, err := http.Get(url)
	if err != nil {
		fmt.Printf("server not responding %s", err.Error())
		os.Exit(1)
	}

	t, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(t, &gameList)

	for _, v := range gameList.Applist.Apps {
		gameModel.Id = int64(v.Appid)
		gameModel.Appid = v.Appid
		gameModel.Provider = "Steam"
		gameModel.Time = time.Now().Unix()
		gameModel.Status = false
		gameModel.Name = v.Name
		_, err := models.GetListById(gameModel.Id)
		if err != nil {
			models.AddList(&gameModel)
		}
	}

	fmt.Println("Game list updated")

	defer r.Body.Close()
}

func SyncSteamGameList() {
	list, _ := models.GetAllList(map[string]string{"status": "false"}, []string{}, []string{}, []string{}, 0, 1000)
	fmt.Println("Game list syncing with by Steam")

	for _, v := range list {
		values := v.(models.List)
		GetSteamGames(strconv.Itoa(values.Appid))
	}
	if len(list) == 0 {
		fmt.Println("No active games were found in the game list, the list is updating.")
		UpdateSteamGameList()
	}
}
