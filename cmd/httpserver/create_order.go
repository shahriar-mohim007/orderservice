package httpserver

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"net/http"
	"orderservice/repository"
	"orderservice/state"
	utils "orderservice/utils"
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

type OrderResponse struct {
	ConsignmentID   string  `json:"consignment_id"`
	MerchantOrderID string  `json:"merchant_order_id"`
	OrderStatus     string  `json:"order_status"`
	DeliveryFee     float64 `json:"delivery_fee"`
}

type ValidationErrorResponse struct {
	Message string              `json:"message"`
	Type    string              `json:"type"`
	Code    int                 `json:"code"`
	Errors  map[string][]string `json:"errors"`
}

func HandleCreateOrder(app *state.State) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		request := OrderRequestPayload{}
		ctx := req.Context()

		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Invalid JSON",
			})
			_ = ValidDataNotFound.WriteToResponse(w, nil)
			return
		}
		validate := validator.New()
		_ = validate.RegisterValidation("bd_phone", validateBDPhone)
		err = validate.Struct(request)
		if err != nil {
			validationErrors := make(map[string][]string)
			for _, fieldErr := range err.(validator.ValidationErrors) {
				var field string
				var message string

				switch fieldErr.Field() {
				case "StoreID":
					field = "store_id"
					validationErrors[field] = []string{"The store field is required", "Wrong Store selected"}
				case "RecipientName":
					field = "recipient_name"
					message = "The recipient name field is required."
				case "RecipientPhone":
					field = "recipient_phone"
					message = "The recipient phone field is required."
				case "RecipientAddress":
					field = "recipient_address"
					message = "The recipient address field is required."
				case "DeliveryType":
					field = "delivery_type"
					message = "The delivery type field is required."
				case "AmountToCollect":
					field = "amount_to_collect"
					message = "The amount to collect field is required."
				case "ItemQuantity":
					field = "item_quantity"
					message = "The item quantity field is required."
				case "ItemWeight":
					field = "item_weight"
					message = "The item weight field is required."
				case "ItemType":
					field = "item_type"
					message = "The item type field is required."
				default:
					field = fieldErr.Field()
					message = fmt.Sprintf("The %s field is required.", field)
				}

				if field != "store_id" {
					validationErrors[field] = append(validationErrors[field], message)
				}
			}

			response := ValidationErrorResponse{
				Message: "Please fix the given errors",
				Type:    "error",
				Code:    422,
				Errors:  validationErrors,
			}

			w.WriteHeader(http.StatusUnprocessableEntity)
			_ = json.NewEncoder(w).Encode(response)
			return
		}
		userID, _ := GetUserIDFromContext(ctx)
		uuID, err := uuid.FromString(userID)

		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error parsing UUID",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}
		ID, err := uuid.NewV4()
		if err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error generating UUID",
			})
			_ = InternalError.WriteToResponse(w, nil)
		}

		consignmentID := utils.GenerateOrderConsignmentID("DA")
		codFee := request.AmountToCollect * 0.01
		deliveryFee := utils.CalculateDeliveryFee(request.RecipientCity, request.ItemWeight)
		totalFee := deliveryFee + codFee

		order := repository.Order{
			ID:                 ID,
			UserID:             uuID,
			OrderConsignmentID: consignmentID,
			StoreID:            request.StoreID,
			MerchantOrderID:    request.MerchantOrderID,
			RecipientName:      request.RecipientName,
			RecipientPhone:     request.RecipientPhone,
			RecipientAddress:   request.RecipientAddress,
			RecipientCity:      request.RecipientCity,
			RecipientZone:      request.RecipientZone,
			RecipientArea:      request.RecipientArea,
			DeliveryType:       request.DeliveryType,
			ItemType:           request.ItemType,
			SpecialInstruction: request.SpecialInstruction,
			ItemQuantity:       request.ItemQuantity,
			ItemWeight:         request.ItemWeight,
			AmountToCollect:    request.AmountToCollect,
			ItemDescription:    request.ItemDescription,
			TotalFee:           totalFee,
			OrderTypeID:        1,
			CODFee:             codFee,
			PromoDiscount:      0.0,
			Discount:           0.0,
			DeliveryFee:        deliveryFee,
			OrderStatus:        "Pending",
			Archive:            false,
		}
		if err = app.Repository.CreateOrder(ctx, &order); err != nil {
			app.Logger.PrintError(err, map[string]string{
				"context": "Error creating order",
			})
			_ = InternalError.WriteToResponse(w, nil)
			return
		}

		responseData := OrderResponse{
			ConsignmentID:   order.OrderConsignmentID,
			MerchantOrderID: order.MerchantOrderID,
			OrderStatus:     order.OrderStatus,
			DeliveryFee:     order.DeliveryFee,
		}

		_ = OrderCreateSuccess.WriteToResponse(w, responseData)
		return

	}
}
