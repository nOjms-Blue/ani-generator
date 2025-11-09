package settings

import (
	"fmt"
)


func PrintHelp(settings OptionSettings) {
	fmt.Println("Usage:")
	fmt.Printf("  ani-generator ")
	for _, setting := range settings.NoKeySettings {
		if (setting.Required) {
			fmt.Printf("%s ", setting.Name)
		} else {
			fmt.Printf("[%s] ", setting.Name)
		}
	}
	fmt.Printf("[options]\n")
	fmt.Println("Options:")
	for _, setting := range settings.KeySettings {
		if setting.LongKey != "" && setting.ShortKey != "" {
			fmt.Printf("  --%s, -%s\n", setting.LongKey, setting.ShortKey)
		} else if setting.LongKey != "" {
			fmt.Printf("  --%s\n", setting.LongKey)
		} else if setting.ShortKey != "" {
			fmt.Printf("  -%s\n", setting.ShortKey)
		}
		if setting.Flag {
			if setting.Description != "" { fmt.Printf("    %s\n", setting.Description) }
		} else {
			if setting.Description != "" { fmt.Printf("    %s\n", setting.Description) }
			if setting.Example != "" { fmt.Printf("    Example: %s\n", setting.Example) }
		}
		fmt.Println()
	}
}