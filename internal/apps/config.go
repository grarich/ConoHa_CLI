/*
Copyright © 2023 grarich <grarich@grawlily.com>
*/
package apps

import (
	"fmt"
	"github.com/spf13/viper"
)

func GetConfig(key string) string {
	return viper.GetString(key)
}

func SetConfig(key string, value string) error {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("設定ファイルの書き込みに失敗しました。")
	}
	return nil
}
