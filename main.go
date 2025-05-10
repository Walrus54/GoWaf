package main

import (
	"log"
	"waf/internal/proxy"
	"waf/internal/waf"
	"waf/internal/web"
)

func main() {
	wafInstance, err := waf.NewWAF()
	if err != nil {
		log.Fatal("WAF init failed:", err)
	}

	// Запуск прокси-сервера
	go func() {
		log.Printf("WAF proxy started on %s", wafInstance.Config.WafPort)
		log.Fatal(proxy.StartServer(wafInstance))
	}()

	// Запуск админ-интерфейса
	log.Printf("Admin interface started on %s", wafInstance.Config.AdminPort)
	log.Fatal(web.StartAdminInterface(wafInstance))
}
