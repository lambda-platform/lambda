package handlers

import (
	"github.com/gofiber/fiber/v2"
	//"os/exec"
	//"github.com/lambda-platform/lambda/tools"
)

func BuildMe(c *fiber.Ctx) (err error) {

	//AbsolutePath  := utils.AbsolutePath()
	//
	//cmd, err := exec.Command("/bin/sh", "-c",AbsolutePath+"start.sh").Output()
	//if err != nil {
	//	fmt.Printf("error %s", err)
	//}
	//output := string(cmd)
	//
	//fmt.Println(output)
	return c.Redirect("/auth/login")

}
