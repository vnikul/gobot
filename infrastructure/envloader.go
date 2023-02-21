package infrastructure

import (
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/joho/godotenv"
	"gobot/entities"
	"reflect"
	"strconv"
)

func LoadConfigFromEnv() (entities.Config, error) {
	env, err := godotenv.Read(".env")
	if err != nil {
		return entities.Config{}, err
	}
	return load(env)
}

func load(configMap map[string]string) (entities.Config, error) {
	config := entities.Config{}

	structType := reflect.TypeOf(&config).Elem()
	structValue := reflect.ValueOf(&config).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		isRequired := true
		tag, ok := field.Tag.Lookup("required")
		if ok && tag == "false" {
			isRequired = false
		}
		value, ok := configMap[strcase.ToScreamingSnake(field.Name)]
		if ok == false {
			return entities.Config{}, errors.New(fmt.Sprintf("the field %s is required", field.Name))
		} else if value == "" && isRequired == true {
			return entities.Config{}, errors.New(fmt.Sprintf("the field %s should not be empty", field.Name))
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(value, 0, field.Type.Bits())
			if err != nil {
				return entities.Config{}, errors.New(fmt.Sprintf("could not parse field %s", field.Name))
			}
			fieldValue.SetInt(intVal)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(value)
			if err != nil {
				return entities.Config{}, errors.New(fmt.Sprintf("could not parse field %s", field.Name))
			}
			fieldValue.SetBool(boolVal)
		}
	}

	return config, nil
}
