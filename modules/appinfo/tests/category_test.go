package appinfo_tests

import (
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/Rayato159/kawaii-shop/modules/appinfo/repositories"
	"github.com/Rayato159/kawaii-shop/modules/appinfo/usecases"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
)

type testFindCateogry struct {
	req *appinfo.CategoryFilter
}

func TestFindCategory(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	appinfoRepo := repositories.AppinfoRepository(db)

	tests := []testFindCateogry{
		{
			req: &appinfo.CategoryFilter{
				Title: "",
			},
		},
		{
			req: &appinfo.CategoryFilter{
				Title: "fashion",
			},
		},
	}

	for _, test := range tests {
		category, err := appinfoRepo.FindCategory(test.req)
		if err != nil {
			t.Errorf("expect: %v, got: %v", "category", category)
		}
		utils.Debug(category)
	}
}

func TestInsertCategory(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	appinfoRepo := repositories.AppinfoRepository(db)
	tests := []*appinfo.Category{
		{
			Title: "gundam",
		},
		{
			Title: "vehicle",
		},
		{
			Title: "game",
		},
	}

	if err := appinfoRepo.InsertCategory(tests); err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
}

func TestInsertCategoryUC(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	appinfoRepo := repositories.AppinfoRepository(db)
	appinfoUC := usecases.AppinfoUsecase(appinfoRepo)
	tests := []*appinfo.Category{
		{
			Title: "gundamx",
		},
		{
			Title: "vehiclex",
		},
		{
			Title: "gamex",
		},
	}

	category, err := appinfoUC.InsertCategory(tests)
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}

	if category == nil {
		t.Errorf("expect: %v, got: %v", "category", category)
		return
	}
	if len(category) == 0 {
		t.Errorf("expect: %v, got: %v", "gt 0", len(category))
	}
	utils.Debug(category)
}

func TestDeleteCategory(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	appinfoRepo := repositories.AppinfoRepository(db)
	categoryId := 5

	if err := appinfoRepo.DeleteCategory(categoryId); err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
}
