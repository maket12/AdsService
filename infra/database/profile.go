package database

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"log"
	"time"
)

type Profile struct {
	UserID               uint64         `gorm:"primaryKey"`
	Name                 string         `gorm:"type:text"`
	Phone                string         `gorm:"type:text"`
	PhotoID              string         `gorm:"type:text"`
	NotificationsEnabled bool           `gorm:"default:true"`
	Subscriptions        pq.StringArray `gorm:"type:jsonb;serializer:json;not null;default:'[]'::jsonb"`
	UpdatedAt            time.Time      `gorm:"not null"`
	Banned               bool           `gorm:"default:false"`
}

func AddProfile(userID uint64, name, phone string) (*Profile, error) {
	var p = Profile{
		UserID:    userID,
		Name:      name,
		Phone:     phone,
		UpdatedAt: time.Now().UTC(),
	}
	log.Printf("Created new profile: %v", userID)
	return &p, DB.Create(&p).Error
}

func UpdateProfileName(userID uint64, name string) (*Profile, error) {
	if result := DB.Model(Profile{}).Where("user_id = ?", userID).
		Update("name", name).Error; result != nil {
		return nil, result
	}
	log.Printf("Updated name = %v of user = %v", name, userID)
	return UpdateProfileTime(userID)
}

func UpdateProfilePhone(userID uint64, phone string) (*Profile, error) {
	if result := DB.Model(Profile{}).Where("user_id = ?", userID).
		Update("phone", phone).Error; result != nil {
		return nil, result
	}
	log.Printf("Updated phone = %v of user = %v", phone, userID)
	return UpdateProfileTime(userID)
}

func UpdateProfilePhoto(userID uint64, photoID string) (*Profile, error) {
	if result := DB.Model(Profile{}).Where("user_id = ?", userID).
		Update("photo_id", photoID).Error; result != nil {
		return nil, result
	}
	log.Printf("Updated photo = %v of user = %v", photoID, userID)
	return UpdateProfileTime(userID)
}

func EnableNotifications(userID uint64) (*Profile, error) {
	if result := DB.Model(Profile{}).Where("user_id = ?", userID).
		Update("notifications_enabled", true).Error; result != nil {
		return nil, result
	}
	log.Printf("Enable notifications for user = %v", userID)
	return UpdateProfileTime(userID)
}

func DisableNotifications(userID uint64) (*Profile, error) {
	if result := DB.Model(Profile{}).Where("user_id = ?", userID).
		Update("notifications_enabled", false).Error; result != nil {
		return nil, result
	}
	log.Printf("Disable notifications for user = %v", userID)
	return UpdateProfileTime(userID)
}

func UpdateProfileSubscriptions(userID uint64, subscriptions []string) (*Profile, error) {
	if result := DB.Model(Profile{}).Where("user_id = ?", userID).
		Update("subscriptions", subscriptions).Error; result != nil {
		return nil, result
	}
	log.Printf("Updated subscriptions = %v of user = %v", subscriptions, userID)
	return UpdateProfileTime(userID)
}

func UpdateProfileTime(userID uint64) (*Profile, error) {
	var p Profile
	result := DB.Model(Profile{}).Where("user_id = ?", userID).
		Update("updated_at", gorm.Expr("now()")).First(&p)
	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
}

func GetProfile(userID uint64) (*Profile, error) {
	var p Profile
	result := DB.Model(Profile{}).Where("user_id = ?", userID).First(&p)
	if result.Error != nil {
		return nil, result.Error
	}
	return &p, nil
}

func GetAllProfiles() ([]Profile, error) {
	var profiles []Profile
	result := DB.Find(&profiles)
	if result.Error != nil {
		return nil, result.Error
	}
	return profiles, nil
}
