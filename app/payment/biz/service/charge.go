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
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/payment/biz/consts"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/common/errno"
	"github.com/cloudwego/kitex/pkg/klog"
	"strconv"
	"time"

	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/payment/biz/dal/mysql"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/app/payment/biz/model"
	"github.com/Vigor-Team/youthcamp-2025-mall-be/rpc_gen/kitex_gen/payment"
	"github.com/cloudwego/kitex/pkg/kerrors"
	creditcard "github.com/durango/go-credit-card"
	"github.com/google/uuid"
)

type ChargeService struct {
	ctx context.Context
} // NewChargeService new ChargeService
func NewChargeService(ctx context.Context) *ChargeService {
	return &ChargeService{ctx: ctx}
}

// Run create note info
func (s *ChargeService) Run(req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	card := creditcard.Card{
		Number: req.CreditCard.CreditCardNumber,
		Cvv:    strconv.Itoa(int(req.CreditCard.CreditCardCvv)),
		Month:  strconv.Itoa(int(req.CreditCard.CreditCardExpirationMonth)),
		Year:   strconv.Itoa(int(req.CreditCard.CreditCardExpirationYear)),
	}

	err = card.Validate(true)
	if err != nil {
		return nil, kerrors.NewBizStatusError(consts.ErrCardValidate, err.Error())
	}

	translationId, err := uuid.NewRandom()
	if err != nil {
		klog.CtxErrorf(s.ctx, "uuid.NewRandom.err: %v", err)
		return nil, kerrors.NewBizStatusError(errno.ErrInternal, "uuid.NewRandom error")
	}
	err = model.CreatePaymentLog(mysql.DB, s.ctx, &model.PaymentLog{
		UserId:        req.UserId,
		OrderId:       req.OrderId,
		TransactionId: translationId.String(),
		Amount:        req.Amount,
		PayAt:         time.Now(),
		PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "model.CreatePaymentLog.err: %v", err)
		return nil, kerrors.NewBizStatusError(consts.ErrCreatePaymentLog, "create payment log error")
	}
	return &payment.ChargeResp{TransactionId: translationId.String()}, nil
}
