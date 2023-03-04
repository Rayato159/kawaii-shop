package users_tests

import (
	"fmt"
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/modules/users/repositories"
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
	r := repositories.UsersRepository(db)

	testsSuccess := testGetProfileSuccess{
		userId: "U000001",
		expect: func() string {
			profile := &users.User{
				Id:       "U000001",
				Email:    "customer001@kawaii.com",
				Username: "customer001",
				Role:     "customer",
			}
			return kawaiitests.ToJsonStringtify(profile)
		}(),
	}

	profile, err := r.GetProfile(testsSuccess.userId)
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	if kawaiitests.ToJsonStringtify(profile) != testsSuccess.expect {
		t.Errorf("expect: %v, got: %v", testsSuccess.expect, kawaiitests.ToJsonStringtify(profile))
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
