package httpserver

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var bdPhoneRegex = regexp.MustCompile(`^(01)[3-9]{1}[0-9]{8}$`)

func validateBDPhone(fl validator.FieldLevel) bool {
	return bdPhoneRegex.MatchString(fl.Field().String())
}

type OrderRequestPayload struct {
	StoreID            int     `json:"store_id" validate:"required"`
	MerchantOrderID    string  `json:"merchant_order_id"`
	RecipientName      string  `json:"recipient_name" validate:"required"`
	RecipientPhone     string  `json:"recipient_phone" validate:"required,bd_phone"`
	RecipientAddress   string  `json:"recipient_address" validate:"required"`
	RecipientCity      int     `json:"recipient_city" validate:"required"`
	RecipientZone      int     `json:"recipient_zone" validate:"required"`
	RecipientArea      int     `json:"recipient_area" validate:"required"`
	DeliveryType       int     `json:"delivery_type" validate:"required"`
	ItemType           int     `json:"item_type" validate:"required"`
	SpecialInstruction string  `json:"special_instruction"`
	ItemQuantity       int     `json:"item_quantity" validate:"required,min=1"`
	ItemWeight         float64 `json:"item_weight" validate:"required,gt=0"`
	AmountToCollect    float64 `json:"amount_to_collect" validate:"required,gt=0"`
	ItemDescription    string  `json:"item_description"`
}
