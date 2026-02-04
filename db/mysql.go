package db

import (
	"database/sql"
	"fmt"
	"log"
	"errors"
	"time"
	"golang.org/x/crypto/bcrypt"

	"crypto-member/config"
	"crypto-member/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

func Connect() {
	user := config.Get("DB_USER")
	pass := config.Get("DB_PASS")
	host := config.Get("DB_HOST")
	port := config.Get("DB_PORT")
	name := config.Get("DB_NAME")

	// DSN tanpa database untuk cek/create DB
	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&charset=utf8mb4&loc=Local", user, pass, host, port)

	sqlDBRoot, err := sql.Open("mysql", dsnRoot)
	if err != nil {
		log.Fatal("❌ Cannot connect to MySQL root:", err)
	}
	defer sqlDBRoot.Close()

	// Buat database kalau belum ada
	_, err = sqlDBRoot.Exec("CREATE DATABASE IF NOT EXISTS " + name + " CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;")
	if err != nil {
		log.Fatal("❌ Failed to create database:", err)
	}
	log.Println("✅ Database ready:", name)

	// DSN lengkap untuk GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local", user, pass, host, port, name)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Connected to database")

	// Auto migrate semua models
	err = DB.AutoMigrate(
		&models.User{},
		&models.DiscordCode{},
		&models.Payment{},
		&models.Coupon{},
		&models.ModuleGroup{},
		&models.Module{},
		&models.ModuleProgress{},
		&models.RulePricing{},
		&models.Expense{},
		&models.CryptoNews{},
	)
	if err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}

	log.Println("✅ Database migrated successfully")
	CreateDefaultAdmin()
	SeedRulePricing()
}

// Buat default admin kalau belum ada
func CreateDefaultAdmin() {
	var admin models.User
	err := DB.First(&admin, "role = ?", "admin").Error

	if err != nil && err == gorm.ErrRecordNotFound {
		// hash password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

		name := "Default Admin"
        expired := time.Now().AddDate(100, 0, 0)
		admin = models.User{
			Email:           "admin@example.com",
			Username:        "admin",
			Password:        string(hashedPassword),
			Role:            "admin",
	        NamaLengkap:     &name,
	        MemberExpiredAt: &expired,
		}

		if err := DB.Create(&admin).Error; err != nil {
			fmt.Println("Gagal buat default admin:", err)
		} else {
			fmt.Println("✅ Default admin berhasil dibuat!")
		}
	} else if err != nil {
		fmt.Println("Error cek default admin:", err)
	} else {
		fmt.Println("Default admin sudah ada:", admin.Username)
	}
}

func GetUserByToken(token string) (*models.User, error) {
	var user models.User
	if err := DB.Where("token = ?", token).First(&user).Error; err != nil {
		return nil, errors.New("user tidak ditemukan atau token salah")
	}
	return &user, nil
}

func UpdateUserDiscordID(userID uint, discordID string) error {
	return DB.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"nama_discord": discordID,
		}).Error
}

func RedeemDiscordCode(code string, discordName string, discordID string) (*models.User, *models.DiscordCode, *models.Payment, error) {
	var dcode models.DiscordCode
	err := DB.Preload("Payment").Where("code = ? AND is_used = false", code).First(&dcode); 
	if err.Error != nil {
		return nil, nil, nil, fmt.Errorf("Kode salah atau sudah digunakan :",err.Error)
	}

	fmt.Println("Kodeku : ", dcode)

	// ambil user via payment
	var payment models.Payment
	if err := DB.First(&payment, dcode.PaymentID).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("Payment tidak ditemukan")
	}

	var user models.User
	if err := DB.First(&user, payment.UserID).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("User tidak ditemukan")
	}
	
	fmt.Println("Mau cek Nama Discord Nih")

	if discordName == "" {
		return nil, nil, nil, fmt.Errorf("Nama Discord wajib diisi")
	}

	if discordID == "" {
		return nil, nil, nil, fmt.Errorf("ID Discord wajib diisi")
	}

	// update user dan discord code
	now := time.Now()
	user.NamaDiscord = &discordName
	user.IDDiscord = &discordID
	// cek apakah user punya MemberExpiredAt
	var newExpiry time.Time
	if user.MemberExpiredAt != nil && user.MemberExpiredAt.After(now) {
		newExpiry = user.MemberExpiredAt.AddDate(0, int(payment.MonthCount), 0)
	} else {
		newExpiry = now.AddDate(0, int(payment.MonthCount), 0)
	}
	user.MemberExpiredAt = &newExpiry
	if err := DB.Save(&user).Error; err != nil {
		return nil, nil, nil, fmt.Errorf("Gagal update user")
	}

	return &user, &dcode, &payment, nil
}

func SeedRulePricing() {
	// Hapus semua data lama
	// if err := DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.RulePricing{}).Error; err != nil {
	// 	log.Printf("Failed to clear RulePricing: %v", err)
	// 	return
	// }
	// log.Println("Semua data RulePricing dihapus ✅")

	var count int64

	if err := DB.Model(&models.RulePricing{}).
		Where("is_active = ?", true).
		Count(&count).Error; err != nil {
		log.Printf("Failed to check active RulePricing: %v", err)
		return
	}

	if count > 0 {
		log.Println("RulePricing aktif sudah ada, skip seeding ⏭️")
		return
	}

	pricings := []models.RulePricing{
		{MinMonth: 1, MaxMonth: intPtr(2), TotalPrice: 200_000, IsActive: true},
		{MinMonth: 3, MaxMonth: intPtr(5), TotalPrice: 550_000, IsActive: true},
		{MinMonth: 6, MaxMonth: intPtr(11), TotalPrice: 1_000_000, IsActive: true},
		{MinMonth: 12, MaxMonth: intPtr(999), TotalPrice: 1_800_000, IsActive: true},
		{MinMonth: 1000, MaxMonth: intPtr(9999), TotalPrice: 2_500_000, IsActive: true},
	}

	for _, p := range pricings {
		if err := DB.Create(&p).Error; err != nil {
			log.Printf("Failed to seed RulePricing: %v", err)
		}
	}

	log.Println("Seeder RulePricing selesai ✅")
}

func intPtr(i int) *int {
	return &i
}