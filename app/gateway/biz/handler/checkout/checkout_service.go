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

package checkout

import (
	"context"

	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/biz/service"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/biz/utils"
	checkout "github.com/Vigor-Team/youthcamp-2025-mall-be/app/gateway/hertz_gen/gateway/checkout"
	"github.com/cloudwego/hertz/pkg/app"
)

// Checkout .
// @router /checkout [GET]
func Checkout(ctx context.Context, c *app.RequestContext) {
	var err error
	var req checkout.CheckoutReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.FailResponse(ctx, c, err)
		return
	}

	resp, err := service.NewCheckoutService(ctx, c).Run(&req)
	if err != nil {
		utils.FailResponse(ctx, c, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// GetCheckoutPreview .
// @router /api/v1/checkout/preview [GET]
func GetCheckoutPreview(ctx context.Context, c *app.RequestContext) {
	var err error
	var req checkout.GetCheckoutPreviewReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.FailResponse(ctx, c, err)
		return
	}

	resp, err := service.NewGetCheckoutPreviewService(ctx, c).Run(&req)

	if err != nil {
		utils.FailResponse(ctx, c, err)
		return
	}
	utils.SuccessResponse(c, resp)
}

// Charge .
// @router /api/v1/payment/charge [POST]
func Charge(ctx context.Context, c *app.RequestContext) {
	var err error
	var req checkout.ChargeReq
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.FailResponse(ctx, c, err)
		return
	}

	resp, err := service.NewChargeService(ctx, c).Run(&req)

	if err != nil {
		utils.FailResponse(ctx, c, err)
		return
	}
	utils.SuccessResponse(c, resp)
}
