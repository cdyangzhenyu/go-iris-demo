package models

import (
	"fpga-bms-server/config"
	"fpga-bms-server/database"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Name     string `gorm:"not null VARCHAR(191)"`
	Username string `gorm:"unique;VARCHAR(191)"`
	Password string `gorm:"not null VARCHAR(191)"`
	RoleID   uint
	Role     Role
}

type UserJson struct {
	Username string `json:"username" validate:"required,gte=2,lte=50"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required,gte=2,lte=50"`
	RoleID   uint   `json:"role_id" validate:"required"`
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}       username string [description]
 */
func UserAdminCheckLogin(username string) User {
	u := User{}
	if err := database.DB.Where("username = ?", username).First(&u).Error; err != nil {
		fmt.Printf("UserAdminCheckLoginErr:%s", err)
	}
	return u
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *User [description]
 */
func GetUserById(id uint) *User {
	user := new(User)
	user.ID = id
	if err := database.DB.Preload("Role").Where(user).First(user).Error; err != nil {
		fmt.Printf("GetUserByIdErr:%s", err)
	}

	return user
}

/**
 * 通过 username 获取 user 记录
 * @method GetUserByUserName
 * @param  {[type]}       user  *User [description]
 */
func GetUserByUserName(username string) *User {
	user := &User{Username: username}
	if err := database.DB.Preload("Role").Where(user).First(user).Error; err != nil {
		fmt.Printf("GetUserByUserNameErr:%s", err)
	}
	return user
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func DeleteUserById(id uint) {
	user := GetUserById(id)
	if user.Username == "admin" || user.Username == "system" {
		return
	}

	u := new(User)
	u.ID = id

	if err := database.DB.Delete(u).Error; err != nil {
		fmt.Printf("DeleteUserByIdErr:%s", err)
	}
}

/**
 * 获取所有的账号
 * @method GetAllUser
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllUsers(name, orderBy string, offset, limit int) (users []*User) {
	if err := database.GetAll(name, orderBy, offset, limit).Preload("Role").Where("name != ?", "system").Find(&users).Error; err != nil {
		fmt.Printf("GetAllUserErr:%s", err)
	}
	return
}

/**
 * 创建
 * @method CreateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateUser(aul *UserJson) (user *User) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = new(User)
	user.Username = aul.Username
	user.Password = string(hash)
	user.Name = aul.Name
	user.RoleID = aul.RoleID

	if err := database.DB.Create(user).Error; err != nil {
		fmt.Printf("CreateUserErr:%s", err)
	}

	return
}

/**
 * 更新
 * @method UpdateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateUser(uj *UserJson, id uint) *User {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(uj.Password, salt)

	user := new(User)
	user.ID = id
	uj.Password = string(hash)

	if err := database.DB.Model(user).Updates(uj).Error; err != nil {
		fmt.Printf("UpdateUserErr:%s", err)
	}

	return user
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func CheckLogin(username, password string) (response Token, status bool, msg string) {
	user := UserAdminCheckLogin(username)
	if user.ID == 0 {
		msg = "用户不存在"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password); ok {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			claims["iat"] = time.Now().Unix()
			token.Claims = claims
			tokenString, err := token.SignedString([]byte("secret"))

			if err != nil {
				msg = err.Error()
				return
			}

			oauth_token := new(OauthToken)
			oauth_token.Token = tokenString
			oauth_token.UserId = user.ID
			oauth_token.Secret = "secret"
			oauth_token.Revoked = false
			oauth_token.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauth_token.CreatedAt = time.Now()

			response = oauth_token.OauthTokenCreate()
			status = true
			msg = "登陆成功"

			return

		} else {
			msg = "用户名或密码错误"
			return
		}
	}
}

/**
* 用户退出登陆
* @method UserAdminLogout
* @param  {[type]} ids string [description]
 */
func UserAdminLogout(user_id uint) bool {
	ot := UpdateOauthTokenByUserId(user_id)
	return ot.Revoked
}

/**
*创建系统管理员，分为登录和接口调用
*@param role_id uint
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreateSystemAdmin(role_id uint) {

	account := []string{"api_account", "admin_account"}

	for _, value := range account {
		aul := new(UserJson)
		aul.Username = config.Conf.Get(fmt.Sprintf("%s.%s", value, "LoginUserName")).(string)
		aul.Password = config.Conf.Get(fmt.Sprintf("%s.%s", value, "LoginPwd")).(string)
		aul.Name = config.Conf.Get(fmt.Sprintf("%s.%s", value, "LoginName")).(string)
		aul.RoleID = role_id
		user := GetUserByUserName(aul.Username)
		if user.ID == 0 {
			fmt.Println("创建账号")
			CreateUser(aul)
		} else {
			fmt.Println("重复初始化账号")
		}
	}
    return
 
}
