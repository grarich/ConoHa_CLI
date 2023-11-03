/*
Copyright © 2023 grarich <grarich@grawlily.com>
*/
package identity

import (
	"bytes"
	"conoha_cli/internal/apps"
	"conoha_cli/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"time"
)

// Request
type passwordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type auth struct {
	PasswordCredentials passwordCredentials `json:"passwordCredentials"`
}

type authBody struct {
	Auth     auth   `json:"auth"`
	TenantID string `json:"tenantId"`
}

// Response
type token struct {
	IssuedAt string    `json:"issued_at"`
	Expires  time.Time `json:"expires"`
	ID       string    `json:"id"`
	AuditIds []string  `json:"audit_ids"`
}

type user struct {
	UserName   string   `json:"username"`
	RolesLinks []string `json:"roles_links"`
	ID         string   `json:"id"`
	Roles      []string `json:"roles"`
	Name       string   `json:"name"`
}

type access struct {
	Token          token    `json:"token"`
	ServiceCatalog []string `json:"serviceCatalog"`
	User           user     `json:"user"`
}

type metadata struct {
	IsAdmin int      `json:"is_admin"`
	Roles   []string `json:"roles"`
}

type tokenResponse struct {
	Access   access   `json:"access"`
	Metadata metadata `json:"metadata"`
}

var (
	uname       string
	GetTokenCmd = &cobra.Command{
		Use:   "get-token",
		Short: "トークン発行を行います。",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			endpoint := apps.GetConfig("endpoint.identity")
			if endpoint == "" {
				return fmt.Errorf("エンドポイントが設定されていません。")
			}
			cmd.Print("Password: ")
			pwd, err := utils.ReadPassword()
			fmt.Println()
			if err != nil {
				return err
			}
			tokenResp, err := callGetTokenAPI(endpoint, string(pwd))
			if err != nil {
				return err
			}

			// 設定ファイル書き込み
			if err := apps.SetConfig("identity.token", tokenResp.Access.Token.ID); err != nil {
				return err
			}
			if err := apps.SetConfig("identity.exp", tokenResp.Access.Token.Expires.Format(time.RFC3339)); err != nil {
				return err
			}
			cmd.Println("設定ファイルにトークンを保存しました。")
			return nil
		},
	}
)

func init() {
	GetTokenCmd.Flags().StringVar(&uname, "username", "", "ユーザー名を指定してください")
	if err := GetTokenCmd.MarkFlagRequired("username"); err != nil {
		return
	}
}

func callGetTokenAPI(endpoint string, pwd string) (*tokenResponse, error) {
	body := new(authBody)
	body.TenantID = ""
	body.Auth.PasswordCredentials.Username = uname
	body.Auth.PasswordCredentials.Password = pwd
	reqBody, _ := json.Marshal(body)
	resp, err := http.Post(fmt.Sprintf("%-v/v2.0/tokens", endpoint), "application/json", bytes.NewReader(reqBody))
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			return
		}
	}(resp.Body)

	if err != nil {
		return new(tokenResponse), err
	}
	if resp.StatusCode == 200 {
		var d tokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
			return new(tokenResponse), err
		}
		return &d, nil
	} else if resp.StatusCode == 401 {
		return new(tokenResponse), fmt.Errorf("認証に失敗しました。ユーザー名またはパスワードが間違っています。")
	}
	return new(tokenResponse), fmt.Errorf("想定外のレスポンスが返されました。StatusCodee: %-v", resp.StatusCode)
}
