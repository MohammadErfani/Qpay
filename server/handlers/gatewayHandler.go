package handlers

import (
	"Qpay/models"
	"Qpay/services"
	"Qpay/utils"
	"fmt"
	"github.com/labstack/echo/v4"
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

type PurchaseAddressRequest struct {
	Route string `json:"route"`
}

// ListAllGateways godoc
// @Summary List all gateways for the authenticated user
// @Description Retrieve a list of all gateways associated with the authenticated user.
// @Tags gateways
// @Accept json
// @Produce json
// @Success 200 {array} GatewayResponse "List of gateways"
// @Failure 400 {object} Response "{"status": "error", "message": "You didn't add any gateway. Please register a gateway!"}"
// @Security ApiKeyAuth
// @Router /gateway [get]
func (h *Handler) ListAllGateways(ctx echo.Context) error {
	h.SetUserID(ctx)
	gateways, err := services.GetUserGateways(h.DB, h.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "You didn't Add Any Gateway. Please Register a Gateway!",
		})
	}
	var GatewayResponses []GatewayResponse
	for _, g := range gateways {
		GatewayResponses = append(GatewayResponses, SetGatewayResponse(g))
	}
	return ctx.JSON(http.StatusOK, GatewayResponses)

}

// FindGateway godoc
// @Summary Find a gateway by ID for the authenticated user
// @Description Retrieve details of a specific gateway associated with the authenticated user.
// @Tags gateways
// @Accept json
// @Produce json
// @Param id path int true "Gateway ID"
// @Success 200 {object} GatewayResponse "Gateway details"
// @Failure 400 {object} Response "{"status": "error", "message": "Gateway ID is not correct"}"
// @Failure 404 {object} Response "{"status": "error", "message": "Gateway does not exist!"}"
// @Security ApiKeyAuth
// @Router /gateway/{id} [get]
func (h *Handler) FindGateway(ctx echo.Context) error {
	h.SetUserID(ctx)
	var gateway models.Gateway
	gatewayID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "gateway is not correct",
		})
	}
	gateway, err = services.GetSpecificGateway(h.DB, h.UserID, uint(gatewayID))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Gateway does not exist!",
		})
	}
	return ctx.JSON(http.StatusOK, SetGatewayResponse(gateway))
}

// RegisterNewGateway godoc
// @Summary Register a new gateway for the authenticated user
// @Description Register a new gateway with the provided details.
// @Tags gateways
// @Accept json
// @Produce json
// @Param gatewayRequest body GatewayRequest true "Gateway details"
// @Success 201 {object} Response "Success message"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 403 {object} Response "{"status": "error", "message": "Gateway doesn't match your credential"}"
// @Failure 403 {object} Response "{"status": "error", "message": "You already have a personal gateway"}"
// @Failure 403 {object} Response "{"status": "error", "message": "Commission is incorrect"}"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in gateway"}"
// @Security ApiKeyAuth
// @Router /gateway [post]
func (h *Handler) RegisterNewGateway(ctx echo.Context) error {
	h.SetUserID(ctx)
	var req GatewayRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
	}
	if err := ValidateGateway(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	//if err := ValidateUniqueGateway(gh.DB, &req); err != nil {
	//	return ctx.JSON(http.StatusConflict, err.Error())
	//}
	_, err := services.CreateGateway(h.DB, h.UserID, req.Name, req.Logo, req.BankAccountID, req.CommissionID, req.IsPersonal)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == "UnAuthorize" {
			return ctx.JSON(http.StatusForbidden, Response{
				Status:  "error",
				Message: "gateway doesn't match your credential",
			})
		}
		if err.Error() == "personal error" {
			return ctx.JSON(http.StatusForbidden, Response{
				Status:  "error",
				Message: "you already have personal gateway",
			})
		}
		if err.Error() == "commission error" {
			return ctx.JSON(http.StatusForbidden, Response{
				Status:  "error",
				Message: "commission is incorrect",
			})
		}
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in gateway",
		})
	}
	return ctx.JSON(http.StatusCreated, Response{
		Status:  "success",
		Message: "You're gateway is successfully registered!",
	})
}

// UpdateGateway godoc
// @Summary Update a gateway for the authenticated user
// @Description Update details of a specific gateway associated with the authenticated user.
// @Tags gateways
// @Accept json
// @Produce json
// @Param id path int true "Gateway ID"
// @Param updateRequest body UpdateGatewayRequest true "Updated gateway details"
// @Success 200 {object} Response "Success message"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Invalid gateway details"}"
// @Security ApiKeyAuth
// @Router /gateway/{id} [patch]
func (h *Handler) UpdateGateway(ctx echo.Context) error {
	h.SetUserID(ctx)
	var req UpdateGatewayRequest
	gatewayID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "gateway is not correct",
		})
	}
	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
	}
	_, err = services.UpdateGateway(h.DB, h.UserID, uint(gatewayID), req.Name, req.Logo, req.BankAccountID, req.CommissionID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Your gateway is successfully updated",
	})
}

// PurchaseAddress godoc
// @Summary Purchase an address for a gateway
// @Description Purchase an address for a specific gateway associated with the authenticated user.
// @Tags gateways
// @Accept json
// @Produce json
// @Param id path int true "Gateway ID"
// @Param purchaseAddressRequest body PurchaseAddressRequest true "Purchase address details"
// @Success 200 {object} Response "Success message"
// @Failure 400 {object} Response "{"status": "error", "message": "Bind Error"}"
// @Failure 400 {object} Response "{"status": "error", "message": "Address already in use"}"
// @Security ApiKeyAuth
// @Router /gateway/{id}/address [patch]
func (h *Handler) PurchaseAddress(ctx echo.Context) error {
	h.SetUserID(ctx)
	var req PurchaseAddressRequest
	gatewayID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "gateway is not correct",
		})
	}
	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Bind Error",
		})
	}
	if err = ValidateRoute(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}
	_, err = services.PurchaseAddress(h.DB, h.UserID, uint(gatewayID), req.Route)
	if err != nil {
		if err.Error() == "already in use" {
			return ctx.JSON(http.StatusConflict, Response{
				Status:  "error",
				Message: err.Error(),
			})
		}
		return ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Address Set Successfully",
	})
}

// ListCommissions godoc
// @Summary List all active commissions
// @Description Retrieve a list of all active commissions.
// @Tags gateways
// @Accept json
// @Produce json
// @Success 200 {array} CommissionResponse "List of commissions"
// @Failure 500 {object} Response "{"status": "error", "message": "Internal server error in getting commissions"}"
// @Security ApiKeyAuth
// @Router /gateway/commission/list [get]
func (h *Handler) ListCommissions(ctx echo.Context) error {
	commissions, err := services.ListActiveCommission(h.DB)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Internal server error in getting commissions",
		})
	}
	var commResponses []CommissionResponse
	for _, comm := range commissions {
		commResponses = append(commResponses, SetCommissionsResponse(comm))
	}
	return ctx.JSON(http.StatusOK, commResponses)
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

func ValidateRoute(address *PurchaseAddressRequest) error {
	requiredFields := map[string]string{
		"route": address.Route,
	}
	if err := utils.IsRequired(requiredFields); err != nil {
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
