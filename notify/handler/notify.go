package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"

	"github.com/lambda-platform/lambda/DB"
	agentModels "github.com/lambda-platform/lambda/agent/models"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"github.com/lambda-platform/lambda/notify/models"
	"net/http"
	"strconv"
	"time"
	//"io/ioutil"
)

func GetNewNotifications(c *fiber.Ctx) error {
	var unseenCount int64
	user_id := c.Params("user_id")
	DB.DB.Table("notification_status").Where("receiver_id = ? and seen = 0", user_id).Count(&unseenCount)

	if config.Config.SysAdmin.UUID {
		var notifications []models.UserNotifactionsUUID
		DB.DB.Table("notification_status as s").Select("n.*, u.first_name, u.login, s.id as sid, s.seen").Joins("left join notifications as n on n.id = s.notif_id left join users as u on u.id = s.receiver_id").Where("receiver_id = ? and seen = 0", user_id).Order("n.created_at DESC").Limit(30).Find(&notifications)

		return c.JSON(map[string]interface{}{
			"count":         unseenCount,
			"notifications": notifications,
		})
	} else {
		var notifications []models.UserNotifactions
		DB.DB.Table("notification_status as s").Select("n.*, u.first_name, u.login, s.id as sid, s.seen").Joins("left join notifications as n on n.id = s.notif_id left join users as u on u.id = s.receiver_id").Where("receiver_id = ? and seen = 0", user_id).Order("n.created_at DESC").Limit(30).Find(&notifications)

		return c.JSON(map[string]interface{}{
			"count":         unseenCount,
			"notifications": notifications,
		})
	}

}
func GetAllNotifications(c *fiber.Ctx) error {

	user_id := c.Params("user_id")

	if config.Config.SysAdmin.UUID {
		var notifications []models.UserNotifactionsUUID

		DB.DB.Table("notification_status as s").Select("n.*, u.first_name, u.login, s.id as sid, s.seen").Joins("left join notifications as n on n.id = s.notif_id left join users as u on u.id = s.receiver_id").Where("receiver_id = ?", user_id).Order("n.created_at DESC").Find(&notifications)

		return c.JSON(map[string]interface{}{
			"count":         0,
			"notifications": notifications,
		})
	} else {
		var notifications []models.UserNotifactions

		DB.DB.Table("notification_status as s").Select("n.*, u.first_name, u.login, s.id as sid, s.seen").Joins("left join notifications as n on n.id = s.notif_id left join users as u on u.id = s.receiver_id").Where("receiver_id = ?", user_id).Order("n.created_at DESC").Find(&notifications)

		return c.JSON(map[string]interface{}{
			"count":         0,
			"notifications": notifications,
		})
	}

}
func SetSeen(c *fiber.Ctx) error {

	id := c.Params("id")

	if config.Config.SysAdmin.UUID {

		authUser := agentUtils.AuthUserUUID(c)

		var status models.NotificationStatusUUID

		DB.DB.Where("notif_id = ? AND receiver_id = ?", id, authUser.ID).First(&status)

		if status.ID != "" {
			status.Seen = 1
			status.SeenTime = time.Now()
			DB.DB.Save(&status)
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": false,
			})
		}

	} else {
		authUser := agentUtils.AuthUser(c)

		var status models.NotificationStatus

		DB.DB.Where("notif_id = ? AND receiver_id = ?", id, authUser.ID).First(&status)

		if status.ID >= 1 {
			status.Seen = 1
			status.SeenTime = time.Now()
			DB.DB.Save(&status)
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": false,
			})
		}
	}

}
func SetToken(c *fiber.Ctx) error {

	user_id := c.Params("user_id")
	token := c.Params("token")

	if config.Config.SysAdmin.UUID {
		var savedToken models.UserFcmTokensUUID

		DB.DB.Where("user_id = ? AND fcm_token = ?", user_id, token).First(&savedToken)

		if savedToken.ID == "" {
			savedToken.FcmToken = token
			savedToken.UserID = user_id
			DB.DB.Save(&savedToken)
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		}
	} else {
		var savedToken models.UserFcmTokens

		DB.DB.Where("user_id = ? AND fcm_token", user_id, token).First(&savedToken)

		if savedToken.ID == 0 {
			savedToken.FcmToken = token
			intID, _ := strconv.Atoi(user_id)
			savedToken.UserID = intID
			DB.DB.Save(&savedToken)
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		}
	}

}
func SetTokenUrlParam(c *fiber.Ctx) error {

	user_id := c.Query("user")
	token := c.Query("token")
	if config.Config.SysAdmin.UUID {
		var User agentModels.UserUUID

		DB.DB.Where("id = ?", user_id).First(&User)

		if User.ID != "" {
			User.FcmToken = token
			DB.DB.Save(&User)
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": false,
			})
		}

	} else {
		var User agentModels.User

		DB.DB.Where("id = ?", user_id).First(&User)

		if User.ID >= 1 {
			User.FcmToken = token
			DB.DB.Save(&User)
			return c.JSON(map[string]interface{}{
				"status": true,
			})
		} else {
			return c.JSON(map[string]interface{}{
				"status": false,
			})
		}
	}

}
func Fcm(c *fiber.Ctx) error {

	receivers := []string{"d3hK8PY53VEUhO1sb2m0pr:APA91bGe_ZU_q91sq_AOgntrK_A_Dv-Piv-AesP5r7T2EgoS2m_ID_ifJ1cZrRdJGhXEABNqA3W-4hCNoJ_RoTnuZCdV9wlMfrDPo44CQHMuo8JQjlk5pgAY4YOM0-eHO6meS7WW8F88"}

	msg := models.FCMData{
		Title: "This is a title. title",
		Body:  "This is a subtitle. subtitle",
		Sound: "/github.com/lambda-platform/lambda/notification.mp3",
		Icon:  "http://localhost/asc/logo.png",
		Link:  "/p/db4172e3-25ba-807f-1c2b-da6a11d10f3b/d7fb539c-8813-5b66-e893-b4d0b1dd971b/9ac627de-77fe-055f-d347-4bdf63513e90",
	}
	notification := models.FCMNotification{
		Title:       "This is a title. title",
		Body:        "This is a subtitle. subtitle",
		Icon:        "http://localhost/asc/logo.png",
		ClickAction: "http://localhost/control#/p/db4172e3-25ba-807f-1c2b-da6a11d10f3b/d7fb539c-8813-5b66-e893-b4d0b1dd971b/9ac627de-77fe-055f-d347-4bdf63513e90",
	}

	SendNotification(receivers, msg, notification)

	return c.JSON(map[string]interface{}{
		"status": true,
	})

}

func CreateNotification(data models.NotificationData) int64 {
	var tokens []string
	var Users []agentModels.User

	if len(data.Roles) >= 1 {
		DB.DB.Where("role IN (?)", data.Roles).Find(&Users)
	} else {
		DB.DB.Where("id IN (?)", data.Users).Find(&Users)
	}

	for _, User := range Users {
		if User.FcmToken != "" {
			tokens = append(tokens, User.FcmToken)
		}
	}

	//authUser := agentUtils.AuthUser(c)

	notification := models.Notification{
		Link:      data.Data.Link,
		Sender:    1,
		Title:     data.Data.Title,
		Body:      data.Data.Body,
		CreatedAt: time.Now(),
	}

	DB.DB.Create(&notification)

	if data.Data.FirstName == "" {
		data.Data.FirstName = "Системээс"
	}

	data.Data.CreatedAt = notification.CreatedAt
	data.Data.ID = notification.ID
	SendNotification(tokens, data.Data, data.Notification)

	for _, User := range Users {
		DB.DB.Table("notification_status")
		NotificationStatus := models.NotificationStatus{
			NotifID:    notification.ID,
			ReceiverID: User.ID,
			Seen:       0,
			SeenTime:   time.Now(),
		}

		DB.DB.Create(&NotificationStatus)
	}

	return notification.ID

}

func SendNotification(receivers []string, msg interface{}, notification models.FCMNotification) {

	data := models.Payload{
		RegistrationIds: receivers,
		Data:            msg,
		Notification:    notification,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "key="+config.LambdaConfig.Notify.ServerKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		fmt.Println("FIREBASE ERROR")
	} else {
		//fmt.Println("FIREBASE RESPONSE")
		//fmt.Println(resp.Body)
		//
		//bodyBytes, _ := ioutil.ReadAll(resp.Body)
		//
		//fmt.Println(string(bodyBytes))
	}

	defer resp.Body.Close()
}
