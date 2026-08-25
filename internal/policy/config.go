package policy

type Config struct{ Rules map[string]bool }

func Load(names []string) Config {
	var rules map[string]bool
	for _, name := range names {
		rules[name] = true
	}
	return Config{Rules: rules}
}
