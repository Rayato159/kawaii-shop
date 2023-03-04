package appinfo_tests

import (
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/appinfo/repositories"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
)

func TestFindCategory(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	appinfoRepo := repositories.AppinfoRepository(db)
	category, err := appinfoRepo.FindCategory()
	if err != nil {
		t.Errorf("expect: %v, got: %v", "category", category)
	}
	utils.Debug(category)
}
