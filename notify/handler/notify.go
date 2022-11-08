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
)

func GetNewNotifications(c *fiber.Ctx) error {
	if config.Config.Database.Connection == "oracle" {
		var unseenCount int64
		user_id := c.Params("user_id")
		DB.DB.Table("NOTIFICATION_STATUS").Where("RECEIVER_ID = ? and SEEN = 0", user_id).Count(&unseenCount)

		var notifications []models.UserNotificationsOracle
		DB.DB.Table("NOTIFICATION_STATUS").Select("NOTIFICATIONS.*, USERS.FIRST_NAME, USERS.LOGIN, NOTIFICATION_STATUS.ID as SID, NOTIFICATION_STATUS.SEEN").Joins("LEFT JOIN NOTIFICATIONS on NOTIFICATIONS.ID = NOTIFICATION_STATUS.NOTIF_ID LEFT JOIN USERS on USERS.ID = NOTIFICATION_STATUS.RECEIVER_ID").Where("RECEIVER_ID = ? and SEEN = 0", user_id).Order("NOTIFICATIONS.CREATED_AT DESC").Limit(30).Find(&notifications)

		return c.JSON(map[string]interface{}{
			"count":         unseenCount,
			"notifications": notifications,
		})
	} else {
		var unseenCount int64
		user_id := c.Params("user_id")
		DB.DB.Table("notification_status").Where("receiver_id = ? and seen = 0", user_id).Count(&unseenCount)

		if config.Config.SysAdmin.UUID {
			var notifications []models.UserNotificationsUUID
			DB.DB.Table("notification_status as s").Select("n.*, u.first_name, u.login, s.id as sid, s.seen").Joins("left join notifications as n on n.id = s.notif_id left join users as u on u.id = s.receiver_id").Where("receiver_id = ? and seen = 0", user_id).Order("n.created_at DESC").Limit(30).Find(&notifications)

			return c.JSON(map[string]interface{}{
				"count":         unseenCount,
				"notifications": notifications,
			})
		} else {
			var notifications []models.UserNotifications
			DB.DB.Table("notification_status as s").Select("n.*, u.first_name, u.login, s.id as sid, s.seen").Joins("left join notifications as n on n.id = s.notif_id left join users as u on u.id = s.receiver_id").Where("receiver_id = ? and seen = 0", user_id).Order("n.created_at DESC").Limit(30).Find(&notifications)

			return c.JSON(map[string]interface{}{
				"count":         unseenCount,
				"notifications": notifications,
			})
		}
	}

}
func GetAllNotifications(c *fiber.Ctx) error {

	user_id := c.Params("user_id")

	if config.Config.SysAdmin.UUID {
		var notifications []models.UserNotificationsUUID

		DB.DB.Table("notification_status as s").Select("n.*, u.first_name, u.login, s.id as sid, s.seen").Joins("left join notifications as n on n.id = s.notif_id left join users as u on u.id = s.receiver_id").Where("receiver_id = ?", user_id).Order("n.created_at DESC").Find(&notifications)

		return c.JSON(map[string]interface{}{
			"count":         0,
			"notifications": notifications,
		})
	} else {
		if config.Config.Database.Connection == "oracle" {
			var notifications []models.UserNotificationsOracle
			DB.DB.Table("NOTIFICATION_STATUS").Select("NOTIFICATIONS.*, USERS.FIRST_NAME, USERS.LOGIN, NOTIFICATION_STATUS.ID as SID, NOTIFICATION_STATUS.SEEN").Joins("LEFT JOIN NOTIFICATIONS on NOTIFICATIONS.ID = NOTIFICATION_STATUS.NOTIF_ID LEFT JOIN USERS on USERS.ID = NOTIFICATION_STATUS.RECEIVER_ID").Where("RECEIVER_ID = ?", user_id).Order("NOTIFICATIONS.CREATED_AT DESC").Find(&notifications)

			return c.JSON(map[string]interface{}{
				"count":         0,
				"notifications": notifications,
			})
		} else {
			var notifications []models.UserNotifications

			DB.DB.Table("notification_status as s").Select("n.*, u.first_name, u.login, s.id as sid, s.seen").Joins("left join notifications as n on n.id = s.notif_id left join users as u on u.id = s.receiver_id").Where("receiver_id = ?", user_id).Order("n.created_at DESC").Find(&notifications)

			return c.JSON(map[string]interface{}{
				"count":         0,
				"notifications": notifications,
			})
		}

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

		if config.Config.Database.Connection == "oracle" {
			var status models.NotificationStatusOracle

			DB.DB.Where("NOTIF_ID = ? AND RECEIVER_ID = ?", id, authUser.ID).First(&status)

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
		} else {
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
		if config.Config.Database.Connection == "oracle" {
			var savedToken models.UserFcmTokens

			DB.DB.Where("user_id = ? AND fcm_token = ?", user_id, token).Find(&savedToken)

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
		} else {
			var savedToken models.UserFcmTokensOracle

			DB.DB.Where("USER_ID = ? AND FCM_TOKEN = ?", user_id, token).Find(&savedToken)

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

}

func Fcm(c *fiber.Ctx) error {

	receivers := []string{"d3hK8PY53VEUhO1sb2m0pr:APA91bGe_ZU_q91sq_AOgntrK_A_Dv-Piv-AesP5r7T2EgoS2m_ID_ifJ1cZrRdJGhXEABNqA3W-4hCNoJ_RoTnuZCdV9wlMfrDPo44CQHMuo8JQjlk5pgAY4YOM0-eHO6meS7WW8F88"}

	msg := models.FCMData{
		Title: "This is a title. title",
		Body:  "This is a subtitle. subtitle",
		Sound: config.LambdaConfig.Notify.Sound,
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
	if config.Config.Database.Connection == "oracle" {
		var tokens []string
		var Users []agentModels.USERSOracle

		if len(data.Roles) >= 1 {
			DB.DB.Where("ROLE IN (?)", data.Roles).Find(&Users)
		} else {
			DB.DB.Where("ID IN (?)", data.Users).Find(&Users)
		}

		for _, User := range Users {

			var savedTokens []models.UserFcmTokensOracle
			DB.DB.Where("USER_ID = ?", User.ID).Find(&savedTokens)

			for _, savedToken := range savedTokens {
				tokens = append(tokens, savedToken.FcmToken)
			}

		}

		//authUser := agentUtils.AuthUser(c)

		notification := models.NotificationOracle{
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
			NotificationStatus := models.NotificationStatusOracle{
				NotifID:    notification.ID,
				ReceiverID: User.ID,
				Seen:       0,
				SeenTime:   time.Now(),
			}

			DB.DB.Create(&NotificationStatus)
		}

		return notification.ID

	} else {
		var tokens []string
		var Users []agentModels.User

		if len(data.Roles) >= 1 {
			DB.DB.Where("role IN (?)", data.Roles).Find(&Users)
		} else {
			DB.DB.Where("id IN (?)", data.Users).Find(&Users)
		}

		for _, User := range Users {

			var savedTokens []models.UserFcmTokens
			DB.DB.Where("user_id = ?", User.ID).Find(&savedTokens)

			for _, savedToken := range savedTokens {
				tokens = append(tokens, savedToken.FcmToken)
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

		//
		//fmt.Println(string(bodyBytes))
	}

	defer resp.Body.Close()
}
