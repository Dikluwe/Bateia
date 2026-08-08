/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/02_shell/cli/struct_flags.md
 * @layer L2
 * @updated 2026-08-08
 */

package cli

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
)

// ParseStructFlags mapeia declarativamente tags de struct em um FlagSet da CLI em L2.
func ParseStructFlags(target interface{}, args []string) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target deve ser um ponteiro para uma struct")
	}

	elem := val.Elem()
	t := elem.Type()
	fs := flag.NewFlagSet("struct_flags", flag.ContinueOnError)

	fieldMap := make(map[string]reflect.Value)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		flagName := field.Tag.Get("flag")
		if flagName == "" {
			continue
		}
		helpText := field.Tag.Get("help")
		fieldVal := elem.Field(i)

		switch field.Type.Kind() {
		case reflect.String:
			_ = fs.String(flagName, fieldVal.String(), helpText)
			fieldMap[flagName] = fieldVal
		case reflect.Int:
			_ = fs.Int(flagName, int(fieldVal.Int()), helpText)
			fieldMap[flagName] = fieldVal
		case reflect.Bool:
			_ = fs.Bool(flagName, fieldVal.Bool(), helpText)
			fieldMap[flagName] = fieldVal
		}
	}

	if err := fs.Parse(args); err != nil {
		return err
	}

	fs.Visit(func(f *flag.Flag) {
		if fVal, ok := fieldMap[f.Name]; ok {
			switch fVal.Kind() {
			case reflect.String:
				fVal.SetString(f.Value.String())
			case reflect.Int:
				if v, err := strconv.Atoi(f.Value.String()); err == nil {
					fVal.SetInt(int64(v))
				}
			case reflect.Bool:
				if v, err := strconv.ParseBool(f.Value.String()); err == nil {
					fVal.SetBool(v)
				}
			}
		}
	})

	return nil
}
