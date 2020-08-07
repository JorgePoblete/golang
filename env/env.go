package env

import (
	"os"
	"reflect"
	"strconv"
)

func Load(conf interface{}) {
	load(reflect.ValueOf(conf), "", "")
}

func load(conf reflect.Value, envTag, envDefault string) {
	// here conf could be either a struct or just a variable
	// if it's a variable we just set its value to the value of the
	// environment variable referenced by its tag, or its default, otherwise we recursively
	// set the struct value to the value returned by load(...) of each of its
	// individual fields

	if conf.Kind() == reflect.Ptr {
		reflectedConf := reflect.Indirect(conf)
		// we should only keep going if we can set values
		if reflectedConf.IsValid() && reflectedConf.CanSet() {
			value, ok := os.LookupEnv(envTag)
			// if the env variable is not set we just use the envDefault
			if !ok {
				value = envDefault
			}
			switch reflectedConf.Kind() {
			case reflect.Struct:
				for i := 0; i < reflectedConf.NumField(); i++ {
					if tag, ok := reflectedConf.Type().Field(i).Tag.Lookup("env"); ok {
						def, _ := reflectedConf.Type().Field(i).Tag.Lookup("envDefault")
						load(reflectedConf.Field(i).Addr(), envTag+tag, def)
					}
				}
				break
			// Here for each type we should make a cast of the env variable and then set the value
			case reflect.String:
				reflectedConf.SetString(value)
				break
			case reflect.Int:
				value, _ := strconv.Atoi(value)
				reflectedConf.Set(reflect.ValueOf(value))
				break
			case reflect.Bool:
				value, _ := strconv.ParseBool(value)
				reflectedConf.Set(reflect.ValueOf(value))
			}
		}

	}

}
