package routes

import (
	"main/lib"
	"main/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AccountHandler struct {
	DB *gorm.DB
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{DB: db}
}

func (h *AccountHandler) GetAccounts(c echo.Context) error {
	var accounts []models.Account

	if err := h.DB.Preload("Player").Preload("Economy").Preload("Setting").Find(&accounts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	var response []models.AccountResponse
	for _, account := range accounts {
		var accountResponse models.AccountResponse
		lib.MapStruct(&account, &accountResponse)
		lib.MapStruct(&account.Player, &accountResponse.Player)
		lib.MapStruct(&account.Economy, &accountResponse.Economy)
		lib.MapStruct(&account.Setting, &accountResponse.Setting)
		response = append(response, accountResponse)
	}

	return c.JSON(http.StatusOK, response)
}

func (h *AccountHandler) GetAccount(c echo.Context) error {
	uuid := c.Param("accountId")

	var account models.Account
	if err := h.DB.Preload("Player").Preload("Economy").Preload("Setting").Where("uuid = ?", uuid).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, models.Error{Reason: "account not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	response := models.AccountResponse{}
	lib.MapStruct(&account, &response)
	lib.MapStruct(&account.Player, &response.Player)
	lib.MapStruct(&account.Economy, &response.Economy)
	lib.MapStruct(&account.Setting, &response.Setting)

	return c.JSON(http.StatusOK, response)
}

func (h *AccountHandler) CreateAccount(c echo.Context) error {
	uuid := c.Param("accountId")
	req := new(models.AccountResponse)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	var existingAccount models.Account
	if err := h.DB.Where("uuid = ?", uuid).First(&existingAccount).Error; err == nil {
		return c.JSON(http.StatusConflict, models.Error{Reason: "Account with this UUID already exists"})
	}

	player := models.Player{}
	economy := models.Economy{}
	setting := models.Setting{}

	lib.MapStruct(&req.Player, &player)
	lib.MapStruct(&req.Economy, &economy)
	lib.MapStruct(&req.Setting, &setting)

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&player).Error; err != nil {
			return err
		}
		if err := tx.Create(&economy).Error; err != nil {
			return err
		}
		if err := tx.Create(&setting).Error; err != nil {
			return err
		}

		account := models.Account{
			UUID:      uuid,
			PlayerID:  player.ID,
			EconomyID: economy.ID,
			SettingID: setting.ID,
		}

		if err := tx.Create(&account).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Account created successfully"})
}

func (h *AccountHandler) UpdateAccount(c echo.Context) error {
	uuid := c.Param("accountId")
	req := new(models.AccountResponse)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	var account models.Account
	if err := h.DB.Where("uuid = ?", uuid).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, models.Error{Reason: "account not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	// Ensure IDs are set correctly before updating
	req.Player.ID = account.PlayerID
	req.Economy.ID = account.EconomyID
	req.Setting.ID = account.SettingID

	if err := h.DB.Model(&models.Player{ID: account.PlayerID}).Updates(req.Player).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	if err := h.DB.Model(&models.Economy{ID: account.EconomyID}).Updates(req.Economy).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	if err := h.DB.Model(&models.Setting{ID: account.SettingID}).Updates(req.Setting).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Account updated successfully"})
}

func (h *AccountHandler) GetHistory(c echo.Context) error {
	return c.JSON(http.StatusOK, &models.APIHelp{
		Version:    "v2.0.0",
		WhatIsThis: "API for this server, don't ask me how to use it... :(",
		Author:     "https://github.com/FantaisieMinecraftServer",
		HomePage:   "https://www.tensyoserver.net",
	})
}

func (h *AccountHandler) CreateHistory(c echo.Context) error {
	return c.JSON(http.StatusOK, &models.APIHelp{
		Version:    "v2.0.0",
		WhatIsThis: "API for this server, don't ask me how to use it... :(",
		Author:     "https://github.com/FantaisieMinecraftServer",
		HomePage:   "https://www.tensyoserver.net",
	})
}

func (h *AccountHandler) UpdateHistory(c echo.Context) error {
	return c.JSON(http.StatusOK, &models.APIHelp{
		Version:    "v2.0.0",
		WhatIsThis: "API for this server, don't ask me how to use it... :(",
		Author:     "https://github.com/FantaisieMinecraftServer",
		HomePage:   "https://www.tensyoserver.net",
	})
}

func (h *AccountHandler) Deposit(c echo.Context) error {
	uuid := c.Param("accountId")
	typeParam := c.QueryParam("type")
	amountParam := c.QueryParam("amount")

	amount, err := strconv.ParseFloat(amountParam, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{Reason: "invalid amount"})
	}

	var account models.Account
	if err := h.DB.Preload("Economy").Where("uuid = ?", uuid).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, models.Error{Reason: "account not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	switch typeParam {
	case "cash":
		account.Economy.Cash += amount
	case "vault":
		account.Economy.Vault += amount
	case "bank":
		account.Economy.Bank += amount
	case "crypto":
		account.Economy.Crypto += amount
	default:
		return c.JSON(http.StatusBadRequest, models.Error{Reason: "invalid type"})
	}

	if err := h.DB.Save(&account.Economy).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "deposit successful"})
}

func (h *AccountHandler) Withdraw(c echo.Context) error {
	uuid := c.Param("accountId")
	typeParam := c.QueryParam("type")
	amountParam := c.QueryParam("amount")

	amount, err := strconv.ParseFloat(amountParam, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{Reason: "invalid amount"})
	}

	var account models.Account
	if err := h.DB.Preload("Economy").Where("uuid = ?", uuid).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, models.Error{Reason: "account not found"})
		}

		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	switch typeParam {
	case "cash":
		if account.Economy.Cash < amount {
			return c.JSON(http.StatusBadRequest, models.Error{Reason: "insufficient funds"})
		}
		account.Economy.Cash -= amount
	case "vault":
		if account.Economy.Vault < amount {
			return c.JSON(http.StatusBadRequest, models.Error{Reason: "insufficient funds"})
		}
		account.Economy.Vault -= amount
	case "bank":
		if account.Economy.Bank < amount {
			return c.JSON(http.StatusBadRequest, models.Error{Reason: "insufficient funds"})
		}
		account.Economy.Bank -= amount
	case "crypto":
		if account.Economy.Crypto < amount {
			return c.JSON(http.StatusBadRequest, models.Error{Reason: "insufficient funds"})
		}
		account.Economy.Crypto -= amount
	default:
		return c.JSON(http.StatusBadRequest, models.Error{Reason: "invalid type"})
	}

	if err := h.DB.Save(&account.Economy).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{Reason: "unknown error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "withdraw successful"})
}
