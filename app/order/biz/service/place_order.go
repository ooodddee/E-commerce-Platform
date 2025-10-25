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
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/order/biz/consts"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/order/biz/dal/redis"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/common/errno"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"strconv"

	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/order/biz/dal/mysql"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/order/biz/model"
	order "github.com/Vigor-Team/youthcamp-2025-mall-be/rpc_gen/kitex_gen/order"
	"gorm.io/gorm"
)

type PlaceOrderService struct {
	ctx context.Context
} // NewPlaceOrderService new PlaceOrderService
func NewPlaceOrderService(ctx context.Context) *PlaceOrderService {
	return &PlaceOrderService{ctx: ctx}
}

// Run create note info
func (s *PlaceOrderService) Run(req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	// Finish your business logic.
	if len(req.OrderItems) == 0 {
		klog.CtxErrorf(s.ctx, "OrderItems empty")
		return nil, kerrors.NewBizStatusError(errno.ErrGRPCRequestParam, "OrderItems empty")
	}

	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		orderId, err := redis.NextId(s.ctx, "order")
		if err != nil {
			klog.CtxErrorf(s.ctx, "redis.NextId.err: %v", err)
			return kerrors.NewBizStatusError(errno.ErrInternal, "redis.NextId error")
		}

		o := &model.Order{
			OrderId:      strconv.Itoa(int(orderId)),
			OrderState:   model.OrderStatePlaced,
			UserId:       req.UserId,
			UserCurrency: req.UserCurrency,
			Consignee: model.Consignee{
				Email: req.Email,
			},
		}
		if req.Address != nil {
			a := req.Address
			o.Consignee.Country = a.Country
			o.Consignee.State = a.State
			o.Consignee.City = a.City
			o.Consignee.StreetAddress = a.StreetAddress
		}
		if err = tx.Create(o).Error; err != nil {
			klog.CtxErrorf(s.ctx, "tx.Create.err: %v", err)
			return kerrors.NewBizStatusError(errno.ErrInternal, "tx.Create error")
		}

		var itemList []*model.OrderItem
		for _, v := range req.OrderItems {
			itemList = append(itemList, &model.OrderItem{
				OrderIdRefer: o.OrderId,
				ProductId:    v.Item.ProductId,
				Quantity:     v.Item.Quantity,
				Cost:         v.Cost,
			})
		}
		if err = tx.Create(&itemList).Error; err != nil {
			klog.CtxErrorf(s.ctx, "tx.Create.err: %v", err)
			return kerrors.NewBizStatusError(consts.ErrCreateOrder, "tx.Create error")
		}
		resp = &order.PlaceOrderResp{
			Order: &order.OrderResult{
				OrderId: strconv.Itoa(int(orderId)),
			},
		}
		return nil
	})
	return
}
