package generator

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/lambda-platform/lambda/DB"
	agentModels "github.com/lambda-platform/lambda/agent/models"
	krudModels "github.com/lambda-platform/lambda/krud/models"
	"github.com/lambda-platform/lambda/models"
)

func RoleData() {

	Kruds := []krudModels.Krud{}
	DB.DB.Find(&Kruds)

	Roles := []agentModels.Role{}
	DB.DB.Where("id >= ?", 1).Find(&Roles)
	roleData := map[int]map[string]interface{}{}

	for _, Role := range Roles {
		Permissions_ := models.Permissions{}
		json.Unmarshal([]byte(Role.Permissions), &Permissions_)
		Menu := models.VBSchema{}
		DB.DB.Where("id = ?", Permissions_.MenuID).Find(&Menu)
		MenuSchema := new(interface{})
		json.Unmarshal([]byte(Menu.Schema), &MenuSchema)

		roleData[Role.ID] = map[string]interface{}{
			"permissions": Permissions_,
			"menu":        MenuSchema,
			"cruds":       Kruds,
		}
	}

	for k, data := range roleData {

		bolB, _ := json.Marshal(data)
		_ = os.WriteFile("lambda/role_"+strconv.Itoa(k)+".json", bolB, 0700)
	}

}
