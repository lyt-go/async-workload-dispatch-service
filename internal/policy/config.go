package policy

type Config struct{ Rules map[string]bool }

func Load(names []string) Config {
	// 即使 names 为空也要初始化规则表，否则后续写入会 panic。
	rules := make(map[string]bool)
	for _, name := range names {
		rules[name] = true
	}
	return Config{Rules: rules}
}
