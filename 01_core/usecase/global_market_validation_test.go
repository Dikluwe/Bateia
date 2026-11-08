/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/01_core/usecase/global_market_validation_test.md
 * @layer L1
 * @updated 2026-10-28
 */
package usecase

import "testing"

// TestStep0033ISO3166AndBCP47ValidatorsAcceptGlobalMarkets confirma que os validadores de L1 e L3
// aceitam códigos de país ISO 3166-1 alpha-2 e locales BCP-47 para além dos 10 piloto originais.
func TestStep0033ISO3166AndBCP47ValidatorsAcceptGlobalMarkets(t *testing.T) {
	// Mercados válidos: originais + expansão global
	validMarkets := []string{
		// Piloto original
		"BR", "DE", "ES", "FR", "IN", "JP", "KR", "RU", "SA", "US",
		// Expansão global — Passo 0033
		"IT", "GR", "MN", "KZ", "NG", "ZA", "EG", "PL", "TR", "MX", "AR",
		"AU", "CA", "CN", "TH", "VN", "SE", "NL", "BE", "PT", "CH", "AT",
	}
	// Locales válidos: piloto original + expansão global
	validLocales := []string{
		// Piloto original
		"en-US", "pt-BR", "de-DE", "fr-FR", "es-ES", "ja-JP", "ko-KR", "ru-RU", "hi-IN", "ar-SA",
		// Expansão global — Passo 0033
		"it-IT", "el-GR", "mn-MN", "kk-KZ", "tr-TR", "pl-PL", "sv-SE", "nl-NL",
		"zh-CN", "th-TH", "vi-VN", "ar-EG", "es-MX", "pt-PT", "en-GB", "en-ZA",
	}
	for _, market := range validMarkets {
		if !validISO3166Alpha2(market) {
			t.Errorf("mercado ISO válido rejeitado por L1: %q", market)
		}
	}
	for _, locale := range validLocales {
		if !validBCP47Locale(locale) {
			t.Errorf("locale BCP-47 válido rejeitado por L1: %q", locale)
		}
	}
}

// TestStep0033InvalidInputsAreRejected confirma que entradas malformadas continuam sendo rejeitadas.
func TestStep0033InvalidInputsAreRejected(t *testing.T) {
	invalidMarkets := []string{"xx", "br", "us", "123", "BR1", "USA", "12", "", "A", "ABC", "b2"}
	for _, market := range invalidMarkets {
		if validISO3166Alpha2(market) {
			t.Errorf("mercado inválido aceito por L1: %q", market)
		}
	}
	invalidLocales := []string{
		"en_US", "EN-US", "brBR", "BR", "123", "", "pt-br", "PT-BR",
		"en-us", "ENGB", "en-123", "abc-DE", "a-DE",
	}
	for _, locale := range invalidLocales {
		if validBCP47Locale(locale) {
			t.Errorf("locale inválido aceito por L1: %q", locale)
		}
	}
}
