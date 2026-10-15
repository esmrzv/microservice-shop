package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/esmrzv/microservice-shop/order-service/internal/middleware"
	"github.com/esmrzv/microservice-shop/order-service/internal/models"
	"github.com/esmrzv/microservice-shop/order-service/internal/service"
)

type fakeOrderService struct {
	createFunc  func(context.Context, uuid.UUID, []*models.OrderItem) (*models.Order, error)
	getByIdFunc func(context.Context, uuid.UUID, uuid.UUID) (*models.Order, error)
	getByUserID func(context.Context, uuid.UUID) ([]*models.Order, error)

	createCalled bool
	gotUserID    uuid.UUID
	gotItems []*models.OrderItem
	gotGetUserID uuid.UUID
	getCalled bool
	gotOrderID uuid.UUID
}

func (s *fakeOrderService) Create(ctx context.Context, userID uuid.UUID, orders []*models.OrderItem) (*models.Order, error) {
	s.createCalled = true
	s.gotUserID = userID
	s.gotItems = orders
	return s.createFunc(ctx, userID, orders)
}

func (s *fakeOrderService) GetByID(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*models.Order, error) {
	s.getCalled = true
	s.gotGetUserID = userID
	s.gotOrderID = orderID
	return s.getByIdFunc(ctx, userID, orderID)
}

func (s *fakeOrderService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	return s.getByUserID(ctx, userID)
}

func TestOrderHandlerCreate(t *testing.T) {
	
	userID := uuid.New()
	productID := uuid.New()

	fakeService := &fakeOrderService{
		createFunc: func(ctx context.Context, gotUserID uuid.UUID, items []*models.OrderItem) (*models.Order, error) {
			return &models.Order{
				ID:         uuid.New(),
				UserID:     gotUserID,
				Status:     "pending",
				TotalPrice: 500,
				Items:      items,
			}, nil
		},
	}
	orderHandler := NewOrderHandler(fakeService)
	body := `{
    	"items": [
        {
            "product_id": "` + productID.String() + `",
            "quantity": 2
        }
    ]
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		bytes.NewBufferString(body),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserIDKey,
		userID,
	)

	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	orderHandler.Create(recorder, req)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status code %d, got %d: %s", http.StatusCreated, recorder.Code, recorder.Body.String())
	}

	if recorder.Header().Get("Content-Type") != "application/json" {
		t.Fatalf(
			"expected Content-Type application/json, got %q",
			recorder.Header().Get("Content-Type"),
		)
	}

	if !fakeService.createCalled {
		t.Fatal("expected service Create to be called")
	}

	if fakeService.gotUserID != userID {
		t.Errorf(
			"expected userID %s, got %s",
			userID,
			fakeService.gotUserID,
		)
	}
	if len(fakeService.gotItems) != 1 {
		t.Fatalf("expected 1 item, got %d", len(fakeService.gotItems))
	}

	gotItem := fakeService.gotItems[0]

	if gotItem.ProductID != productID {
		t.Errorf("expected productID %s, got %s", productID, gotItem.ProductID)
	}
	if gotItem.Quantity != 2 {
    	t.Errorf("expected quantity 2, got %d", gotItem.Quantity)
	}
	var gotOrder models.Order
	if err := json.NewDecoder(recorder.Body).Decode(&gotOrder); err != nil{
		t.Fatalf("failed to decode %v", err)
	}
	if gotOrder.ID == uuid.Nil {
		t.Error("expected non-empty order ID")
	}
	if gotOrder.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, gotOrder.UserID)
	}
	if gotOrder.Status != "pending"{
		t.Errorf("expected status pending, got %s", gotOrder.Status)
	}
	if gotOrder.TotalPrice != 500 {
		t.Errorf("expected total price 500, got %v", gotOrder.TotalPrice)
	}
	if len(gotOrder.Items) != 1{
		t.Fatalf("expected 1 item in response, got %d", len(gotOrder.Items))
	}

}


func TestOrderHandlerInvalidJson(t *testing.T) {
	userID := uuid.New()

	fakeService := &fakeOrderService{
		createFunc: func(ctx context.Context, gotUserID uuid.UUID, items []*models.OrderItem) (*models.Order, error){
			return nil, nil
		},
			
		}
	orderHandler := NewOrderHandler(fakeService)
	body := `{
    	"items": [
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	orderHandler.Create(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if fakeService.createCalled{
		t.Fatal("service Create must not be called for invalid JSON")
	}
}


func TestCreateHandler_UserUnauthorized(t *testing.T){
	fakeService := &fakeOrderService{
		createFunc: func(ctx context.Context,gotUserID uuid.UUID, items[]*models.OrderItem) (*models.Order, error){
			return nil, nil
		},
	}
	orderHandler := NewOrderHandler(fakeService)

	body := `{"items": []}`
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
	recorder := httptest.NewRecorder()
	orderHandler.Create(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Error("expected unuthorized status code")
	}

	if fakeService.createCalled{
		t.Fatal("service Create must not be called for unuthorized")
	}
}

func TestOrderHandler_ProductNotFound(t *testing.T){
	userID := uuid.New()
	productID := uuid.New()

	fakeService := &fakeOrderService{
		createFunc: func(ctx context.Context, userID uuid.UUID, items []*models.OrderItem) (*models.Order, error){
			return nil, service.ErrProductNotFound
		},
	}
	orderHandler := NewOrderHandler(fakeService)

	body := `{
		"items": [
			{
				"product_id": "` + productID.String() + `",
				"quantity": 2
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	orderHandler.Create(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusNotFound,
			recorder.Code,
			recorder.Body.String(),
		)
	}
	if !fakeService.createCalled{
		t.Fatal("expected service Create to be called")
	}
}


func TestOrderHandler_Create_InvalidQuantity(t *testing.T) {
	userID := uuid.New()
	productID := uuid.New()

	fakeService := &fakeOrderService{
		createFunc: func(
			ctx context.Context,
			gotUserID uuid.UUID,
			items []*models.OrderItem,
		) (*models.Order, error) {
			return nil, service.ErrInvalidQuantity
		},
	}

	handler := NewOrderHandler(fakeService)

	body := `{
		"items": [
			{
				"product_id": "` + productID.String() + `",
				"quantity": 0
			}
		]
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		bytes.NewBufferString(body),
	)

	ctx := context.WithValue(
		req.Context(),
		middleware.UserIDKey,
		userID,
	)
	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()
	handler.Create(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			recorder.Code,
			recorder.Body.String(),
		)
	}

	if !fakeService.createCalled {
		t.Fatal("expected service Create to be called")
	}
}


func TestOrderHandler_GetByID(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()
	expectedOrder := &models.Order{
		ID: orderID,
		UserID: userID,
		Status: "pending",
		TotalPrice: 500,
		Items: []*models.OrderItem{},
	}

	fakeService := &fakeOrderService{
		getByIdFunc: func(ctx context.Context, gotUserID uuid.UUID, gotOrderID uuid.UUID) (*models.Order, error){
			return expectedOrder, nil
		},
	}

	orderHandler := NewOrderHandler(fakeService)

	req := httptest.NewRequest(http.MethodGet, "/orders/" + orderID.String(), nil)

	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)

	req.SetPathValue("id", orderID.String())

	req = req.WithContext(ctx)

	recorder := httptest.NewRecorder()

	orderHandler.GetByID(recorder, req)

	if recorder.Code != http.StatusOK {
    	t.Fatalf("expected status %d, got %d",
        http.StatusOK, recorder.Code)
	}
	if !fakeService.getCalled{
		t.Fatal("expected GetByID service method to be called")
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
    	t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	if fakeService.gotGetUserID != userID {

		t.Errorf("expected userID %s, got %s",
        	userID, fakeService.gotGetUserID)
	}

	if fakeService.gotOrderID != orderID{
		t.Errorf("expected orderID %s, got %s",
        orderID, fakeService.gotOrderID)
	}
	var gotOrder models.Order

	if err:= json.NewDecoder(recorder.Body).Decode(&gotOrder);err != nil{
		t.Fatalf("fail decode response %v", err)
	}
	if gotOrder.ID != expectedOrder.ID {
    	t.Errorf("expected order ID %s, got %s",
        expectedOrder.ID, gotOrder.ID)
	}
	if gotOrder.UserID != expectedOrder.UserID {
    	t.Errorf("expected user ID %s, got %s",
        expectedOrder.UserID, gotOrder.UserID)
	}
	if gotOrder.Status != expectedOrder.Status {
    	t.Errorf("expected status %q, got %q",
        expectedOrder.Status, gotOrder.Status)
	}

	if gotOrder.TotalPrice != expectedOrder.TotalPrice {
    	t.Errorf("expected total price %v, got %v",
        expectedOrder.TotalPrice, gotOrder.TotalPrice)
	}
}


func TestOrderHandler_Unauthorized(t *testing.T){
	orderID := uuid.New()

	fakeService := &fakeOrderService{}
	orderHandler := NewOrderHandler(fakeService)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	req.SetPathValue("id", orderID.String())
	recorder := httptest.NewRecorder()
	orderHandler.GetByID(recorder, req)

	if recorder.Code != http.StatusUnauthorized{
		t.Fatalf("expected status %d, got %d",
            http.StatusUnauthorized, recorder.Code)
	}

	if fakeService.getCalled {
		t.Fatal("service must not be called without authentication")
	}
}

func TestOrderHandler_InvalidUUID(t *testing.T){
	fakeService := &fakeOrderService{}
	orderHandler := NewOrderHandler(fakeService)

	req := httptest.NewRequest(http.MethodGet, "/orders/abc", nil)
	userID := uuid.New()
	req.SetPathValue("id", "abc")
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	orderHandler.GetByID(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	if fakeService.getCalled{
		t.Fatal("service must not be called without incorrect orderID")
	}
}



func TestOrderHandler_GetByID_ProductNotFound(t *testing.T) {
	userID := uuid.New()
	orderID := uuid.New()

	fakeService := &fakeOrderService{
		getByIdFunc: func(ctx context.Context, gotUserID uuid.UUID, gotOrderID uuid.UUID) (*models.Order, error){
			return nil, service.ErrProductNotFound
		},
	}
	orderHandler := NewOrderHandler(fakeService)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	req.SetPathValue("id", orderID.String())
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
	req = req.WithContext(ctx)
	recorder := httptest.NewRecorder()
	orderHandler.GetByID(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	if !fakeService.getCalled{
		t.Fatal("service must not be called without incorrect orderID")
	}
}


func TestOrderHandler_GetByID_InternalError(t *testing.T) {
    userID := uuid.New()
    orderID := uuid.New()

    fakeService := &fakeOrderService{
        getByIdFunc: func(
            ctx context.Context,
            gotUserID uuid.UUID,
            gotOrderID uuid.UUID,
        ) (*models.Order, error) {
            return nil, errors.New("database unavailable")
        },
    }

    handler := NewOrderHandler(fakeService)

    req := httptest.NewRequest(
        http.MethodGet,
        "/orders/"+orderID.String(),
        nil,
    )
    req.SetPathValue("id", orderID.String())

    ctx := context.WithValue(
        req.Context(),
        middleware.UserIDKey,
        userID,
    )
    req = req.WithContext(ctx)

    recorder := httptest.NewRecorder()
    handler.GetByID(recorder, req)

    if recorder.Code != http.StatusInternalServerError {
        t.Fatalf("expected status %d, got %d",
            http.StatusInternalServerError, recorder.Code)
    }

    if !fakeService.getCalled {
        t.Fatal("expected GetByID service method to be called")
    }
}

func TestOrderHandler_GetByUserID(t *testing.T) {
    userID := uuid.New()

    expectedOrders := []*models.Order{
        {
            ID:         uuid.New(),
            UserID:     userID,
            Status:     "pending",
            TotalPrice: 500,
            Items:      []*models.OrderItem{},
        },
        {
            ID:         uuid.New(),
            UserID:     userID,
            Status:     "completed",
            TotalPrice: 750,
            Items:      []*models.OrderItem{},
        },
    }

    serviceCalled := false
    var gotUserID uuid.UUID

    fakeService := &fakeOrderService{
        getByUserID: func(
            ctx context.Context,
            id uuid.UUID,
        ) ([]*models.Order, error) {
            serviceCalled = true
            gotUserID = id
            return expectedOrders, nil
        },
    }

    handler := NewOrderHandler(fakeService)

    req := httptest.NewRequest(
        http.MethodGet,
        "/orders",
        nil,
    )

    ctx := context.WithValue(
        req.Context(),
        middleware.UserIDKey,
        userID,
    )
    req = req.WithContext(ctx)

    recorder := httptest.NewRecorder()
    handler.GetByUserID(recorder, req)

	if recorder.Code != http.StatusOK {
   		 t.Fatalf("expected status %d, got %d",
        http.StatusOK, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
    	t.Fatalf("expected Content-Type application/json, got %q", contentType)
	}

	if !serviceCalled {
    	t.Fatal("expected GetByUserID service method to be called")
	}

	if gotUserID != userID {
    	t.Errorf("expected userID %s, got %s", userID, gotUserID)
	}

	var gotOrders []*models.Order
	if err := json.NewDecoder(recorder.Body).Decode(&gotOrders); err != nil {
    	t.Fatalf("failed to decode response: %v", err)
	}

	if len(gotOrders) != len(expectedOrders) {
    	t.Fatalf("expected %d orders, got %d",
        len(expectedOrders), len(gotOrders))
	}

	for i, expected := range expectedOrders {
    	if gotOrders[i].ID != expected.ID {
        	t.Errorf("order %d: expected ID %s, got %s",
            	i, expected.ID, gotOrders[i].ID)
    	}

    	if gotOrders[i].UserID != expected.UserID {
        	t.Errorf("order %d: expected userID %s, got %s",
            i, expected.UserID, gotOrders[i].UserID)
    	}

    	if gotOrders[i].Status != expected.Status {
        	t.Errorf("order %d: expected status %q, got %q",
            i, expected.Status, gotOrders[i].Status)
    	}

    	if gotOrders[i].TotalPrice != expected.TotalPrice {
        	t.Errorf("order %d: expected total price %v, got %v",
            i, expected.TotalPrice, gotOrders[i].TotalPrice)
    	}
	}

    
}

func TestOrderHandler_GetByUserID_Unauthorized(t *testing.T) {
    fakeService := &fakeOrderService{}
    handler := NewOrderHandler(fakeService)

    req := httptest.NewRequest(
        http.MethodGet,
        "/orders",
        nil,
    )

    recorder := httptest.NewRecorder()
    handler.GetByUserID(recorder, req)

    if recorder.Code != http.StatusUnauthorized {
        t.Fatalf("expected status %d, got %d",
            http.StatusUnauthorized, recorder.Code)
    }
}

func TestOrderHandler_GetByUserID_InternalError(t *testing.T) {
    userID := uuid.New()

    fakeService := &fakeOrderService{
        getByUserID: func(
            ctx context.Context,
            id uuid.UUID,
        ) ([]*models.Order, error) {
            return nil, errors.New("database unavailable")
        },
    }

    handler := NewOrderHandler(fakeService)

    req := httptest.NewRequest(
        http.MethodGet,
        "/orders",
        nil,
    )

    ctx := context.WithValue(
        req.Context(),
        middleware.UserIDKey,
        userID,
    )
    req = req.WithContext(ctx)

    recorder := httptest.NewRecorder()
    handler.GetByUserID(recorder, req)

    if recorder.Code != http.StatusInternalServerError {
        t.Fatalf("expected status %d, got %d",
            http.StatusInternalServerError, recorder.Code)
    }
}

func TestOrderHandler_GetByUserID_Empty(t *testing.T) {
    userID := uuid.New()

    fakeService := &fakeOrderService{
        getByUserID: func(
            ctx context.Context,
            id uuid.UUID,
        ) ([]*models.Order, error) {
            return []*models.Order{}, nil
        },
    }

    handler := NewOrderHandler(fakeService)

    req := httptest.NewRequest(
        http.MethodGet,
        "/orders",
        nil,
    )

    ctx := context.WithValue(
        req.Context(),
        middleware.UserIDKey,
        userID,
    )
    req = req.WithContext(ctx)

    recorder := httptest.NewRecorder()
    handler.GetByUserID(recorder, req)

    if recorder.Code != http.StatusOK {
        t.Fatalf("expected status %d, got %d",
            http.StatusOK, recorder.Code)
    }

    var gotOrders []*models.Order
    if err := json.NewDecoder(recorder.Body).Decode(&gotOrders); err != nil {
        t.Fatalf("failed to decode response: %v", err)
    }

    if len(gotOrders) != 0 {
        t.Fatalf("expected empty orders, got %d", len(gotOrders))
    }
}