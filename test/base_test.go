package test

import (
	"fpga-bms-server/router"
	"fpga-bms-server/database"
	"fpga-bms-server/models"
	"flag"
	"fmt"
	"os"
	"testing"

	"fpga-bms-server/config"
	"github.com/gavv/httpexpect"
	"github.com/kataras/iris"
	"github.com/kataras/iris/httptest"
)

var (
	app         *iris.Application // iris.Applications
	testRole    *models.Role
	testPerm    *models.Permission
	testPermIds []uint
	testUser    *models.User
)

//单元测试基境
func TestMain(m *testing.M) {

	// 初始化app
	app = router.Router

	baseCase()

	flag.Parse()
	exitCode := m.Run()

	// 删除测试数据表，保持测试环境
	database.DB.DropTable("users", "roles", "permissions", "role_perms", "oauth_tokens", &models.OauthToken{})

	os.Exit(exitCode)

}

// 单元测试 login 方法
func login(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: config.Conf.Get("app.debug").(bool)})
	if Data != nil {
		e.POST(url).WithJSON(Object).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.POST(url).WithJSON(Object).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

	return
}

// 单元测试 create 方法
func create(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: config.Conf.Get("app.debug").(bool)})
	at := GetLoginToken()

	ob := e.POST(url).WithHeader("Authorization", "Bearer "+at.Token).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()

	ob.Value("status").Equal(Status)
	ob.Value("msg").Equal(Msg)

	for k, v := range Data {
		ob.Value("data").Object().Value(k).Equal(v)
	}

	return
}

// 单元测试 update 方法
func update(t *testing.T, url string, Object interface{}, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: config.Conf.Get("app.debug").(bool)})
	at := GetLoginToken()

	ob := e.PUT(url).WithHeader("Authorization", "Bearer "+at.Token).WithJSON(Object).
		Expect().Status(StatusCode).JSON().Object()

	ob.Value("status").Equal(Status)
	ob.Value("msg").Equal(Msg)

	for k, v := range Data {
		ob.Value("data").Object().Value(k).Equal(v)
	}

	return
}

// 单元测试 getOne 方法
func getOne(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: config.Conf.Get("app.debug").(bool)})
	at := GetLoginToken()
	if Data != nil {
		e.GET(url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.GET(url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

	return
}

// 单元测试 getMore 方法
func getMore(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: config.Conf.Get("app.debug").(bool)})
	at := GetLoginToken()
	if Data != nil {
		e.GET(url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg, Data)
	} else {
		e.GET(url).WithHeader("Authorization", "Bearer "+at.Token).
			Expect().Status(StatusCode).
			JSON().Object().Values().Contains(Status, Msg)
	}

	return
}

// 单元测试 delete 方法
func delete(t *testing.T, url string, StatusCode int, Status bool, Msg string, Data map[string]interface{}) (e *httpexpect.Expect) {
	e = httptest.New(t, app, httptest.Configuration{Debug: config.Conf.Get("app.debug").(bool)})
	at := GetLoginToken()

	e.DELETE(url).WithHeader("Authorization", "Bearer "+at.Token).
		Expect().Status(StatusCode).
		JSON().Object().Values().Contains(Status, Msg)

	return
}

/**
*登陆用户
*@return   Token 返回登陆后的token
 */
func GetLoginToken() models.Token {
	response, status, msg := models.CheckLogin(
		config.Conf.Get("api_account.LoginUserName").(string),
		config.Conf.Get("api_account.LoginPwd").(string),
	)

	// 打印错误信息
	if !status {
		fmt.Println(msg)
	}

	return response
}

func baseCase() {
	perm_json := &models.PermissionJson{
		Name:        "test_update_user",
		Description: "访客",
		DisplayName: "访客",
	}

	testPerm = models.CreatePermission(perm_json)
	testPermIds = []uint{testPerm.ID}

	role_json := &models.RoleJson{
		Name:        "test_update_user",
		Description: "访客",
		DisplayName: "访客",
	}

	testRole = models.CreateRole(role_json, testPermIds)

	aul := &models.UserJson{
		Username: "guest",
		Name:     "访客",
		Password: "guest111",
		RoleID:   testRole.ID,
	}

	testUser = models.CreateUser(aul)
}
