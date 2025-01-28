package repository

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Order struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	OrderConsignmentID string    `json:"order_consignment_id" db:"order_consignment_id"`
	StoreID            int       `json:"store_id" db:"store_id"`
	MerchantOrderID    string    `json:"merchant_order_id,omitempty" db:"merchant_order_id"`
	RecipientName      string    `json:"recipient_name" db:"recipient_name"`
	RecipientPhone     string    `json:"recipient_phone" db:"recipient_phone"`
	RecipientAddress   string    `json:"recipient_address" db:"recipient_address"`
	RecipientCity      int       `json:"recipient_city" db:"recipient_city"`
	RecipientZone      int       `json:"recipient_zone" db:"recipient_zone"`
	RecipientArea      int       `json:"recipient_area" db:"recipient_area"`
	DeliveryType       int       `json:"delivery_type" db:"delivery_type"`
	ItemType           int       `json:"item_type" db:"item_type"`
	SpecialInstruction string    `json:"special_instruction,omitempty" db:"special_instruction"`
	ItemQuantity       int       `json:"item_quantity" db:"item_quantity"`
	ItemWeight         float64   `json:"item_weight" db:"item_weight"`
	AmountToCollect    float64   `json:"amount_to_collect" db:"amount_to_collect"`
	ItemDescription    string    `json:"item_description,omitempty" db:"item_description"`
	TotalFee           float64   `json:"total_fee" db:"total_fee"`
	OrderTypeID        int       `json:"order_type_id" db:"order_type_id"`
	CODFee             float64   `json:"cod_fee,omitempty" db:"cod_fee"`
	PromoDiscount      float64   `json:"promo_discount,omitempty" db:"promo_discount"`
	Discount           float64   `json:"discount,omitempty" db:"discount"`
	DeliveryFee        float64   `json:"delivery_fee" db:"delivery_fee"`
	OrderStatus        string    `json:"order_status" db:"order_status"`
	Archive            bool      `json:"archive" db:"archive"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
