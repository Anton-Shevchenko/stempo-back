package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stempo/backend/internal/domain/entity"
	"github.com/stempo/backend/internal/infrastructure/database"
	"github.com/stempo/backend/internal/infrastructure/email"
	"github.com/stempo/backend/internal/infrastructure/oauth"
	postgresRepo "github.com/stempo/backend/internal/repository/postgres"
	"github.com/stempo/backend/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const testJWTSecret = "test-secret-key-for-testing-only"

type APITestSuite struct {
	suite.Suite
	db         *gorm.DB
	router     *gin.Engine
	userID     uint
	businessID uint
}

func (suite *APITestSuite) SetupSuite() {
	// Set test environment variables
	os.Setenv("JWT_SECRET", testJWTSecret)
	gin.SetMode(gin.TestMode)

	// Setup test database
	var err error
	suite.db, err = setupTestDB()
	if err != nil {
		suite.T().Fatal("Failed to setup test database:", err)
	}

	// Setup test router
	suite.router, err = setupTestRouter(suite.db)
	if err != nil {
		suite.T().Fatal("Failed to setup test router:", err)
	}
}

func (suite *APITestSuite) TearDownSuite() {
	if suite.db != nil {
		clearTestDB(suite.db)
	}
}

func setupTestDB() (*gorm.DB, error) {
	dsn := "host=localhost user=stempo_test password=stempo_test_password dbname=stempo_test_db port=5433 sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgresDriver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		return nil, err
	}

	// Clear all tables
	if err := clearTestDB(db); err != nil {
		return nil, err
	}

	// Seed test data
	if err := seedTestData(db); err != nil {
		return nil, err
	}

	return db, nil
}

func setupTestRouter(db *gorm.DB) (*gin.Engine, error) {
	userRepo := postgresRepo.NewUserRepository(db)
	businessRepo := postgresRepo.NewBusinessRepository(db)
	programRepo := postgresRepo.NewBonusProgramRepository(db)
	cardRepo := postgresRepo.NewLoyaltyCardRepository(db)
	categoryRepo := postgresRepo.NewCategoryRepository(db)
	cityRepo := postgresRepo.NewCityRepository(db)
	employeeRepo := postgresRepo.NewEmployeeRepository(db)
	redemptionRepo := postgresRepo.NewBonusRedemptionRepository(db)
	qrCodeRepo := postgresRepo.NewQRCodeRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, oauth.NewGoogleTokenVerifier())
	businessUsecase := usecase.NewBusinessUsecase(businessRepo)
	programUsecase := usecase.NewBonusProgramUsecase(programRepo, businessRepo)
	employeeUsecase := usecase.NewEmployeeUsecase(employeeRepo, businessRepo, userRepo, email.NewMailjetService())
	cardUsecase := usecase.NewLoyaltyCardUsecase(cardRepo, businessRepo, employeeUsecase)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	cityUsecase := usecase.NewCityUsecase(cityRepo)
	redemptionUsecase := usecase.NewBonusRedemptionUsecase(redemptionRepo, cardRepo, businessRepo, employeeUsecase)
	qrCodeUsecase := usecase.NewQRCodeUsecase(qrCodeRepo, programRepo, businessRepo)

	router := SetupRouter(authUsecase, businessUsecase, programUsecase, cardUsecase, categoryUsecase, cityUsecase, employeeUsecase, redemptionUsecase, qrCodeUsecase)
	return router, nil
}

func clearTestDB(db *gorm.DB) error {
	tables := []string{
		"loyalty_cards",
		"bonus_programs",
		"businesses",
		"users",
		"cities",
		"categories",
	}

	for _, table := range tables {
		if err := db.Exec("TRUNCATE TABLE " + table + " CASCADE").Error; err != nil {
			return err
		}
	}

	return nil
}

func seedTestData(db *gorm.DB) error {
	// Seed cities
	kyivLat := 50.4501
	kyivLng := 30.5234
	cities := []entity.City{
		{Name: "Kyiv", Slug: "kyiv", Latitude: &kyivLat, Longitude: &kyivLng},
		{Name: "Lviv", Slug: "lviv"},
	}
	for i := range cities {
		cities[i].CreatedAt = time.Now()
		cities[i].UpdatedAt = time.Now()
		if err := db.Create(&cities[i]).Error; err != nil {
			return err
		}
	}

	// Seed categories
	categories := []entity.Category{
		{Name: "Coffee", Slug: "coffee", Icon: "☕", IconColor: "#8B5CF6"},
		{Name: "Sports", Slug: "sports", Icon: "💪", IconColor: "#EF4444"},
		{Name: "Food", Slug: "food", Icon: "🍝", IconColor: "#F59E0B"},
	}
	for i := range categories {
		categories[i].CreatedAt = time.Now()
		categories[i].UpdatedAt = time.Now()
		if err := db.Create(&categories[i]).Error; err != nil {
			return err
		}
	}

	// Seed test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
	name := "Test User"
	phone := "+380501234567"
	testUser := entity.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     &name,
		Phone:    &phone,
		CityID:   &cities[0].ID,
	}
	testUser.CreatedAt = time.Now()
	testUser.UpdatedAt = time.Now()
	if err := db.Create(&testUser).Error; err != nil {
		return err
	}

	// Seed test business
	testBusiness := entity.Business{
		Name:              "Test Business",
		Category:          "coffee",
		Address:           "Test Address 123",
		Rating:            4.5,
		IsOpen:            true,
		HasLoyaltyProgram: true,
		Featured:          false,
		Status:            entity.BusinessStatusApproved,
		OwnerID:           testUser.ID,
		Icon:              "☕",
		IconColor:         "#8B5CF6",
	}
	testBusiness.CreatedAt = time.Now()
	testBusiness.UpdatedAt = time.Now()
	if err := db.Create(&testBusiness).Error; err != nil {
		return err
	}

	// Seed test bonus programs
	programs := []entity.BonusProgram{
		{
			Name:           "Test Program 1",
			Description:    "Test description 1",
			BusinessID:     testBusiness.ID,
			PointsRequired: 10,
			Discount:       20,
			DiscountType:   "percentage",
			Status:         entity.BonusProgramStatusPending,
		},
		{
			Name:           "Test Program 2",
			Description:    "Test description 2",
			BusinessID:     testBusiness.ID,
			PointsRequired: 15,
			Discount:       30,
			DiscountType:   "percentage",
			Status:         entity.BonusProgramStatusApproved,
		},
	}
	for i := range programs {
		programs[i].CreatedAt = time.Now()
		programs[i].UpdatedAt = time.Now()
		if err := db.Create(&programs[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

func createAuthRequest(method, url string, body interface{}, userID uint) *http.Request {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = testJWTSecret
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	req.Header.Set("Authorization", "Bearer "+tokenString)

	return req
}

func createRequest(method, url string, body interface{}) *http.Request {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req = httptest.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	return req
}

func executeRequest(router *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

func (suite *APITestSuite) getTestUserID() uint {
	var user entity.User
	suite.db.Where("email = ?", "test@example.com").First(&user)
	return user.ID
}

func (suite *APITestSuite) getTestBusinessID() uint {
	var business entity.Business
	suite.db.Where("name = ?", "Test Business").First(&business)
	return business.ID
}

// ========== Auth Routes Tests ==========

func (suite *APITestSuite) TestAuthRegister() {
	// Get first city ID from database
	var city entity.City
	suite.db.First(&city)

	body := map[string]interface{}{
		"email":    "newuser@example.com",
		"password": "password123",
		"name":     "New User",
		"phone":    "+380501111111",
		"cityId":   city.ID,
	}
	req := createRequest("POST", "/api/auth/register", body)
	recorder := executeRequest(suite.router, req)

	if recorder.Code != 201 {
		var errorResp map[string]interface{}
		json.Unmarshal(recorder.Body.Bytes(), &errorResp)
		suite.T().Logf("Registration failed: %v", errorResp)
	}
	assert.Equal(suite.T(), 201, recorder.Code)
}

func (suite *APITestSuite) TestAuthLogin() {
	body := map[string]string{
		"email":    "test@example.com",
		"password": "test123",
	}
	req := createRequest("POST", "/api/auth/login", body)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
	var response map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.Contains(suite.T(), response, "token")
}

func (suite *APITestSuite) TestAuthGetCurrentUser() {
	userID := suite.getTestUserID()
	req := createAuthRequest("GET", "/api/auth/me", nil, userID)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

// ========== Business Routes Tests ==========

func (suite *APITestSuite) TestBusinessGetAll() {
	req := createRequest("GET", "/api/businesses", nil)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

func (suite *APITestSuite) TestBusinessGetByID() {
	businessID := suite.getTestBusinessID()
	req := createRequest("GET", fmt.Sprintf("/api/businesses/%d", businessID), nil)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

func (suite *APITestSuite) TestBusinessGetMyBusiness() {
	userID := suite.getTestUserID()
	req := createAuthRequest("GET", "/api/businesses/my-business", nil, userID)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

func (suite *APITestSuite) TestBusinessCreate() {
	userID := suite.getTestUserID()
	body := map[string]interface{}{
		"name":              "New Test Business",
		"category":          "coffee",
		"address":           "New Address 456",
		"hasLoyaltyProgram": true,
	}
	req := createAuthRequest("POST", "/api/businesses", body, userID)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 201, recorder.Code)
}

// ========== Bonus Program Routes Tests ==========

func (suite *APITestSuite) TestBonusProgramGetAll() {
	req := createRequest("GET", "/api/bonus-programs", nil)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

func (suite *APITestSuite) TestBonusProgramGetByID() {
	// Get first program ID from database
	var program entity.BonusProgram
	suite.db.First(&program)

	req := createRequest("GET", fmt.Sprintf("/api/bonus-programs/%d", program.ID), nil)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

func (suite *APITestSuite) TestBonusProgramGetByBusinessID() {
	userID := suite.getTestUserID()
	businessID := suite.getTestBusinessID()
	req := createAuthRequest("GET", fmt.Sprintf("/api/bonus-programs/by-business/%d", businessID), nil, userID)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
	var programs []interface{}
	json.Unmarshal(recorder.Body.Bytes(), &programs)
	assert.IsType(suite.T(), []interface{}{}, programs)
}

func (suite *APITestSuite) TestBonusProgramGetPublicByBusinessID() {
	businessID := suite.getTestBusinessID()
	req := createRequest("GET", fmt.Sprintf("/api/bonus-programs/public/by-business/%d", businessID), nil)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
	var programs []interface{}
	json.Unmarshal(recorder.Body.Bytes(), &programs)
	assert.IsType(suite.T(), []interface{}{}, programs)
}

func (suite *APITestSuite) TestBonusProgramCreate() {
	userID := suite.getTestUserID()
	businessID := suite.getTestBusinessID()
	body := map[string]interface{}{
		"name":           "New Test Program",
		"description":    "Test description",
		"businessId":     businessID,
		"pointsRequired": 20,
		"discount":       25,
		"discountType":   "percentage",
	}
	req := createAuthRequest("POST", "/api/bonus-programs", body, userID)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 201, recorder.Code)
}

// ========== Category Routes Tests ==========

func (suite *APITestSuite) TestCategoryGetAll() {
	req := createRequest("GET", "/api/categories", nil)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

func (suite *APITestSuite) TestCategoryGetBySlug() {
	req := createRequest("GET", "/api/categories/coffee", nil)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

// ========== City Routes Tests ==========

func (suite *APITestSuite) TestCityGetAll() {
	req := createRequest("GET", "/api/cities", nil)
	recorder := executeRequest(suite.router, req)

	assert.Equal(suite.T(), 200, recorder.Code)
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}
