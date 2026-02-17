package database

import (
	"log"
	"time"

	"github.com/stempo/backend/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	log.Println("Starting database seeding...")

	if err := seedCities(db); err != nil {
		return err
	}

	if err := seedCategories(db); err != nil {
		return err
	}

	if err := seedUsers(db); err != nil {
		return err
	}

	if err := seedBusinesses(db); err != nil {
		return err
	}

	if err := seedBonusPrograms(db); err != nil {
		return err
	}

	if err := seedLoyaltyCards(db); err != nil {
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}

func seedCities(db *gorm.DB) error {
	var count int64
	db.Model(&entity.City{}).Count(&count)
	if count > 0 {
		log.Println("Cities already seeded, skipping...")
		return nil
	}

	kyivLat := 50.4501
	kyivLng := 30.5234
	lvivLat := 49.8397
	lvivLng := 24.0297
	kharkivLat := 49.9935
	kharkivLng := 36.2304
	odesaLat := 46.4825
	odesaLng := 30.7233

	cities := []entity.City{
		{
			Name:      "Kyiv",
			Slug:      "kyiv",
			Latitude:  &kyivLat,
			Longitude: &kyivLng,
		},
		{
			Name:      "Lviv",
			Slug:      "lviv",
			Latitude:  &lvivLat,
			Longitude: &lvivLng,
		},
		{
			Name:      "Kharkiv",
			Slug:      "kharkiv",
			Latitude:  &kharkivLat,
			Longitude: &kharkivLng,
		},
		{
			Name:      "Odesa",
			Slug:      "odesa",
			Latitude:  &odesaLat,
			Longitude: &odesaLng,
		},
		{
			Name:      "Dnipro",
			Slug:      "dnipro",
		},
		{
			Name:      "Zaporizhzhia",
			Slug:      "zaporizhzhia",
		},
		{
			Name:      "Vinnytsia",
			Slug:      "vinnytsia",
		},
		{
			Name:      "Poltava",
			Slug:      "poltava",
		},
		{
			Name:      "Chernivtsi",
			Slug:      "chernivtsi",
		},
		{
			Name:      "Ivano-Frankivsk",
			Slug:      "ivano-frankivsk",
		},
	}

	for i := range cities {
		cities[i].CreatedAt = time.Now()
		cities[i].UpdatedAt = time.Now()
		if err := db.Create(&cities[i]).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d cities", len(cities))
	return nil
}

func seedCategories(db *gorm.DB) error {
	var count int64
	db.Model(&entity.Category{}).Count(&count)
	if count > 0 {
		log.Println("Categories already seeded, skipping...")
		return nil
	}

	categories := []entity.Category{
		{
			Name:      "Coffee",
			Slug:      "coffee",
			Icon:      "☕",
			IconColor: "#8B5CF6",
		},
		{
			Name:      "Sports",
			Slug:      "sports",
			Icon:      "💪",
			IconColor: "#EF4444",
		},
		{
			Name:      "Courses",
			Slug:      "courses",
			Icon:      "💻",
			IconColor: "#6366F1",
		},
		{
			Name:      "Bar",
			Slug:      "bar",
			Icon:      "🍺",
			IconColor: "#F59E0B",
		},
		{
			Name:      "Food",
			Slug:      "food",
			Icon:      "🍝",
			IconColor: "#F59E0B",
		},
	}

	for i := range categories {
		categories[i].CreatedAt = time.Now()
		categories[i].UpdatedAt = time.Now()
		if err := db.Create(&categories[i]).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d categories", len(categories))
	return nil
}

func seedUsers(db *gorm.DB) error {
	var count int64
	db.Model(&entity.User{}).Count(&count)
	if count > 0 {
		log.Println("Users already seeded, skipping...")
		return nil
	}

	hashedPassword, err := hashPassword("admin123")
	if err != nil {
		return err
	}

	users := []entity.User{
		{
			Email:    "admin@example.com",
			Password: hashedPassword,
			Name:     stringPtr("Admin User"),
			Phone:    stringPtr("+1 (555) 000-0000"),
		},
		{
			Email:    "user@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			Name:     stringPtr("John Doe"),
			Phone:    stringPtr("+1 (555) 123-4567"),
		},
		{
			Email:    "business@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
			Name:     stringPtr("Business Owner"),
			Phone:    stringPtr("+1 (555) 987-6543"),
		},
	}

	for i := range users {
		users[i].CreatedAt = time.Now()
		users[i].UpdatedAt = time.Now()
		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d users", len(users))
	return nil
}

func seedBusinesses(db *gorm.DB) error {
	var count int64
	db.Model(&entity.Business{}).Count(&count)
	if count > 0 {
		log.Println("Businesses already seeded, skipping...")
		return nil
	}

	var owner entity.User
	if err := db.Where("email = ?", "business@example.com").First(&owner).Error; err != nil {
		return err
	}

	businesses := []entity.Business{
		{
			Name:            "Joe's Coffee",
			Category:        "coffee",
			Address:         "123 Main St, Downtown",
			Rating:          4.8,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1554118811-1e0d58224f24?w=400&h=400&fit=crop"),
			Icon:            "☕",
			IconColor:       "#8B5CF6",
			Description:     stringPtr("Artisan coffee and fresh pastries"),
			HasLoyaltyProgram: true,
			Featured:        true,
			Status:          entity.BusinessStatusApproved,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Zen Yoga",
			Category:        "sports",
			Address:         "456 Oak Ave, Midtown",
			Rating:          4.9,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1506126613408-eca07ce68773?w=400&h=400&fit=crop"),
			Icon:            "🧘",
			IconColor:       "#10B981",
			Description:     stringPtr("Peaceful yoga classes for all levels"),
			HasLoyaltyProgram: true,
			Featured:        true,
			Status:          entity.BusinessStatusApproved,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Bella Pasta",
			Category:        "food",
			Address:         "789 Pine St, Uptown",
			Rating:          4.7,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1414235077428-338989a2e8c0?w=400&h=400&fit=crop"),
			Icon:            "🍝",
			IconColor:       "#F59E0B",
			Description:     stringPtr("Authentic Italian cuisine"),
			HasLoyaltyProgram: false,
			Featured:        true,
			Status:          entity.BusinessStatusApproved,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Green Smoothie Bar",
			Category:        "food",
			Address:         "321 Elm St, Downtown",
			Rating:          4.6,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1600271886742-f049cd451bba?w=400&h=400&fit=crop"),
			Icon:            "🥤",
			IconColor:       "#10B981",
			Description:     stringPtr("Fresh juices and healthy smoothies"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Bookworm Café",
			Category:        "coffee",
			Address:         "654 University Blvd",
			Rating:          4.5,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1481627834876-b7833e8f5570?w=400&h=400&fit=crop"),
			Icon:            "📚",
			IconColor:       "#6366F1",
			Description:     stringPtr("Books, coffee, and cozy reading nooks"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Fit Zone Gym",
			Category:        "sports",
			Address:         "890 Fitness Ave",
			Rating:          4.4,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1534438327276-14e5300c3a48?w=400&h=400&fit=crop"),
			Icon:            "💪",
			IconColor:       "#EF4444",
			Description:     stringPtr("Modern fitness equipment and classes"),
			HasLoyaltyProgram: false,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Brew & Bean",
			Category:        "coffee",
			Address:         "234 Coffee Lane",
			Rating:          4.7,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1509042239860-f550ce710b93?w=400&h=400&fit=crop"),
			Icon:            "☕",
			IconColor:       "#8B5CF6",
			Description:     stringPtr("Specialty coffee roasters"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "CrossFit Central",
			Category:        "sports",
			Address:         "567 Workout St",
			Rating:          4.6,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1534438327276-14e5300c3a48?w=400&h=400&fit=crop"),
			Icon:            "🏋️",
			IconColor:       "#EF4444",
			Description:     stringPtr("High-intensity training"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Coding Academy",
			Category:        "courses",
			Address:         "890 Tech Blvd",
			Rating:          4.9,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=400&h=400&fit=crop"),
			Icon:            "💻",
			IconColor:       "#6366F1",
			Description:     stringPtr("Programming and web development courses"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "The Craft Bar",
			Category:        "bar",
			Address:         "123 Bar Street",
			Rating:          4.5,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1514933651103-005eec06c04b?w=400&h=400&fit=crop"),
			Icon:            "🍺",
			IconColor:       "#F59E0B",
			Description:     stringPtr("Craft cocktails and local beers"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Sushi Master",
			Category:        "food",
			Address:         "456 Sushi Ave",
			Rating:          4.8,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1579584425555-c3ce17fd4351?w=400&h=400&fit=crop"),
			Icon:            "🍣",
			IconColor:       "#EF4444",
			Description:     stringPtr("Fresh sushi and Japanese cuisine"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Espresso Corner",
			Category:        "coffee",
			Address:         "789 Bean Road",
			Rating:          4.6,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1442512595331-e89e73853f31?w=400&h=400&fit=crop"),
			Icon:            "☕",
			IconColor:       "#8B5CF6",
			Description:     stringPtr("Quick espresso and pastries"),
			HasLoyaltyProgram: false,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Martial Arts Dojo",
			Category:        "sports",
			Address:         "234 Fight St",
			Rating:          4.7,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1549060279-7e168fce2090?w=400&h=400&fit=crop"),
			Icon:            "🥋",
			IconColor:       "#EF4444",
			Description:     stringPtr("Karate, Judo, and Taekwondo classes"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Language School",
			Category:        "courses",
			Address:         "567 Learn Ave",
			Rating:          4.8,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1434030216411-0b793f4b4173?w=400&h=400&fit=crop"),
			Icon:            "📖",
			IconColor:       "#6366F1",
			Description:     stringPtr("English, Spanish, French courses"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Wine & Dine",
			Category:        "bar",
			Address:         "890 Vineyard Lane",
			Rating:          4.6,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1510812431401-41d2bd2722f3?w=400&h=400&fit=crop"),
			Icon:            "🍷",
			IconColor:       "#F59E0B",
			Description:     stringPtr("Fine wines and gourmet food"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Burger House",
			Category:        "food",
			Address:         "123 Burger Blvd",
			Rating:          4.5,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=400&h=400&fit=crop"),
			Icon:            "🍔",
			IconColor:       "#F59E0B",
			Description:     stringPtr("Gourmet burgers and fries"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Morning Brew",
			Category:        "coffee",
			Address:         "456 Sunrise Ave",
			Rating:          4.7,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1501339847302-ac426a4c7c6e?w=400&h=400&fit=crop"),
			Icon:            "☕",
			IconColor:       "#8B5CF6",
			Description:     stringPtr("Breakfast and coffee"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Swimming Pool",
			Category:        "sports",
			Address:         "789 Aqua St",
			Rating:          4.5,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1576610616656-d3aa5d1f4534?w=400&h=400&fit=crop"),
			Icon:            "🏊",
			IconColor:       "#10B981",
			Description:     stringPtr("Swimming lessons and pool access"),
			HasLoyaltyProgram: false,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Music Academy",
			Category:        "courses",
			Address:         "234 Melody Road",
			Rating:          4.8,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=400&h=400&fit=crop"),
			Icon:            "🎵",
			IconColor:       "#6366F1",
			Description:     stringPtr("Piano, guitar, and vocal lessons"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Cocktail Lounge",
			Category:        "bar",
			Address:         "567 Mix St",
			Rating:          4.6,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1551538827-9c037cb4f32a?w=400&h=400&fit=crop"),
			Icon:            "🍸",
			IconColor:       "#F59E0B",
			Description:     stringPtr("Creative cocktails and live music"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Pizza Place",
			Category:        "food",
			Address:         "890 Slice Ave",
			Rating:          4.6,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1513104890138-7c749659a591?w=400&h=400&fit=crop"),
			Icon:            "🍕",
			IconColor:       "#F59E0B",
			Description:     stringPtr("Wood-fired pizza and Italian dishes"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Coffee Lab",
			Category:        "coffee",
			Address:         "123 Experiment St",
			Rating:          4.9,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1447933601403-0c6688de566e?w=400&h=400&fit=crop"),
			Icon:            "☕",
			IconColor:       "#8B5CF6",
			Description:     stringPtr("Experimental coffee brewing"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Rock Climbing",
			Category:        "sports",
			Address:         "234 Climb Blvd",
			Rating:          4.7,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1516567727245-ad8c68f3ec93?w=400&h=400&fit=crop"),
			Icon:            "🧗",
			IconColor:       "#EF4444",
			Description:     stringPtr("Indoor rock climbing facility"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Art Studio",
			Category:        "courses",
			Address:         "567 Canvas Road",
			Rating:          4.7,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?w=400&h=400&fit=crop"),
			Icon:            "🎨",
			IconColor:       "#6366F1",
			Description:     stringPtr("Painting and drawing classes"),
			HasLoyaltyProgram: true,
			Featured:        true,
			OwnerID:         owner.ID,
		},
		{
			Name:            "Beer Garden",
			Category:        "bar",
			Address:         "890 Hops Lane",
			Rating:          4.5,
			IsOpen:          true,
			ImageURL:        stringPtr("https://images.unsplash.com/photo-1551218808-94e220e084d2?w=400&h=400&fit=crop"),
			Icon:            "🍻",
			IconColor:       "#F59E0B",
			Description:     stringPtr("Outdoor beer garden"),
			HasLoyaltyProgram: false,
			Featured:        true,
			OwnerID:         owner.ID,
		},
	}

	for i := range businesses {
		businesses[i].CreatedAt = time.Now()
		businesses[i].UpdatedAt = time.Now()
		// Set status to approved if not already set
		if businesses[i].Status == "" {
			businesses[i].Status = entity.BusinessStatusApproved
		}
		if err := db.Create(&businesses[i]).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d businesses", len(businesses))
	return nil
}

func seedBonusPrograms(db *gorm.DB) error {
	var count int64
	db.Model(&entity.BonusProgram{}).Count(&count)
	if count > 0 {
		log.Println("Bonus programs already seeded, skipping...")
		return nil
	}

	var joesCoffee, zenYoga, bellaPasta, smoothieBar, bookworm entity.Business
	if err := db.Where("name = ?", "Joe's Coffee").First(&joesCoffee).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Zen Yoga").First(&zenYoga).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Bella Pasta").First(&bellaPasta).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Green Smoothie Bar").First(&smoothieBar).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Bookworm Café").First(&bookworm).Error; err != nil {
		return err
	}

	expiresAt := time.Now().AddDate(1, 0, 0)

	programs := []entity.BonusProgram{
		{
			Name:           "Coffee Lovers",
			Status:         entity.BonusProgramStatusApproved,
			Description:   "Get 1 stamp for every coffee purchase. 10 stamps = free coffee!",
			BusinessID:     joesCoffee.ID,
			Points:         7,
			PointsRequired: 10,
			Discount:       100,
			DiscountType:   "percentage",
			ExpiresAt:      &expiresAt,
		},
		{
			Name:           "Loyalty Rewards",
			Description:   "Collect 5 stamps and get 20% off your next order",
			BusinessID:     joesCoffee.ID,
			Points:         3,
			PointsRequired: 5,
			Discount:       20,
			DiscountType:   "percentage",
			Status:         entity.BonusProgramStatusActive,
			ExpiresAt:      &expiresAt,
		},
		{
			Name:           "Yoga Pass",
			Description:   "Attend 10 classes and get 1 free class",
			BusinessID:     zenYoga.ID,
			Points:         5,
			PointsRequired: 10,
			Discount:       100,
			DiscountType:   "percentage",
			Status:         entity.BonusProgramStatusActive,
			ExpiresAt:      &expiresAt,
		},
		{
			Name:           "Pasta Lovers",
			Description:   "Order 5 pasta dishes, get 20% off your next visit",
			BusinessID:     bellaPasta.ID,
			Points:         2,
			PointsRequired: 5,
			Discount:       20,
			DiscountType:   "percentage",
			Status:         entity.BonusProgramStatusActive,
			ExpiresAt:      &expiresAt,
		},
		{
			Name:           "Smoothie Rewards",
			Description:   "Buy 8 smoothies, get the 9th one free",
			BusinessID:     smoothieBar.ID,
			Points:         4,
			PointsRequired: 8,
			Discount:       100,
			DiscountType:   "percentage",
			Status:         entity.BonusProgramStatusActive,
			ExpiresAt:      &expiresAt,
		},
		{
			Name:           "Book & Brew",
			Description:   "Buy 5 books and get 15% off your next purchase",
			BusinessID:     bookworm.ID,
			Points:         2,
			PointsRequired: 5,
			Discount:       15,
			DiscountType:   "percentage",
			Status:         entity.BonusProgramStatusActive,
			ExpiresAt:      &expiresAt,
		},
	}

	for i := range programs {
		programs[i].CreatedAt = time.Now()
		programs[i].UpdatedAt = time.Now()
		if err := db.Create(&programs[i]).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d bonus programs", len(programs))
	return nil
}

func seedLoyaltyCards(db *gorm.DB) error {
	var count int64
	db.Model(&entity.LoyaltyCard{}).Count(&count)
	if count > 0 {
		log.Println("Loyalty cards already seeded, skipping...")
		return nil
	}

	var user entity.User
	if err := db.Where("email = ?", "user@example.com").First(&user).Error; err != nil {
		return err
	}

	var business1, business2 entity.Business
	if err := db.Where("name = ?", "Joe's Coffee").First(&business1).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Zen Yoga").First(&business2).Error; err != nil {
		return err
	}

	cards := []entity.LoyaltyCard{
		{
			UserID:         user.ID,
			BusinessID:     business1.ID,
			Stamps:         8,
			StampsRequired: 10,
		},
		{
			UserID:         user.ID,
			BusinessID:     business2.ID,
			Stamps:         5,
			StampsRequired: 10,
		},
	}

	for i := range cards {
		cards[i].CreatedAt = time.Now()
		cards[i].UpdatedAt = time.Now()
		if err := db.Create(&cards[i]).Error; err != nil {
			return err
		}
	}

	log.Printf("Seeded %d loyalty cards", len(cards))
	return nil
}

func stringPtr(s string) *string {
	return &s
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
