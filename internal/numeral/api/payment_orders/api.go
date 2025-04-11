package payment_orders

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/erecarte/showcase/pkg/api/models"
	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"net/http"
	"strconv"
	"time"
)

//go:embed request_schema.json
var requestSchema embed.FS

type Api struct {
	router     gin.IRoutes
	service    *Service
	jsonSchema gojsonschema.JSONLoader
}

type PaymentOrderStorage interface {
	Store(ctx context.Context, po *PaymentOrder) error
	Get(ctx context.Context, id string) (*PaymentOrder, error)
	UpdateStatus(ctx context.Context, id string, status string) error
}

func NewApi(router gin.IRoutes, service *Service) (*Api, error) {
	schemaFile, err := requestSchema.ReadFile("request_schema.json")
	if err != nil {

	}
	schemaLoader := gojsonschema.NewBytesLoader(schemaFile)
	api := &Api{
		router:     router,
		service:    service,
		jsonSchema: schemaLoader,
	}
	router.POST("/payment_orders", errorHandler(api.createPaymentOrder))
	router.GET("/payment_orders/:id", errorHandler(api.retrievePaymentOrder))
	//router.GET("/payment_orders", errorHandler(api.listPaymentOrders))
	//router.POST("/payment_orders/:id", errorHandler(api.updatePaymentOrder))
	return api, nil
}

func (a *Api) retrievePaymentOrder(c *gin.Context) error {
	id := c.Param("id")
	if id == "" {
		return ErrInvalidRequest
	}
	po, err := a.service.GetPaymentOrder(c.Request.Context(), id)
	if err != nil {
		return ErrRecordNotFound
	}
	c.JSON(http.StatusOK, po)
	return nil
}

func (a *Api) createPaymentOrder(c *gin.Context) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ErrInvalidRequest
	}
	result, err := gojsonschema.Validate(a.jsonSchema, gojsonschema.NewBytesLoader(body))
	if err != nil {
		return ErrInvalidRequest
	}
	if !result.Valid() {
		return ErrInvalidRequest
	}
	apiModel := &models.PaymentOrderApiModel{}
	err = json.Unmarshal(body, apiModel)
	if err != nil {
		return ErrInvalidRequest
	}
	amount := fmt.Sprintf("%f", apiModel.Ammount)
	now := time.Now().UTC()
	po, err := a.service.CreatePaymentOrder(c.Request.Context(), &PaymentOrder{
		DebtorIban:           apiModel.DebtorIban,
		DebtorName:           apiModel.DebtorName,
		CreditorIban:         apiModel.CreditorIban,
		CreditorName:         apiModel.CreditorName,
		Amount:               amount,
		IdempotencyUniqueKey: apiModel.IdempotencyUniqueKey,
		Status:               "PENDING",
		CreatedAt:            &now,
	})
	if err != nil {
		return ErrRecordAlreadyExists
	}
	ammount, err := strconv.ParseFloat(po.Amount, 64)
	if err != nil {
		return errors.New("invalid amount")
	}
	c.JSON(http.StatusOK, &models.PaymentOrderApiModel{
		DebtorIban:           po.DebtorIban,
		DebtorName:           po.DebtorName,
		CreditorIban:         po.CreditorIban,
		CreditorName:         po.CreditorName,
		Ammount:              ammount,
		IdempotencyUniqueKey: po.IdempotencyUniqueKey,
	})
	return nil
}

func errorHandler(handler func(c *gin.Context) error) func(*gin.Context) {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			switch {
			case errors.Is(err, ErrRecordNotFound):
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.Is(err, ErrRecordAlreadyExists):
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			case errors.Is(err, ErrInvalidRequest):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
		}
	}
}

func (a *Api) Start() {

}
