package cli

import (
	"fmt"
)


// コマンドライン引数の解析とチェックの結果
type CheckResult struct {
	NoKeyValues []string          // オプションが指定されなかった引数
	KeyValues map[string][]string // オプションの値
	KeyFlags map[string]int32     // フラグオプションの回数
}

// コマンドライン引数の解析とチェック
func CheckArguments(args []string, settings OptionSettings) (CheckResult, error) {
	var minNoKeyValues int = 0
	var maxNoKeyValues int = 0
	var nowKey string = ""
	var nowSetting OptionKeySetting = OptionKeySetting{}
	var noKeyValues []string = []string{}
	var keyValues map[string][]string = map[string][]string{}
	var keyFlags map[string]int32 = map[string]int32{}
	
	// コマンドライン引数のオプション設定のチェック
	for _, setting := range settings.NoKeySettings {
		if setting.Required {
			minNoKeyValues++
		}
		maxNoKeyValues++
	}
	for i, setting := range settings.KeySettings {
		if setting.LongKey == "" && setting.ShortKey == "" {
			return CheckResult{}, fmt.Errorf("empty option is not allowed")
		}
		
		for _, setting2 := range settings.KeySettings[i+1:] {
			if setting.LongKey != "" {
				if setting2.LongKey != "" {
					if setting.LongKey == setting2.LongKey {
						return CheckResult{}, fmt.Errorf("duplicate option: %s", setting.LongKey)
					}
				}
				if setting2.ShortKey != "" {
					if setting.LongKey == setting2.ShortKey {
						return CheckResult{}, fmt.Errorf("duplicate option: %s", setting.LongKey)
					}
				}
			} else {
				if setting2.ShortKey != "" {
					if setting.ShortKey == setting2.ShortKey {
						return CheckResult{}, fmt.Errorf("duplicate option: %s", setting.ShortKey)
					}
				}
				if setting2.LongKey != "" {
					if setting.ShortKey == setting2.LongKey {
						return CheckResult{}, fmt.Errorf("duplicate option: %s", setting.ShortKey)
					}
				}
			}
		}
	}
	
	// コマンドライン引数の解析
	for _, arg := range args {
		// オプションのキーかどうか
		keyFound := false
		for _, setting := range settings.KeySettings {
			if fmt.Sprintf("--%s", setting.LongKey) == arg {
				keyFound = true
				nowKey = setting.LongKey
				nowSetting = setting
				break
			}
			if fmt.Sprintf("-%s", setting.ShortKey) == arg {
				keyFound = true
				nowKey = setting.LongKey
				nowSetting = setting
				break
			}
		}
		
		// 値の設定
		if keyFound {
			if nowSetting.Flag {
				_, ok := keyFlags[nowKey]
				if !ok { keyFlags[nowKey] = 0 }
				keyFlags[nowKey]++
				
				nowKey = ""
			}
		} else {
			if nowKey != "" {
				_, ok := keyValues[nowKey]
				if !ok { keyValues[nowKey] = []string{} }
				
				if !nowSetting.Multiple {
					if len(keyValues[nowKey]) > 0 {
						return CheckResult{}, fmt.Errorf("multiple values are not allowed for non-multiple option: %s", nowKey)
					}
				}
				keyValues[nowKey] = append(keyValues[nowKey], arg)
				nowKey = ""
			} else {
				noKeyValues = append(noKeyValues, arg)
			}
		}
	}
	
	// 解析後のオプション設定のチェック
	for _, setting := range settings.KeySettings {
		if setting.Flag {
			_, ok := keyFlags[setting.LongKey]
			if !ok { keyFlags[setting.LongKey] = 0 }
		} else {
			if setting.Required {
				v, ok := keyValues[setting.LongKey]
				if !ok {
					return CheckResult{}, fmt.Errorf("required option is not set: %s", setting.LongKey)
				}
				if len(v) == 0 {
					return CheckResult{}, fmt.Errorf("required option is not set: %s", setting.LongKey)
				}
			}
		}
	}
	if len(noKeyValues) < minNoKeyValues {
		return CheckResult{}, fmt.Errorf("too few options: %d", len(noKeyValues))
	}
	if len(noKeyValues) > maxNoKeyValues {
		return CheckResult{}, fmt.Errorf("too many options: %d", len(noKeyValues))
	}
	
	return CheckResult{
		NoKeyValues: noKeyValues,
		KeyValues: keyValues,
		KeyFlags: keyFlags,
	}, nil
}
