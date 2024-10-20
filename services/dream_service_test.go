package services_test

import (
	"context"
	"database/sql"
	"eatwo/db"
	"eatwo/models"
	"eatwo/services"
	"os"
	"testing"
	"time"
)

func getDreamService(t *testing.T) (*services.DreamService, func()) {
	sqlDB, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	dreamRepository := db.NewDreamRepository(sqlDB)
	err = dreamRepository.Migrate(context.Background())
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	dreamService := services.NewDreamService(dreamRepository)

	close := func() {
		defer sqlDB.Close()
		os.Remove(dbFileName)
	}

	return dreamService, close
}

func TestDreamService_Create(t *testing.T) {
	dreamService, close := getDreamService(t)
	defer close()

	err := dreamService.Create(&models.Dream{
		ID:          "1",
		UserID:      "1",
		Description: "description",
		Explanation: "explanation",
		Date:        time.Now(),
	})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}

func TestDreamService_GetByUserID(t *testing.T) {
	dreamService, close := getDreamService(t)
	defer close()
	seedDB(t, dreamService)

	dreams, err := dreamService.GetByUserID("22")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if len(dreams) != 1 {
		t.Errorf("Expected 1 dream, but got: %v", len(dreams))
	}
}

func TestDreamService_GetByUserID_OrderdBy_Date(t *testing.T) {
	dreamService, close := getDreamService(t)
	defer close()
	seedDB(t, dreamService)

	dreams, err := dreamService.GetByUserID("1")
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	if dreams[0].ID != "1" {
		t.Errorf("Expected dream with ID 1, but got: %v", dreams[0].ID)
	}

	if dreams[1].ID != "2" {
		t.Errorf("Expected dream with ID 2, but got: %v", dreams[1].ID)
	}
}

func seedDB(t *testing.T, dreamService *services.DreamService) {
	err := dreamService.Create(&models.Dream{
		ID:          "1",
		UserID:      "1",
		Description: "description",
		Explanation: "explanation",
		Date:        time.Now().AddDate(0, 0, -1),
	})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = dreamService.Create(&models.Dream{
		ID:          "2",
		UserID:      "1",
		Description: "description",
		Explanation: "explanation",
		Date:        time.Now(),
	})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	err = dreamService.Create(&models.Dream{
		ID:          "3",
		UserID:      "22",
		Description: "description",
		Explanation: "explanation",
		Date:        time.Now(),
	})
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
}
