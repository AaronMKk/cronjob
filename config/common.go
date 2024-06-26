package config

type configValidate interface {
	Validate() error
}

type configSetDefault interface {
	SetDefault()
}

type configItems interface {
	ConfigItems() []interface{}
}

// SetDefault is set item config with default value
func SetDefault(cfg interface{}) {
	if f, ok := cfg.(configSetDefault); ok {
		f.SetDefault()
	}

	if f, ok := cfg.(configItems); ok {
		items := f.ConfigItems()

		for i := range items {
			SetDefault(items[i])
		}
	}
}

// Validate for validate item cfg
func Validate(cfg interface{}) error {
	if f, ok := cfg.(configValidate); ok {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	if f, ok := cfg.(configItems); ok {
		items := f.ConfigItems()

		for i := range items {
			if err := Validate(items[i]); err != nil {
				return err
			}
		}
	}

	return nil
}
