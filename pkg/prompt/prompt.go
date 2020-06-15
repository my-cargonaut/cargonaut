package prompt

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structtag"
	"github.com/manifoldco/promptui"
)

// Available struct tag values for "prompt".
const (
	TagValueSkip     = "-"
	TagValueRequired = "required"
	TagValueValidate = "validate"
	TagValueName     = "name"
	TagValueMask     = "mask"
)

// ValidatorFunc validates prompt input.
type ValidatorFunc = promptui.ValidateFunc

// A Prompter prompts the user for input.
type Prompter struct {
	validators map[string]ValidatorFunc
}

// New creates a new Prompter.
func New() *Prompter {
	return new(Prompter)
}

// NewWithValidators creates a new Prompter with the given validator functions.
func NewWithValidators(validators map[string]ValidatorFunc) *Prompter {
	p := New()
	p.validators = validators
	return p
}

// Run executes the Prompter for the given struct.
func (p *Prompter) Run(v interface{}) error {
	elem := reflect.ValueOf(v).Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("can only process struct types but provided type is %q", elem.Kind())
	}

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)

		tags, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return err
		}

		tag, err := tags.Get("prompt")
		if err != nil {
			// No tag present, initialize an empty one.
			tag = new(structtag.Tag)
		} else if tag.Name == TagValueSkip {
			continue
		}

		fieldVal := elem.Field(i)

		prompt := promptui.Prompt{
			Label: field.Name,
		}

		ft := field.Type
		switch ft.Kind() {
		case reflect.String:
			prompt.Default = fieldVal.String()
			prompt.AllowEdit = true
		case reflect.Bool:
			prompt.IsConfirm = true
		default:
			return fmt.Errorf("can only process string and bool fields but field %q is type %q", field.Name, ft.Kind())
		}

		if tag.Name == TagValueRequired {
			prompt.Validate = validateRequired
		}

		for _, o := range tag.Options {
			var val string
			optVal := strings.Split(o, "=")
			opt := optVal[0]
			if len(optVal) > 1 {
				val = optVal[1]
			}

			switch opt {
			case TagValueValidate:
				valFunc, ok := p.validators[val]
				if !ok {
					return fmt.Errorf("unknown validation function %q", val)
				}
				prompt.Validate = valdatorChain(tag.Name != TagValueRequired, prompt.Validate, valFunc)
			case TagValueName:
				prompt.Label = val
			case TagValueMask:
				prompt.Mask = rune(val[0])
			}
		}

		res, err := prompt.Run()
		if err == promptui.ErrAbort {
			fieldVal.SetBool(false)
			continue
		} else if err != nil {
			return err
		} else if ft.Kind() == reflect.Bool {
			fieldVal.SetBool(true)
			continue
		}
		fieldVal.SetString(res)
	}
	return nil
}
