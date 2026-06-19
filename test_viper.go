package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func main() {
	os.Setenv("MY_VAR", "") // Set to empty string
	
	viper.SetDefault("my.var", "default")
	viper.Set("my.var", "config_file_value") // simulate read from yaml
	
	viper.BindEnv("my.var", "MY_VAR")
	
	fmt.Printf("Viper value is: %q\n", viper.GetString("my.var"))
}
