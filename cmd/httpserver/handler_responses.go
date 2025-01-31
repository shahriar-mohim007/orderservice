package httpserver

import (
	"net/http"
	utilis "orderservice/utils"
)

var ValidDataNotFound = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "The provided information is invalid. Please recheck and try again.",
	Type:       "error",
}

var InvalidEmailPassword = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "The user credentials were incorrect.",
	Type:       "error",
}

var InternalError = utilis.ResponseState{
	StatusCode: http.StatusInternalServerError,
	Message:    "Internal server error",
	Type:       "error",
}
var UserAlreadyExist = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "User Already Exist With this Email",
	Type:       "error",
}
var UserCreated = utilis.ResponseState{
	StatusCode: http.StatusCreated,
	Message:    "User created successfully",
	Type:       "success",
}

var loginSuccess = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "Login Successful",
	Type:       "success",
}

var Unauthorized = utilis.ResponseState{
	StatusCode: http.StatusUnauthorized,
	Message:    "Unauthorized",
	Type:       "error",
}
var OrderCreateSuccess = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "Order Created Successfully",
	Type:       "success",
}

var BadRequestError = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "Bad Request",
}

var OrderFetchedSuccess = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "Orders successfully fetched.",
	Type:       "success",
}
var OrderCancelSuccess = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "Order Cancelled Successfully",
	Type:       "success",
}

var OrderCancelError = utilis.ResponseState{
	StatusCode: http.StatusBadRequest,
	Message:    "Please contact cx to cancel order",
	Type:       "error",
}

var LogoutSuccess = utilis.ResponseState{
	StatusCode: http.StatusOK,
	Message:    "Successfully logged out",
	Type:       "success",
}
