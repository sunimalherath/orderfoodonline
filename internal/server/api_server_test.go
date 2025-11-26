package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/sunimalherath/orderfoodonline/internal/core/entities"
)

type mockProductService struct {
	listProductsFunc    func(ctx context.Context) ([]entities.Product, error)
	findProductByIDFunc func(ctx context.Context, productID int64) (*entities.Product, error)
}

func (m *mockProductService) ListProducts(ctx context.Context) ([]entities.Product, error) {
	if m.listProductsFunc != nil {
		return m.listProductsFunc(ctx)
	}

	return nil, nil
}

func (m *mockProductService) FindProductByID(ctx context.Context, productID int64) (*entities.Product, error) {
	if m.findProductByIDFunc != nil {
		return m.findProductByIDFunc(ctx, productID)
	}

	return nil, nil
}

type mockOrderService struct {
	placeAnOrderFunc func(ctx context.Context, orderReq entities.OrderReq) (*entities.Order, error)
}

func (m *mockOrderService) PlaceAnOrder(ctx context.Context, orderReq entities.OrderReq) (*entities.Order, error) {
	if m.placeAnOrderFunc != nil {
		return m.placeAnOrderFunc(ctx, orderReq)
	}

	return nil, nil
}

func newTestServer(prodSvc *mockProductService, orderSvc *mockOrderService) *apiServer {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	return &apiServer{
		prodSvc:  prodSvc,
		orderSvc: orderSvc,
		apiKey:   "test-api-key",
		logger:   logger,
	}
}

func TestListProducts_Success(t *testing.T) {
	expectedStatus := http.StatusOK

	mockFunc := func(ctx context.Context) ([]entities.Product, error) {
		return []entities.Product{
			{ID: "1", Name: "Product 1", Category: "Category A", Price: 10.99},
			{ID: "2", Name: "Product 2", Category: "Category B", Price: 20.99},
		}, nil
	}

	mockProdSvc := &mockProductService{
		listProductsFunc: mockFunc,
	}

	server := newTestServer(mockProdSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/product", nil)
	w := httptest.NewRecorder()

	server.ListProducts(w, req)

	if w.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
	}
}

func TestListProducts_Failure(t *testing.T) {
	expectedStatus := http.StatusInternalServerError

	mockFunc := func(ctx context.Context) ([]entities.Product, error) {
		return nil, errors.New("random test error")
	}

	mockProdSvc := &mockProductService{
		listProductsFunc: mockFunc,
	}

	server := newTestServer(mockProdSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/product", nil)
	w := httptest.NewRecorder()

	server.ListProducts(w, req)

	if w.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
	}
}

func TestFindProductByID_Success(t *testing.T) {
	expectedStatus := http.StatusOK
	productID := "1"

	mockFunc := func(ctx context.Context, productID int64) (*entities.Product, error) {
		return &entities.Product{
			ID:       "1",
			Name:     "Product 1",
			Category: "Category A",
			Price:    10.99,
		}, nil
	}

	mockProdSvc := &mockProductService{
		findProductByIDFunc: mockFunc,
	}
	server := newTestServer(mockProdSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/product/"+productID, nil)
	req.SetPathValue("productId", productID)
	w := httptest.NewRecorder()

	server.FindProductByID(w, req)

	if w.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
	}
}

func TestFindProductByID_Failure(t *testing.T) {
	expectedStatus := http.StatusInternalServerError
	productID := "1"

	mockFunc := func(ctx context.Context, productID int64) (*entities.Product, error) {
		return nil, nil
	}

	mockProdSvc := &mockProductService{
		findProductByIDFunc: mockFunc,
	}

	server := newTestServer(mockProdSvc, nil)

	req := httptest.NewRequest(http.MethodGet, "/product/"+productID, nil)
	req.SetPathValue("productId", productID)
	w := httptest.NewRecorder()

	server.FindProductByID(w, req)

	if w.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
	}
}

func TestPlaceAnOrder_Success(t *testing.T) {
	expectedStatus := http.StatusOK
	mockOrderID := uuid.New().String()

	var orderReq any = entities.OrderReq{
		Items: []entities.OrderItem{
			{ProductID: "1", Quantity: 2},
		},
		CouponCode: "",
	}

	mockFunc := func(ctx context.Context, orderReq entities.OrderReq) (*entities.Order, error) {
		return &entities.Order{
			ID: mockOrderID,
			Items: []entities.OrderItem{
				{ProductID: "1", Quantity: 2},
			},
		}, nil
	}

	mockOrderSvc := &mockOrderService{
		placeAnOrderFunc: mockFunc,
	}
	server := newTestServer(nil, mockOrderSvc)

	var body []byte
	var err error

	if str, ok := orderReq.(string); ok {
		body = []byte(str)
	} else {
		body, err = json.Marshal(orderReq)
		if err != nil {
			t.Fatalf("failed to marshal order request: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.PlaceAnOrder(w, req)

	if w.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
	}
}

func TestPlaceAnOrder_Failure(t *testing.T) {
	expectedStatus := http.StatusUnprocessableEntity

	var orderReq any = entities.OrderReq{
		Items: []entities.OrderItem{
			{ProductID: "", Quantity: 2},
		},
		CouponCode: "",
	}

	mockOrderSvc := &mockOrderService{}

	server := newTestServer(nil, mockOrderSvc)

	var body []byte
	var err error

	if str, ok := orderReq.(string); ok {
		body = []byte(str)
	} else {
		body, err = json.Marshal(orderReq)
		if err != nil {
			t.Fatalf("failed to marshal order request: %v", err)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.PlaceAnOrder(w, req)

	if w.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, w.Code)
	}
}
