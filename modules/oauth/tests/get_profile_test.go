package oauth_tests

import (
	"fmt"
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/oauth"
	"github.com/Rayato159/kawaii-shop/modules/oauth/repositories"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
)

type testGetProfileSuccess struct {
	userId string
	expect string
}

type testGetProfileFailed struct {
	userId string
	expect string
}

func TestGetProfile(t *testing.T) {
	db := kawaiitests.Setup().GetDb()
	r := repositories.OauthRepository(db)

	testsSuccess := testGetProfileSuccess{
		userId: "U000001",
		expect: func() string {
			profile := &oauth.User{
				Id:       "U000001",
				Email:    "customer001@kawaii.com",
				Username: "customer001",
				Role:     "customer",
			}
			return profile.ToJsonStringtify()
		}(),
	}

	profile, err := r.GetProfile(testsSuccess.userId)
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	if profile.ToJsonStringtify() != testsSuccess.expect {
		t.Errorf("expect: %v, got: %v", testsSuccess.expect, profile.ToJsonStringtify())
	}
	utils.Debug(profile)

	testsFailed := testGetProfileFailed{
		userId: "U999999",
		expect: "get user profile failed: sql: no rows in result set",
	}

	profile, err = r.GetProfile(testsFailed.userId)
	if err == nil {
		t.Errorf("expect: %v, got: %v", "err", nil)
	}
	if err.Error() != testsFailed.expect {
		t.Errorf("expect: %v, got: %v", testsFailed.expect, err.Error())
	}
	fmt.Println(err)
}
