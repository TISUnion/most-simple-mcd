package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/grpc/api"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	"time"
)

type UserService struct {
}

func (u *UserService) Login(ctx context.Context, req *api.LoginReq) (resp *api.LoginResp, err error) {
	token, e := u._login(req.Account, req.Password)
	err = e
	resp = &api.LoginResp{
		Token: token,
	}
	return
}

func (u *UserService) Logout(ctx context.Context, req *api.LogoutReq) (resp *api.LogoutResp, err error) {
	u._logout()
	resp = new(api.LogoutResp)
	return
}

func (u *UserService) Info(ctx context.Context, req *api.InfoReq) (resp *api.InfoResp, err error) {
	userModel := u._info()
	resp = &api.InfoResp{
		Account:  userModel.Account,
		Password: userModel.Password,
		Nickname: userModel.Nickname,
		Roles:    userModel.Roles,
		Avatar:   userModel.Avatar,
	}
	return
}

func (u *UserService) Update(ctx context.Context, req *api.UpdateReq) (resp *api.UpdateResp, err error) {
	err = u._update(&models.AdminUser{
		Nickname: req.Nickname,
		Account:  req.Account,
		Password: req.Password,
		Roles:    req.Roles,
		Avatar:   req.Avatar,
	})

	resp = new(api.UpdateResp)
	return
}

func (u *UserService) _login(account, passwd string) (token string, err error) {
	var adminObj models.AdminUser
	adminJson := modules.GetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY)
	if adminJson == "" {
		adminObj = *u.setDefaultAccount()
	} else {
		if err = json.Unmarshal([]byte(adminJson), &adminObj); err != nil {
			modules.WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
			return
		}
	}

	if account == adminObj.Account && utils.Md5(passwd) == adminObj.Password {
		token = modules.GetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY)
		if token == "" {
			token = utils.Md5(fmt.Sprintf("%v%s", time.Now().UnixNano(), passwd))
			modules.SetWiteTTLFromDatabase(constant.DEFAULT_TOKEN_DB_KEY, token, constant.DEFAULT_TOKEN_DB_KEY_EXPIRE)
		}
		return
	}
	err = errors.New(constant.PASSWORD_ERROR_MESSAGE)
	return
}

func (u *UserService) _logout() {
	modules.SetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY, "")
}

func (u *UserService) _info() models.AdminUser {
	var adminObj models.AdminUser
	adminJson := modules.GetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY)
	_ = json.Unmarshal([]byte(adminJson), &adminObj)
	adminObj.Password = ""
	return adminObj
}

func (u *UserService) _update(reqInfo *models.AdminUser) error {
	// 获取数据库信息
	var adminObj models.AdminUser
	adminJson := modules.GetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY)
	if err := json.Unmarshal([]byte(adminJson), &adminObj); err != nil {
		modules.WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		return errors.New(constant.HTTP_SYSTEM_ERROR_MESSAGE)
	}
	if reqInfo.Account != "" {
		adminObj.Account = reqInfo.Account
	}
	if reqInfo.Nickname != "" {
		adminObj.Nickname = reqInfo.Nickname
	}
	if reqInfo.Password != "" {
		adminObj.Password = utils.Md5(reqInfo.Password)
	}
	jsonStr, _ := json.Marshal(adminObj)
	modules.SetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY, string(jsonStr))
	// 清空token
	modules.SetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY, "")
	return nil
}

// 设置初始账号密码
func (u *UserService) setDefaultAccount() *models.AdminUser {
	pwd := utils.Md5(constant.DEFAULT_PASSWORD)
	adminObj := &models.AdminUser{
		Nickname: constant.DEFAULT_ACCOUNT,
		Account:  constant.DEFAULT_ACCOUNT,
		Password: pwd,
		Roles:    nil,
	}
	adminJson, _ := json.Marshal(adminObj)
	modules.SetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY, string(adminJson))
	return adminObj
}
