package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/middleware"
	"github.com/esmrzv/product-service/internal/service"
	"github.com/google/uuid"
)

func TestProductHandler_CreateProduct(t *testing.T) {
	categoryID := uuid.New()

	reqBody := dto.CreateProductRequest{
		CategoryID:  categoryID,
		Name:        "test product",
		Description: "test description",
		Price:       1000,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	request = request.WithContext(
		context.WithValue(
			request.Context(),
			middleware.UserIDKey,
			uuid.New(),
		),
	)
	recorder := httptest.NewRecorder()
	product := &dto.ProductResponse{
		ID:          uuid.New(),
		CategoryID:  reqBody.CategoryID,
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Price:       reqBody.Price,
	}
	productService := &mockProductService{
		createProductFunc: func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
			return product, nil
		},
	}
	handler := NewProductHandler(productService)
	handler.CreateProduct(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Errorf("got %d, want %d", recorder.Code, http.StatusCreated)
	}
}

func TestProductHandler_CreateProduct_InvalidJSON(t *testing.T) {
	body := []byte("")
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	productService := &mockProductService{
		createProductFunc: func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
			t.Fatal("CreateProduct should not be called")
			return nil, nil
		},
	}
	handler := NewProductHandler(productService)
	handler.CreateProduct(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("got %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestProductHandler_CreateProduct_CategoryNotFound(t *testing.T) {
	categoryID := uuid.New()
	reqBody := &dto.CreateProductRequest{
		CategoryID:  categoryID,
		Name:        "test product",
		Description: "test description",
		Price:       1000,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req = req.WithContext(
		context.WithValue(
			req.Context(),
			middleware.UserIDKey,
			uuid.New(),
		),
	)
	recorder := httptest.NewRecorder()
	productService := &mockProductService{
		createProductFunc: func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {

			return nil, service.ErrCategoryNotFound
		},
	}
	handler := NewProductHandler(productService)
	handler.CreateProduct(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Errorf("got %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestProductHandler_CreateProduct_ServiceError(t *testing.T) {
	repoErr := errors.New("database error")
	categoryID := uuid.New()
	reqBody := dto.CreateProductRequest{
		CategoryID:  categoryID,
		Name:        "test product",
		Description: "test description",
		Price:       1000,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(body))
	req = req.WithContext(
		context.WithValue(
			req.Context(),
			middleware.UserIDKey,
			uuid.New(),
		),
	)
	recorder := httptest.NewRecorder()
	productService := &mockProductService{
		createProductFunc: func(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {

			return nil, repoErr
		},
	}
	handler := NewProductHandler(productService)
	handler.CreateProduct(recorder, req)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("got %d, want %d", recorder.Code, http.StatusInternalServerError)
	}
}
