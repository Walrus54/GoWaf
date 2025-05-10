package waf

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"waf/internal/logger"
	"waf/internal/rules"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TargetURL string `yaml:"target_url"`
	WafPort   string `yaml:"waf_port"`
	AdminPort string `yaml:"admin_port"`
	LogFile   string `yaml:"log_file"`
	Mode      string `yaml:"mode"`
	RulesFile string `yaml:"rules_file"`
}

type WAF struct {
	Config Config
	Rules  map[string]rules.Rule
	Logger *logger.Logger
	mu     sync.RWMutex
	nextID int // Добавляем счетчик для ID
}

func loadConfig(path string) (Config, error) {
	var config Config
	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

func NewWAF() (*WAF, error) {
	config, err := loadConfig("config/config.yaml")
	if err != nil {
		return nil, err
	}

	logger, err := logger.NewLogger(config.LogFile)
	if err != nil {
		return nil, err
	}

	// Получаем 3 значения из LoadRulesMap
	rulesMap, nextID, err := rules.LoadRulesMap(config.RulesFile)
	if err != nil {
		return nil, err
	}

	return &WAF{
		Config: config,
		Rules:  rulesMap,
		Logger: logger,
		nextID: nextID, // Инициализируем счетчик
	}, nil
}

// Методы для веб-интерфейса
func (w *WAF) SetMode(mode string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Config.Mode = mode
}

func (w *WAF) AddRule(name, pattern string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Генерируем числовой ID
	id := strconv.Itoa(w.nextID)
	w.nextID++

	rule := rules.Rule{
		ID:         id,
		Name:       name,
		RawPattern: pattern, // Используем RawPattern вместо Pattern
	}

	// Компилируем регулярное выражение
	if err := rule.Compile(); err != nil {
		return err
	}

	w.Rules[rule.ID] = rule
	return w.saveRules()
}

func (w *WAF) DeleteRule(id string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.Rules, id)
	return w.saveRules()
}

func (w *WAF) GetRules() []rules.Rule {
	w.mu.RLock()
	defer w.mu.RUnlock()

	result := make([]rules.Rule, 0, len(w.Rules))
	for _, rule := range w.Rules {
		result = append(result, rule)
	}
	return result
}

func (w *WAF) saveRules() error {
	data, err := yaml.Marshal(w.Rules)
	if err != nil {
		return err
	}
	return os.WriteFile(w.Config.RulesFile, data, 0644)
}

func (w *WAF) CheckRequest(r *http.Request) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	checkParts := []string{
		r.URL.String(),
		r.Header.Get("User-Agent"),
		r.Header.Get("Referer"),
	}

	// Добавляем тело запроса для POST
	if r.Method == http.MethodPost {
		body, _ := io.ReadAll(r.Body)
		checkParts = append(checkParts, string(body))
		r.Body = io.NopCloser(bytes.NewBuffer(body)) // Восстанавливаем тело
	}

	for _, part := range checkParts {
		for _, rule := range w.Rules {
			if rule.Pattern.MatchString(part) {
				return false
			}
		}
	}
	return true
}
