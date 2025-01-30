package httpserver

import (
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"
	"orderservice/state"
	"strconv"
)

var statusMap = map[string]string{
	"1": "Pending",
	"0": "Cancel",
}
var archiveMap = map[string]bool{
	"0": false,
	"1": true,
}

var ItemTypeMap = map[int]string{
	1: "Other",
	2: "Parcel",
}

var OrderTypeMap = map[int]string{
	1: "Delivery",
}

type GetOrderResponse struct {
	Data        []Order `json:"data"`
	Total       int     `json:"total"`
	CurrentPage int     `json:"current_page"`
	PerPage     int     `json:"per_page"`
	TotalInPage int     `json:"total_in_page"`
	LastPage    int     `json:"last_page"`
}

type Order struct {
	OrderConsignmentID string  `json:"order_consignment_id"`
	OrderCreatedAt     string  `json:"order_created_at"`
	OrderDescription   string  `json:"order_description"`
	MerchantOrderID    string  `json:"merchant_order_id"`
	RecipientName      string  `json:"recipient_name"`
	RecipientAddress   string  `json:"recipient_address"`
	RecipientPhone     string  `json:"recipient_phone"`
	OrderAmount        float64 `json:"order_amount"`
	TotalFee           float64 `json:"total_fee"`
	OrderTypeID        int     `json:"order_type_id"`
	CODFee             float64 `json:"cod_fee"`
	PromoDiscount      float64 `json:"promo_discount"`
	Discount           float64 `json:"discount"`
	DeliveryFee        float64 `json:"delivery_fee"`
	OrderStatus        string  `json:"order_status"`
	OrderType          string  `json:"order_type"`
	ItemType           string  `json:"item_type"`
}

func HandlerGetAllOrders(app *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		userID, _ := GetUserIDFromContext(ctx)
		uuID, err := uuid.FromString(userID)

		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"Context": "Error parsing UUID",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}
		archive := req.URL.Query().Get("archive")
		transferStatus := req.URL.Query().Get("transfer_status")
		limitParam := req.URL.Query().Get("limit")
		pageParam := req.URL.Query().Get("page")

		limit := 10
		page := 0
		if limitParam != "" {
			limit, err = strconv.Atoi(limitParam)
			if err != nil {
				app.Logger.PrintError(fmt.Errorf("invalid limit value"), map[string]string{
					"context": "pagination",
				})
				_ = BadRequestError.WriteToResponse(w, nil)
				return
			}
		}

		if pageParam != "" {
			page, err = strconv.Atoi(pageParam)
			if err != nil {
				app.Logger.PrintError(fmt.Errorf("invalid page value"), map[string]string{
					"context": "pagination",
				})
				_ = BadRequestError.WriteToResponse(w, nil)
				return
			}
		}

		order, err := app.Repository.GetOrders(ctx, statusMap[transferStatus], archiveMap[archive], uuID, limit, page)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"Context": "GetAllOrders",
			})
			_ = InternalError.WriteToResponse(w, nil)
		}

		var orders []Order
		for _, o := range order {
			orders = append(orders, Order{
				OrderConsignmentID: o.OrderConsignmentID,
				OrderCreatedAt:     o.CreatedAt.Format("2006-01-02 15:04:05"),
				OrderDescription:   o.ItemDescription,
				MerchantOrderID:    o.MerchantOrderID,
				RecipientName:      o.RecipientName,
				RecipientAddress:   o.RecipientAddress,
				RecipientPhone:     o.RecipientPhone,
				OrderAmount:        o.AmountToCollect,
				TotalFee:           o.TotalFee,
				OrderTypeID:        o.OrderTypeID,
				CODFee:             o.CODFee,
				PromoDiscount:      o.PromoDiscount,
				Discount:           o.Discount,
				DeliveryFee:        o.DeliveryFee,
				OrderStatus:        o.OrderStatus,
				OrderType:          OrderTypeMap[o.OrderTypeID],
				ItemType:           ItemTypeMap[o.ItemType],
			})
		}

		totalCount, err2 := app.Repository.GetOrderCount(ctx, uuID)
		if err2 != nil {
			app.Logger.PrintError(err, map[string]string{
				"Context": "GetAllOrdersCount",
			})
			_ = InternalError.WriteToResponse(w, nil)
		}

		lastPage := (totalCount + limit - 1) / limit

		response := GetOrderResponse{
			Data:        orders,
			Total:       totalCount,
			CurrentPage: page,
			PerPage:     limit,
			TotalInPage: len(orders),
			LastPage:    lastPage,
		}
		_ = OrderFetchedSuccess.WriteToResponse(w, response)
		return

	}
}
