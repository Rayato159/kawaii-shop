package utils_tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/Rayato159/kawaii-shop/modules/users"
	"github.com/Rayato159/kawaii-shop/pkg/kawaiiauth"
	kawaiitests "github.com/Rayato159/kawaii-shop/tests"
)

type testParseToken struct {
	token  string
	isErr  bool
	expect string
}

func TestSignAccessToken(t *testing.T) {
	cfg := kawaiitests.Setup()

	tokenStack := make([]string, 0)
	tokenStack = append(tokenStack, "helloaccess")

	// Expires
	cfg.SetJwtAccessExpires(0)
	token, err := kawaiiauth.NewKawaiiAuth(kawaiiauth.Access, cfg.GetConfig().Jwt(), &users.UserClaims{
		Id: "U000001",
	})
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	if token == nil {
		t.Errorf("expect: %v, got: %v", "obj", token)
	}
	if token.SignToken() == "" {
		t.Errorf("expect: %v, got: %v", "xxxxxxxxx", token.SignToken())
	}
	tokenStack = append(tokenStack, token.SignToken())

	// Alive
	cfg.SetJwtAccessExpires(99999999)
	token, err = kawaiiauth.NewKawaiiAuth(kawaiiauth.Access, cfg.GetConfig().Jwt(), &users.UserClaims{
		Id: "U000001",
	})
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	if token == nil {
		t.Errorf("expect: %v, got: %v", "obj", token)
	}
	if token.SignToken() == "" {
		t.Errorf("expect: %v, got: %v", "xxxxxxxxx", token.SignToken())
	}
	tokenStack = append(tokenStack, token.SignToken())

	// Write file
	tokenJsonBytes, err := json.MarshalIndent(&tokenStack, "", "\t")
	if err != nil {
		t.Errorf("marshal indent token failed: %v", err)
	}
	if err := os.WriteFile("./access_token.json", tokenJsonBytes, 0777); err != nil {
		t.Errorf("export token failed: %v", err)
	}
}

func TestRefreshToken(t *testing.T) {
	cfg := kawaiitests.Setup()

	tokenStack := make([]string, 0)
	tokenStack = append(tokenStack, "hellorefresh")

	// Expires
	cfg.SetJwtRefreshExpires(0)
	token, err := kawaiiauth.NewKawaiiAuth(kawaiiauth.Refresh, cfg.GetConfig().Jwt(), &users.UserClaims{
		Id: "U000001",
	})
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	if token == nil {
		t.Errorf("expect: %v, got: %v", "obj", token)
	}
	if token.SignToken() == "" {
		t.Errorf("expect: %v, got: %v", "xxxxxxxxx", token.SignToken())
	}
	tokenStack = append(tokenStack, token.SignToken())

	// Alive
	cfg.SetJwtRefreshExpires(99999999)
	token, err = kawaiiauth.NewKawaiiAuth(kawaiiauth.Refresh, cfg.GetConfig().Jwt(), &users.UserClaims{
		Id: "U000001",
	})
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	if token == nil {
		t.Errorf("expect: %v, got: %v", "obj", token)
	}
	if token.SignToken() == "" {
		t.Errorf("expect: %v, got: %v", "xxxxxxxxx", token.SignToken())
	}
	tokenStack = append(tokenStack, token.SignToken())

	// Write file
	tokenJsonBytes, err := json.MarshalIndent(&tokenStack, "", "\t")
	if err != nil {
		t.Errorf("marshal indent token failed: %v", err)
	}
	if err := os.WriteFile("./refresh_token.json", tokenJsonBytes, 0777); err != nil {
		t.Errorf("export token failed: %v", err)
	}
}

func TestAdminToken(t *testing.T) {
	cfg := kawaiitests.Setup()

	tokenStack := make([]string, 0)
	tokenStack = append(tokenStack, "helloadmin")

	token, err := kawaiiauth.NewKawaiiAuth(kawaiiauth.Admin, cfg.GetConfig().Jwt(), nil)
	if err != nil {
		t.Errorf("expect: %v, got: %v", nil, err)
	}
	if token == nil {
		t.Errorf("expect: %v, got: %v", "obj", token)
	}
	if token.SignToken() == "" {
		t.Errorf("expect: %v, got: %v", "xxxxxxxxx", token.SignToken())
	}
	tokenStack = append(tokenStack, token.SignToken())

	// Write file
	tokenJsonBytes, err := json.MarshalIndent(&tokenStack, "", "\t")
	if err != nil {
		t.Errorf("marshal indent token failed: %v", err)
	}
	if err := os.WriteFile("./admin_token.json", tokenJsonBytes, 0777); err != nil {
		t.Errorf("export token failed: %v", err)
	}
}

func TestParseAccessToken(t *testing.T) {
	cfg := kawaiitests.Setup()
	tests := make([]testParseToken, 0)

	accessTokenJsonBytes, err := os.ReadFile("./access_token.json")
	if err != nil {
		t.Errorf("read file failed: %v", err)
	}
	accessToken := make([]string, 0)
	if err := json.Unmarshal(accessTokenJsonBytes, &accessToken); err != nil {
		t.Errorf("unmarshal access_token failed: %v", err)
	}
	tests = append(tests, testParseToken{
		token:  accessToken[0],
		isErr:  true,
		expect: "token format is invalid",
	})
	tests = append(tests, testParseToken{
		token:  accessToken[1],
		isErr:  true,
		expect: "token had expired",
	})
	tests = append(tests, testParseToken{
		token:  accessToken[2],
		isErr:  false,
		expect: "",
	})

	for _, req := range tests {
		if req.isErr {
			_, err := kawaiiauth.ParseToken(cfg.GetConfig().Jwt(), req.token)
			if err == nil {
				fmt.Println(req.token)
				t.Errorf("expect: %v, got: %v", "err", err)
			}
			if err != nil && err.Error() != req.expect {
				t.Errorf("expect: %v, got: %v", req.expect, err)
			}
		} else {
			_, err := kawaiiauth.ParseToken(cfg.GetConfig().Jwt(), req.token)
			if err != nil {
				t.Errorf("expect: %v, got: %v", nil, err)
			}
		}
	}
}
func TestParseRefreshToken(t *testing.T) {
	cfg := kawaiitests.Setup()
	tests := make([]testParseToken, 0)

	refreshTokenJsonBytes, err := os.ReadFile("./refresh_token.json")
	if err != nil {
		t.Errorf("read file failed: %v", err)
	}
	refreshToken := make([]string, 0)
	if err := json.Unmarshal(refreshTokenJsonBytes, &refreshToken); err != nil {
		t.Errorf("unmarshal refresh_token failed: %v", err)
	}
	tests = append(tests, testParseToken{
		token:  refreshToken[0],
		isErr:  true,
		expect: "token format is invalid",
	})
	tests = append(tests, testParseToken{
		token:  refreshToken[1],
		isErr:  true,
		expect: "token had expired",
	})
	tests = append(tests, testParseToken{
		token:  refreshToken[2],
		isErr:  false,
		expect: "",
	})

	for _, req := range tests {
		if req.isErr {
			_, err := kawaiiauth.ParseToken(cfg.GetConfig().Jwt(), req.token)
			if err == nil {
				fmt.Println(req.token)
				t.Errorf("expect: %v, got: %v", "err", err)
			}
			if err != nil && err.Error() != req.expect {
				t.Errorf("expect: %v, got: %v", req.expect, err)
			}
		} else {
			_, err := kawaiiauth.ParseToken(cfg.GetConfig().Jwt(), req.token)
			if err != nil {
				t.Errorf("expect: %v, got: %v", nil, err)
			}
		}
	}
}

func TestParseAdminToken(t *testing.T) {
	cfg := kawaiitests.Setup()
	tests := make([]testParseToken, 0)

	adminTokenJsonBytes, err := os.ReadFile("./admin_token.json")
	if err != nil {
		t.Errorf("read file failed: %v", err)
	}
	adminToken := make([]string, 0)
	if err := json.Unmarshal(adminTokenJsonBytes, &adminToken); err != nil {
		t.Errorf("unmarshal refresh_token failed: %v", err)
	}
	tests = append(tests, testParseToken{
		token:  adminToken[0],
		isErr:  true,
		expect: "token format is invalid",
	})
	tests = append(tests, testParseToken{
		token:  adminToken[1],
		isErr:  false,
		expect: "",
	})

	for _, req := range tests {
		if req.isErr {
			_, err := kawaiiauth.ParseAdminToken(cfg.GetConfig().Jwt(), req.token)
			if err == nil {
				fmt.Println(req.token)
				t.Errorf("expect: %v, got: %v", "err", err)
			}
			if err != nil && err.Error() != req.expect {
				t.Errorf("expect: %v, got: %v", req.expect, err)
			}
		} else {
			_, err := kawaiiauth.ParseAdminToken(cfg.GetConfig().Jwt(), req.token)
			if err != nil {
				t.Errorf("expect: %v, got: %v", nil, err)
			}
		}
	}
}

func TestRepeatToken(t *testing.T) {
	cfg := kawaiitests.Setup()
	token := kawaiiauth.RepeatToken(
		cfg.GetConfig().Jwt(),
		&users.UserClaims{
			Id: "U000001",
		},
		1777875301,
	)
	if token == "" {
		t.Errorf("expect: %v, got: %v", "xxxxxxxxxx", "")
	}
	fmt.Println(token)
}
