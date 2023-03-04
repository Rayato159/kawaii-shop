package middlewares_tests

import (
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"
	"github.com/Rayato159/kawaii-shop/pkg/utils"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
)

type testFindAccessToken struct {
	userId      string
	accessToken string
	isErr       bool
	expect      bool
}

type testAuthorize struct {
	userRole   int
	expectRole int
	expect     int
}

func TestFindAccessToken(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	middlewareRepo := repositories.MiddlewareRepository(db)

	tests := []testFindAccessToken{
		{
			userId:      "xxxxxxxxxxxx",
			accessToken: "xxxxxxxxxxx",
			isErr:       true,
			expect:      false,
		},
		{
			userId:      "U000001",
			accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGFpbXMiOnsiaWQiOiJVMDAwMDAxIn0sImlzcyI6Imthd2FpaXNob3AtYXBpIiwic3ViIjoiYWNjZXNzLXRva2VuIiwiYXVkIjpbImN1c3RvbWVyIiwiYWRtaW4iXSwiZXhwIjoxNjc4MDI4MTI0LCJuYmYiOjE2Nzc5NDE3MjQsImlhdCI6MTY3Nzk0MTcyNH0.5Gu0fMyFFNV3QAGaQslnigP7N-GQWVZQi-9TkmDy4Fw",
			isErr:       false,
			expect:      true,
		},
	}

	for _, req := range tests {
		check := middlewareRepo.FindAccessToken(req.userId, req.accessToken)
		if req.isErr {
			if check {
				t.Errorf("expect: %v, got: %v", req.expect, check)
			}
		} else {
			if !check {
				t.Errorf("expect: %v, got: %v", req.expect, check)
			}
		}
	}

}

func TestBinaryConverter(t *testing.T) {
	number := 1 + 2
	bits := 3
	expect := []int{1, 1}

	binary := utils.BinaryConverter(number, bits)
	for i := range binary {
		if binary[i] != expect[i] {
			t.Errorf("expect: %v, got: %v", expect, binary)
		}
	}
}

func TestFindRole(t *testing.T) {
	db := kawaiitests.Setup().GetDb()

	middlewareRepo := repositories.MiddlewareRepository(db)
	roles, err := middlewareRepo.FindRole()
	if err != nil {
		t.Errorf("expect: %v, err: %v", nil, err)
	}
	utils.Debug(roles)
}

func TestAuthorize(t *testing.T) {
	length := 2
	tests := []testAuthorize{
		{
			userRole:   1,
			expectRole: 2,
			expect:     0,
		},
		{
			userRole:   2,
			expectRole: 2,
			expect:     1,
		},
		{
			userRole:   2,
			expectRole: 3,
			expect:     1,
		},
	}

	for _, test := range tests {
		expectValueBinary := utils.BinaryConverter(test.expectRole, length)
		userValueBinary := utils.BinaryConverter(test.userRole, length)

		for j := range userValueBinary {
			if userValueBinary[j]&expectValueBinary[j] == 1 {
				t.Errorf("expect: %v, got: %v", test.expect, userValueBinary[j]&expectValueBinary[j])
			}
		}
	}
}
