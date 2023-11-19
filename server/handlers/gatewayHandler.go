package handlers

import (
	"Qpay/database"
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type GatewayRequest struct {
	Name string `json:"name" xml:"name" form:"name" query:"name"`
}

type GatewayResponse struct {
	UserID        uint   `json:"user_id"`
	CommissionID  uint   `json:"commission_id"`
	BankAccountID uint   `json:"bank_account_id"`
	Name          string `json:"name"`
	Logo          string `json:"logo"`
	Route         string `json:"route"`
	Status        string `json:"status"`
	Type          string `json:"type"`
}

func ListAllGateways(ctx echo.Context) error {
	db := database.DB()
	var userID uint = 1
	gateways, err := services.GetUserGateways(db, userID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "You Aren't Add Any Gateway. Please Register a Gateway!")
	}
	var GatewayResponses []GatewayResponse
	for _, ba := range gateways {
		GatewayResponses = append(GatewayResponses, SetGatewayResponse(ba))
	}
	return ctx.JSON(http.StatusOK, GatewayResponses)

}

func FindGateway(ctx echo.Context) error {
	db := database.DB()
	var gateway models.Gateway
	var userID uint = 1
	gatewayID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "bank gateway is not correct")
	}
	gateway, err = services.GetSpecificGateway(db, userID, uint(gatewayID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "You're Gateway is not exist!")
	}
	return ctx.JSON(http.StatusOK, SetGatewayResponse(gateway))
}

func RegisterNewGateway(ctx echo.Context) error {
	db := database.DB()
	var req GatewayRequest
	var userID uint = 1 //user id ro bayad tashkhis bedim o inja vared konim.
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	err := utils.CheckGateway(req.Name)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := ValidateUniqueGateway(db, &req); err != nil {
		return ctx.JSON(http.StatusConflict, err.Error())
	}
	_, err = services.CreateGateway(db, userID, req.Name)
	if err != nil {
		if err.Error() == "UnAuthorize" {
			return ctx.JSON(http.StatusForbidden, "gateway doesn't match your credential")
		}
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in gateway")
	}
	return ctx.JSON(http.StatusOK, "You're gateway is successfully registered!")
}

func ValidateUniqueGateway(db *gorm.DB, gateway *GatewayRequest) error {
	if _, err := services.GetGateway(db, "name", gateway.Name); err == nil {
		return errors.New("name already exist")
	}
	return nil
}

func SetGatewayResponse(gateway models.Gateway) GatewayResponse {
	var status, GatewayType string
	if gateway.Status == models.StatusGatewayActive {
		status = "active"
	} else if gateway.Status == models.StatusGatewayInActive {
		status = "inactive"
	} else if gateway.Status == models.StatusGatewayUnapproved {
		status = "UnApproved"
	} else if gateway.Status == models.StatusGatewayDraft {
		status = "Draft"
	}
	if gateway.Type == models.PersonalTypeGateway {
		GatewayType = "Personal"
	} else if gateway.Type == models.CompanyTypeGateway {
		GatewayType = "Company"
	}
	return GatewayResponse{
		UserID:        gateway.UserID,
		CommissionID:  gateway.CommissionID,
		BankAccountID: gateway.BankAccountID,
		Name:          gateway.Name,
		Logo:          gateway.Logo,
		Route:         gateway.Route,
		Status:        status,
		Type:          GatewayType,
	}
}
