package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/esmrzv/product-service/internal/dto"
	"github.com/esmrzv/product-service/internal/middleware"
	"github.com/esmrzv/product-service/internal/models"
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

func TestProductHandler_GetProductByID(t *testing.T) {
	productID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)
	req = req.WithContext(
		context.WithValue(
			req.Context(),
			middleware.UserIDKey,
			uuid.New(),
		),
	)

	recorder := httptest.NewRecorder()
	productResponce := dto.ProductResponse{
		ID:          productID,
		CategoryID:  uuid.New(),
		Name:        "geforce",
		Description: "rtx 2030",
		Price:       23232,
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
	}
	productService := &mockProductService{
		getProductByIDFunc: func(ctx context.Context, productID uuid.UUID) (*dto.ProductResponse, error) {
			return &productResponce, nil
		},
	}

	router := setupHandlerRoute(productService)

	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", recorder.Code, http.StatusOK)
	}
	var response dto.ProductResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatal("failed to decode body")
	}
	if response.ID != productResponce.ID {
		t.Fatalf("got ID %v, want %v", response.ID, productResponce.ID)
	}
	if response.Name != productResponce.Name {
		t.Fatalf("got name %q, want %q", response.Name, productResponce.Name)
	}

}

func TestProductHandler_GetProductByID_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/not-uuid", nil)
	recorder := httptest.NewRecorder()

	productService := &mockProductService{
		getProductByIDFunc: func(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
			t.Fatal("GetProductByID should not be called")
			return nil, nil
		},
	}
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestProductHandler_GetProductByID_NotFound(t *testing.T) {
	productID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)
	recorder := httptest.NewRecorder()

	productService := &mockProductService{
		getProductByIDFunc: func(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
			return nil, service.ErrProductNotFound
		},
	}

	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusNotFound)
	}
}

func TestProductHandler_GetProductByID_RepoError(t *testing.T) {
	repoErr := errors.New("database error")
	productID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)
	recorder := httptest.NewRecorder()

	productService := &mockProductService{
		getProductByIDFunc: func(ctx context.Context, id uuid.UUID) (*dto.ProductResponse, error) {
			return nil, repoErr
		},
	}

	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusInternalServerError)
	}
}

func TestProductHandler_GetAll(t *testing.T) {
	product1ID := uuid.New()
	product2ID := uuid.New()

	categoryID := uuid.New()
	product1 := models.Product{

		ID:          product1ID,
		CategoryID:  categoryID,
		Name:        "iphone15",
		Description: "pro max",
		Price:       202000,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	product2 := models.Product{
		ID:          product2ID,
		CategoryID:  categoryID,
		Name:        "iphone17",
		Description: "pro",
		Price:       202002,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	productService := &mockProductService{
		getAllProductsFunc: func(ctx context.Context) ([]dto.ProductResponse, error) {
			return []dto.ProductResponse{
				{
					ID:          product1.ID,
					CategoryID:  product1.CategoryID,
					Name:        product1.Name,
					Description: product1.Description,
					Price:       product1.Price,
					CreatedAt:   product1.CreatedAt,
					UpdatedAt:   product1.UpdatedAt,
				},
				{
					ID:          product2.ID,
					CategoryID:  product2.CategoryID,
					Name:        product2.Name,
					Description: product2.Description,
					Price:       product2.Price,
					CreatedAt:   product2.CreatedAt,
					UpdatedAt:   product2.UpdatedAt,
				},
			}, nil
		}}

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("got %d, want %d", recorder.Code, http.StatusOK)
	}
	var response []dto.ProductResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatal("failed to decode recorder body")
	}
	if len(response) != 2 {
		t.Fatalf("expected 2 products")
	}
	if response[0].ID != product1ID {
		t.Fatalf("got ID %v, want %v", response[0].ID, product1.ID)
	}
	if response[0].Name != product1.Name {
		t.Fatalf("got name %s ,want %s", response[0].Name, product1.Name)
	}
	if response[0].Price != product1.Price {
		t.Fatalf("got price %f, want %f", response[0].Price, product1.Price)
	}

}

func TestProductHandler_GetAll_ServiceError(t *testing.T) {
	repoErr := errors.New("database error")
	productService := &mockProductService{
		getAllProductsFunc: func(ctx context.Context) ([]dto.ProductResponse, error) {
			return nil, repoErr
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusInternalServerError)
	}
}

func TestProductHandler_GetAll_Empty(t *testing.T) {
	productService := &mockProductService{
		getAllProductsFunc: func(ctx context.Context) ([]dto.ProductResponse, error) {
			return []dto.ProductResponse{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusOK)
	}
	var response []dto.ProductResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	if err != nil {
		t.Fatal("failed to decode recorder body")
	}
	if len(response) != 0 {
		t.Fatalf("got %d products, want 0", len(response))
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	productID := uuid.New()
	updateReq := dto.UpdateProductRequest{
		Name:        "iphone15",
		Description: "pro max",
		Price:       202000,
	}
	body, err := json.Marshal(updateReq)
	if err != nil {
		t.Fatal("failed to marshal request body")
	}

	req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewReader(body))
	req = req.WithContext(
		context.WithValue(
			req.Context(),
			middleware.UserIDKey,
			uuid.New(),
		),
	)
	recorder := httptest.NewRecorder()
	productService := &mockProductService{
		updateProductFunc: func(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error {
			return nil
		},
	}
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusNoContent)
	}

}
func TestProductHandler_UpdateProduct_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/products/uuid", nil)
	recorder := httptest.NewRecorder()
	productService := &mockProductService{
		updateProductFunc: func(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error {
			t.Fatal("service do not call")
			return nil
		},
	}
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusBadRequest)
	}
}

func TestProductHandler_UpdateProduct_InvalidJson(t *testing.T) {
	productID := uuid.New()
	request := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), nil)
	recorder := httptest.NewRecorder()
	productService := &mockProductService{
		updateProductFunc: func(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error {
			t.Fatal("service do not call")
			return nil
		},
	}
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusBadRequest)
	}

}

func TestProductHandler_UpdateProduct_ServiceError(t *testing.T) {
	repoErr := errors.New("database error")
	productID := uuid.New()
	updateReq := dto.UpdateProductRequest{
		Name:        "iphone15",
		Description: "pro max",
		Price:       202000,
	}
	body, err := json.Marshal(updateReq)
	if err != nil {
		t.Fatal("failed to marshal request body")
	}

	req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewReader(body))
	productService := &mockProductService{
		updateProductFunc: func(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error {
			return repoErr
		},
	}
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusInternalServerError)
	}

}

func TestProductHandler_UpdateProduct_ProductNotFound(t *testing.T) {
	productID := uuid.New()
	updateReq := dto.UpdateProductRequest{
		Name:        "iphone15",
		Description: "pro max",
		Price:       202000,
	}
	body, err := json.Marshal(updateReq)
	if err != nil {
		t.Fatal("failed to marshal request body")
	}
	req := httptest.NewRequest(http.MethodPut, "/products/"+productID.String(), bytes.NewReader(body))
	productService := &mockProductService{
		updateProductFunc: func(ctx context.Context, productID uuid.UUID, product *dto.UpdateProductRequest) error {
			return service.ErrProductNotFound
		},
	}
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusNotFound)
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	productID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/products/"+productID.String(), nil)
	productService := &mockProductService{
		deleteProductFunc: func(ctx context.Context, productID uuid.UUID) error {
			return nil
		},
	}
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusNoContent)
	}
}

func TestProductHandler_DeleteProduct_InvalidUUID(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/products/invalid-uuid", nil)
	productService := &mockProductService{
		deleteProductFunc: func(ctx context.Context, productID uuid.UUID) error {
			t.Fatal("service do not call")
			return nil
		},
	}
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusBadRequest)
	}
}

func TestProductHandler_DeleteProduct_ProductNotFound(t *testing.T) {
	productID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/products/"+productID.String(), nil)
	productService := &mockProductService{
		deleteProductFunc: func(ctx context.Context, productID uuid.UUID) error {
			return service.ErrProductNotFound
		},
	}
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusNotFound)
	}
}

func TestNewProductHandler_DeleteProduct_ServerError(t *testing.T) {
	repoErr := errors.New("database error")
	req := httptest.NewRequest(http.MethodDelete, "/products/"+uuid.New().String(), nil)
	productService := &mockProductService{
		deleteProductFunc: func(ctx context.Context, productID uuid.UUID) error {
			return repoErr
		},
	}
	recorder := httptest.NewRecorder()
	router := setupHandlerRoute(productService)
	router.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("got %d, want %d ", recorder.Code, http.StatusInternalServerError)
	}
}
