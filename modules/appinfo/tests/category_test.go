package appinfo_tests

import (
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/appinfo"
	"github.com/Rayato159/kawaii-shop/modules/appinfo/repositories"
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
