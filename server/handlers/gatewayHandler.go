package handlers

import (
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type GatewayRequest struct {
	Name          string `json:"name" xml:"name" form:"name" query:"name"`
	Logo          string `json:"logo" xml:"logo" form:"logo" query:"logo"`
	BankAccountID uint   `json:"bank_account_id" xml:"bank_account_id" form:"bank_account_id" query:"bank_account_id"`
	CommissionID  uint   `json:"commission_id" xml:"commission_id" form:"commission_id" query:"commission_id"`
	IsPersonal    bool   `json:"is_personal"`
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

type GatewayHandler struct {
	DB     *gorm.DB
	UserID uint
}
type ChangeGatewayRequest struct {
	BankAccountID uint `json:"bank_account_id"`
	GatewayID     uint `json:"gateway_id"`
}

type UpdateGatewayRequest struct {
	CommissionID  uint   `json:"commission_id"`
	BankAccountID uint   `json:"bank_account_id"`
	Name          string `json:"name"`
	Logo          string `json:"logo"`
}

func (gh *GatewayHandler) ListAllGateways(ctx echo.Context) error {
	gateways, err := services.GetUserGateways(gh.DB, gh.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "You didn't Add Any Gateway. Please Register a Gateway!")
	}
	var GatewayResponses []GatewayResponse
	for _, ba := range gateways {
		GatewayResponses = append(GatewayResponses, SetGatewayResponse(ba))
	}
	return ctx.JSON(http.StatusOK, GatewayResponses)

}

func (gh *GatewayHandler) FindGateway(ctx echo.Context) error {
	var gateway models.Gateway
	gatewayID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "gateway is not correct")
	}
	gateway, err = services.GetSpecificGateway(gh.DB, gh.UserID, uint(gatewayID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Gateway does not exist!")
	}
	return ctx.JSON(http.StatusOK, SetGatewayResponse(gateway))
}

func (gh *GatewayHandler) RegisterNewGateway(ctx echo.Context) error {
	var req GatewayRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	if err := ValidateGateway(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	//if err := ValidateUniqueGateway(gh.DB, &req); err != nil {
	//	return ctx.JSON(http.StatusConflict, err.Error())
	//}
	_, err := services.CreateGateway(gh.DB, gh.UserID, req.Name, req.Logo, req.BankAccountID, req.CommissionID, req.IsPersonal)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == "UnAuthorize" {
			return ctx.JSON(http.StatusForbidden, "gateway doesn't match your credential")
		}
		if err.Error() == "personal error" {
			return ctx.JSON(http.StatusForbidden, "you already have personal gateway")
		}
		if err.Error() == "commission error" {
			return ctx.JSON(http.StatusForbidden, "commission is incorrect")
		}
		return ctx.JSON(http.StatusInternalServerError, "Internal server error in gateway")
	}
	return ctx.JSON(http.StatusCreated, "You're gateway is successfully registered!")
}

func (gh *GatewayHandler) UpdateGateway(ctx echo.Context) error {
	var req UpdateGatewayRequest
	gatewayID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "gateway is not correct")
	}
	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	_, err = services.UpdateGateway(gh.DB, gh.UserID, uint(gatewayID), req.Name, req.Logo, req.BankAccountID, req.CommissionID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, "You're gateway is successfully updated")
}

func (gh *GatewayHandler) DeleteGateway(ctx echo.Context) error {
	// Todo
	return nil
}

//	func ValidateUniqueGateway(db *gorm.DB, gateway *GatewayRequest) error {
//		if _, err := services.GetGateway(db, "name", gateway.Name); err == nil {
//			return errors.New("name already exist")
//		}
//		if _, err := services.GetGateway(db, "route", gateway.Route); err == nil {
//			return errors.New("route already exist")
//		}
//		return nil
//	}
func ValidateGateway(gateway *GatewayRequest) error {

	requiredFields := map[string]string{
		"name": gateway.Name,
	}
	requiredIDFields := map[string]uint{
		"bank_account_id": gateway.BankAccountID,
		"commission_id":   gateway.CommissionID,
	}
	if err := utils.IsRequired(requiredFields); err != nil {
		return err
	}
	if err := utils.IsRequiredID(requiredIDFields); err != nil {
		return err
	}
	if err := utils.CheckGateway(gateway.Name); err != nil {
		return err
	}

	return nil
}
func SetGatewayResponse(gateway models.Gateway) GatewayResponse {
	var status, GatewayType string
	if gateway.Status == models.StatusGatewayActive {
		status = "active"
	} else if gateway.Status == models.StatusGatewayInActive {
		status = "inactive"
	} else if gateway.Status == models.StatusGatewayBlocked {
		status = "blocked"
	} else if gateway.Status == models.StatusGatewayDraft {
		status = "Draft"
	}
	if gateway.Type == models.PersonalTypeGateway {
		GatewayType = "Personal"
	} else if gateway.Type == models.BusinessTypeGateway {
		GatewayType = "Business"
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
func (gh *GatewayHandler) ChangeGetawayBankAccount(ctx echo.Context) error {
	var req ChangeGatewayRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Bind Error")
	}
	gateway, err := services.GetSpecificGateway(gh.DB, gh.UserID, req.GatewayID)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, "gateway doesn't exist for this user")
	}
	_, err = services.GetSpecificBankAccount(gh.DB, gh.UserID, req.BankAccountID)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, "bank account doesn't exist for this user")
	}
	gateway.BankAccountID = req.BankAccountID
	return ctx.JSON(http.StatusForbidden, "Bank account updated successfully")

}
