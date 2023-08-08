package testutil

import "testing"

func Setenv(t *testing.T, envs map[string]string) {
	for k, v := range envs {
		t.Setenv(k, v)
	}
}

// SetLocalMockEnv テスト環境向けに、環境変数を設定します
func SetLocalMockEnv(t *testing.T) {
	Setenv(t, map[string]string{
		"AWS_DEFAULT_REGION": "ap-northeast-1",
		"DDB_ENDPOINT":       "http://localhost:4566",
	})
}
