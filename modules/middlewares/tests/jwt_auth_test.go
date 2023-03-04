package middlewares_tests

import (
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/middlewares/repositories"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
)

type testFindAccessToken struct {
	userId      string
	accessToken string
	isErr       bool
	expect      bool
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
