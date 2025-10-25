// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"

	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/biz/consts"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/biz/dal/mysql"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/biz/model"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/pkg/kerrors"

	auth "github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/hertz_gen/gateway/auth"
	common "github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/hertz_gen/gateway/common"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/infra/rpc"
	rpcuser "github.com/Vigor-Team/youthcamp-2025-mall-be/rpc_gen/kitex_gen/user"
	"github.com/cloudwego/hertz/pkg/app"
)

type RegisterService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRegisterService(Context context.Context, RequestContext *app.RequestContext) *RegisterService {
	return &RegisterService{RequestContext: RequestContext, Context: Context}
}

func (h *RegisterService) Run(req *auth.RegisterReq) (resp *common.Empty, err error) {
	userID, err := rpc.UserClient.Register(h.Context, &rpcuser.RegisterReq{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.Password,
	})
	if err != nil {
		return nil, err
	}
	err = func() error {
		role, err := model.GetRoleByName(mysql.DB, h.Context, "Registered User")
		if err != nil {
			hlog.CtxErrorf(h.Context, "GetRoleByName failed, err: %v", err)
			return err
		}
		err = model.BindUserRole(mysql.DB, h.Context, &model.UserRole{
			UID: int64(userID.UserId),
			RID: role.ID,
		})
		return err
	}()
	if err != nil {
		hlog.CtxErrorf(h.Context, "BindUserRole failed, err: %v", err)
		_, err := rpc.UserClient.DeleteUser(h.Context, &rpcuser.UserDeleteReq{
			UserId: userID.UserId,
		})
		if err != nil {
			hlog.CtxErrorf(h.Context, "DeleteUser failed, err: %v", err)
			return nil, err
		}
		return nil, kerrors.NewBizStatusError(consts.ErrBindRoleUser, "BindUserRole failed")
	}
	return
}
