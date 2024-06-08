package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2/google"
	"log"
	"os"

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
		User, err := agentUtils.AuthUserObject(c)

		if err != nil {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":  err.Error(),
				"status": false,
			})
		}

		if config.Config.Database.Connection == "oracle" {
			var status models.NotificationStatusOracle

			DB.DB.Where("NOTIF_ID = ? AND RECEIVER_ID = ?", id, User["id"]).First(&status)

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

			DB.DB.Where("notif_id = ? AND receiver_id = ?", id, User["id"]).First(&status)

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
			var savedToken models.UserFcmTokensOracle

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
			var savedToken models.UserFcmTokens

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

//func Fcm(c *fiber.Ctx) error {
//
//	receivers := []string{"d3hK8PY53VEUhO1sb2m0pr:APA91bGe_ZU_q91sq_AOgntrK_A_Dv-Piv-AesP5r7T2EgoS2m_ID_ifJ1cZrRdJGhXEABNqA3W-4hCNoJ_RoTnuZCdV9wlMfrDPo44CQHMuo8JQjlk5pgAY4YOM0-eHO6meS7WW8F88"}
//
//	msg := models.FCMData{
//		Title: "This is a title. title",
//		Body:  "This is a subtitle. subtitle",
//		Sound: config.LambdaConfig.Notify.Sound,
//		Icon:  "http://localhost/asc/logo.png",
//		Link:  "/p/db4172e3-25ba-807f-1c2b-da6a11d10f3b/d7fb539c-8813-5b66-e893-b4d0b1dd971b/9ac627de-77fe-055f-d347-4bdf63513e90",
//	}
//	notification := models.FCMNotification{
//		Title:       "This is a title. title",
//		Body:        "This is a subtitle. subtitle",
//		Icon:        "http://localhost/asc/logo.png",
//		ClickAction: "http://localhost/control#/p/db4172e3-25ba-807f-1c2b-da6a11d10f3b/d7fb539c-8813-5b66-e893-b4d0b1dd971b/9ac627de-77fe-055f-d347-4bdf63513e90",
//	}
//
//	SendNotification(receivers, msg, notification)
//
//	return c.JSON(map[string]interface{}{
//		"status": true,
//	})
//
//}

func CreateNotification(notification models.NotificationData, data map[string]interface{}) int64 {
	accessToken, err := getAccessToken(config.LambdaConfig.Notify.ServerKey)
	if err != nil {
		log.Fatalf("Error getting access token: %v", err)
	}

	jsonData, _ := json.Marshal(data)
	if config.Config.Database.Connection == "oracle" {
		var Users []agentModels.USERSOracle

		if len(notification.Roles) >= 1 {
			DB.DB.Where("ROLE IN (?)", notification.Roles).Find(&Users)
		} else {
			DB.DB.Where("ID IN (?)", notification.Users).Find(&Users)
		}

		//authUser := agentUtils.AuthUser(c)

		link, ok := data["link"].(string)
		if !ok {
			link = ""
		}

		notificationDB := models.NotificationOracle{
			Link:      link,
			Sender:    1,
			Title:     notification.Notification.Title,
			Body:      notification.Notification.Body,
			Data:      string(jsonData),
			CreatedAt: time.Now(),
		}

		DB.DB.Create(&notificationDB)

		if _, exists := data["first_name"]; !exists {
			data["first_name"] = "Системээс"
		}

		data["created_at"] = notificationDB.CreatedAt
		data["id"] = notificationDB.ID

		for _, User := range Users {
			var fcmToken string
			DB.DB.Model(&models.UserFcmTokens{}).Where("user_id = ?", User.ID).Pluck("fcm_token", &fcmToken)

			if fcmToken != "" {
				SendNotification(accessToken, fcmToken, notification.Notification, data)
			}

			DB.DB.Table("notification_status")
			NotificationStatus := models.NotificationStatus{
				NotifID:    notificationDB.ID,
				ReceiverID: User.ID,
				Seen:       0,
				SeenTime:   time.Now(),
			}

			DB.DB.Create(&NotificationStatus)
		}

		return notificationDB.ID

	} else {
		var Users []agentModels.User

		if len(notification.Roles) >= 1 {
			DB.DB.Where("role IN (?)", notification.Roles).Find(&Users)
		} else {
			DB.DB.Where("id IN (?)", notification.Users).Find(&Users)
		}

		//authUser := agentUtils.AuthUser(c)

		link, ok := data["link"].(string)
		if !ok {
			link = ""
		}

		notificationDB := models.Notification{
			Link:      link,
			Sender:    1,
			Title:     notification.Notification.Title,
			Body:      notification.Notification.Body,
			Data:      string(jsonData),
			CreatedAt: time.Now(),
		}

		DB.DB.Create(&notificationDB)

		if _, exists := data["first_name"]; !exists {
			data["first_name"] = "Системээс"
		}

		data["created_at"] = notificationDB.CreatedAt
		data["id"] = strconv.FormatInt(notificationDB.ID, 10)

		for _, User := range Users {
			var fcmToken string
			DB.DB.Model(&models.UserFcmTokens{}).Where("user_id = ?", User.ID).Pluck("fcm_token", &fcmToken)

			if fcmToken != "" {
				SendNotification(accessToken, fcmToken, notification.Notification, data)
			}

			DB.DB.Table("notification_status")
			NotificationStatus := models.NotificationStatus{
				NotifID:    notificationDB.ID,
				ReceiverID: User.ID,
				Seen:       0,
				SeenTime:   time.Now(),
			}

			DB.DB.Create(&NotificationStatus)
		}

		return notificationDB.ID
	}

}

func getAccessToken(serviceAccountFile string) (string, error) {
	ctx := context.Background()

	creds, err := os.ReadFile(serviceAccountFile)
	if err != nil {
		return "", fmt.Errorf("failed to read service account file: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(creds, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT config: %v", err)
	}

	token, err := conf.TokenSource(ctx).Token()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token.AccessToken, nil
}

func SendNotification(accessToken string, receiver string, notification models.FCMNotification, data map[string]interface{}) {
	payload := models.FCMHTTPRequest{
		Message: models.Message{
			Token:        receiver,
			Notification: notification,
			Data:         data,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	body := bytes.NewReader(payloadBytes)
	url := "https://fcm.googleapis.com/v1/projects/" + config.LambdaConfig.Notify.FirebaseConfig.ProjectID + "/messages:send"

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Fatalf("Error creating new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request to Firebase: %v", err)
		return
	}

	defer resp.Body.Close()

	//// Check the response status code
	//if resp.StatusCode == http.StatusOK {
	//	log.Println("Notification sent successfully.")
	//} else {
	//	log.Printf("Failed to send notification. Status code: %d", resp.StatusCode)
	//}
	//
	//// Read the response body
	//bodyBytes, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	log.Printf("Error reading response body: %v", err)
	//} else {
	//	log.Printf("Response body: %s", string(bodyBytes))
	//}
}
