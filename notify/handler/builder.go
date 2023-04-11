package handler

import (
	"bytes"
	"fmt"
	"github.com/PaesslerAG/gval"
	"github.com/lambda-platform/lambda/DB"
	models2 "github.com/lambda-platform/lambda/agent/models"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/notify/models"
	"regexp"
	"strings"
	"text/template"
)

func BuildNotification(dataJson map[string]interface{}, schemaId int64, action string, userPre map[string]interface{}) {

	if config.Config.Database.Connection == "oracle" {
		target := models.NotificationTargetOracle{}

		DB.DB.Where("SCHEMA_ID = ? AND TARGET_ACTIONS LIKE ?", schemaId, "%"+action+"%").Find(&target)

		if target.ID >= 1 {

			user := models2.USERSOracle{}

			DB.DB.Where("ID = ?", userPre["id"]).Find(&user)

			var re1 = regexp.MustCompile(`{`)
			template := re1.ReplaceAllString(target.Condition, ``)
			var re2 = regexp.MustCompile(`}`)
			template = re2.ReplaceAllString(template, ``)
			var re3 = regexp.MustCompile(`'`)
			template = re3.ReplaceAllString(template, `"`)

			value, _ := gval.Evaluate(template, dataJson)

			Body := Execute(dataJson, target.Body)

			if value == true {

				FCMData := models.FCMData{
					Title:       target.Title,
					Body:        Body,
					FirstName:   user.FirstName,
					Sound:       config.LambdaConfig.Notify.Sound,
					Icon:        config.LambdaConfig.Favicon,
					Link:        target.Link,
					ClickAction: config.LambdaConfig.Domain + "/" + target.Link,
				}

				FCMNotification := models.FCMNotification{
					Title:       target.Title,
					Body:        Body,
					Icon:        config.LambdaConfig.Domain + "/" + config.LambdaConfig.Favicon,
					ClickAction: config.LambdaConfig.Domain + "/" + target.Link,
				}

				data := models.NotificationData{
					Roles:        []int{target.TargetRole},
					Data:         FCMData,
					Notification: FCMNotification,
				}
				CreateNotification(data, dataJson)

			}
		}
	} else {
		target := models.NotificationTarget{}

		DB.DB.Where("schema_id = ? AND target_actions LIKE ?", schemaId, "%"+action+"%").Find(&target)

		if target.ID >= 1 {

			user := models2.User{}

			DB.DB.Where("id = ?", userPre["id"]).Find(&user)

			var re1 = regexp.MustCompile(`{`)
			template := re1.ReplaceAllString(target.Condition, ``)
			var re2 = regexp.MustCompile(`}`)
			template = re2.ReplaceAllString(template, ``)
			var re3 = regexp.MustCompile(`'`)
			template = re3.ReplaceAllString(template, `"`)

			value, _ := gval.Evaluate(template, dataJson)

			fmt.Println(template)
			fmt.Println(dataJson)

			Body := Execute(dataJson, target.Body)

			if value == true {

				FCMData := models.FCMData{
					Title:       target.Title,
					Body:        Body,
					FirstName:   user.FirstName,
					Sound:       config.LambdaConfig.Notify.Sound,
					Icon:        config.LambdaConfig.Favicon,
					Link:        target.Link,
					ClickAction: config.LambdaConfig.Domain + "/" + target.Link,
				}

				FCMNotification := models.FCMNotification{
					Title:       target.Title,
					Body:        Body,
					Icon:        config.LambdaConfig.Domain + "/" + config.LambdaConfig.Favicon,
					ClickAction: config.LambdaConfig.Domain + "/" + target.Link,
				}

				data := models.NotificationData{
					Roles:        []int{target.TargetRole},
					Data:         FCMData,
					Notification: FCMNotification,
				}
				CreateNotification(data, dataJson)

			}
		}
	}

}

func Execute(data interface{}, TBody string) string {
	TBody = strings.Replace(TBody, "{", "{{.", -1)
	TBody = strings.Replace(TBody, "}", "}}", -1)
	t := template.Must(template.New("").Parse(TBody))
	buf := bytes.Buffer{}
	t.Execute(&buf, data)
	return buf.String()
}
