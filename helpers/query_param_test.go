package helpers

import (
	"net/http"
	"net/url"
	"testing"
)

func TestNewQueryParams(t *testing.T) {
	// Criar uma requisição HTTP para teste
	req, err := http.NewRequest("GET", "http://example.com?existing=value", nil)
	if err != nil {
		t.Fatalf("Erro ao criar requisição: %v", err)
	}

	// Testar criação de novo QueryParams
	qp := NewQueryParams(req)
	if qp == nil {
		t.Fatal("NewQueryParams retornou nil")
	}

	// Verificar se o tipo retornado implementa a interface QueryParams
	_, ok := qp.(QueryParams)
	if !ok {
		t.Fatal("NewQueryParams não retornou um tipo que implementa QueryParams")
	}

	// Verificar se mantém parâmetros existentes da URL
	encoded := qp.Encode()
	if encoded != "existing=value" {
		t.Errorf("Esperado 'existing=value', obtido '%s'", encoded)
	}
}

func TestQueryParam_Add(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	qp := NewQueryParams(req)

	t.Run("Adicionar valor string válido", func(t *testing.T) {
		value := "test_value"
		qp.Add("test_param", &value)

		encoded := qp.Encode()
		if encoded != "test_param=test_value" {
			t.Errorf("Esperado 'test_param=test_value', obtido '%s'", encoded)
		}
	})

	t.Run("Adicionar valor nil", func(t *testing.T) {
		// Criar nova instância para teste isolado
		req2, _ := http.NewRequest("GET", "http://example.com", nil)
		qp2 := NewQueryParams(req2)

		qp2.Add("nil_param", nil)

		encoded := qp2.Encode()
		if encoded != "" {
			t.Errorf("Esperado string vazia para valor nil, obtido '%s'", encoded)
		}
	})

	t.Run("Adicionar múltiplos parâmetros", func(t *testing.T) {
		req3, _ := http.NewRequest("GET", "http://example.com", nil)
		qp3 := NewQueryParams(req3)

		param1 := "value1"
		param2 := "value2"
		qp3.Add("param1", &param1)
		qp3.Add("param2", &param2)

		encoded := qp3.Encode()
		// A ordem pode variar, então verificar se ambos estão presentes
		expectedValues := map[string]string{
			"param1": "value1",
			"param2": "value2",
		}

		parsedValues, err := url.ParseQuery(encoded)
		if err != nil {
			t.Fatalf("Erro ao fazer parse do query string: %v", err)
		}

		for key, expectedValue := range expectedValues {
			if parsedValues.Get(key) != expectedValue {
				t.Errorf("Para %s: esperado '%s', obtido '%s'", key, expectedValue, parsedValues.Get(key))
			}
		}
	})

	t.Run("Adicionar string vazia", func(t *testing.T) {
		req4, _ := http.NewRequest("GET", "http://example.com", nil)
		qp4 := NewQueryParams(req4)

		emptyValue := ""
		qp4.Add("empty_param", &emptyValue)

		encoded := qp4.Encode()
		if encoded != "empty_param=" {
			t.Errorf("Esperado 'empty_param=', obtido '%s'", encoded)
		}
	})
}

func TestQueryParam_AddReflect(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	qp := NewQueryParams(req)

	t.Run("Adicionar string usando reflection", func(t *testing.T) {
		stringValue := "test_string"
		qp.AddReflect("string_param", stringValue)

		encoded := qp.Encode()
		if encoded != "string_param=test_string" {
			t.Errorf("Esperado 'string_param=test_string', obtido '%s'", encoded)
		}
	})

	t.Run("Adicionar int usando reflection", func(t *testing.T) {
		req2, _ := http.NewRequest("GET", "http://example.com", nil)
		qp2 := NewQueryParams(req2)

		intValue := 42
		qp2.AddReflect("int_param", intValue)

		encoded := qp2.Encode()
		if encoded != "int_param=42" {
			t.Errorf("Esperado 'int_param=42', obtido '%s'", encoded)
		}
	})

	t.Run("Adicionar valor nil usando reflection", func(t *testing.T) {
		req3, _ := http.NewRequest("GET", "http://example.com", nil)
		qp3 := NewQueryParams(req3)

		qp3.AddReflect("nil_param", nil)

		encoded := qp3.Encode()
		if encoded != "" {
			t.Errorf("Esperado string vazia para valor nil, obtido '%s'", encoded)
		}
	})

	t.Run("Adicionar tipo não suportado", func(t *testing.T) {
		req4, _ := http.NewRequest("GET", "http://example.com", nil)
		qp4 := NewQueryParams(req4)

		// Testar com um tipo não suportado (float64)
		floatValue := 3.14
		qp4.AddReflect("float_param", floatValue)

		encoded := qp4.Encode()
		if encoded != "" {
			t.Errorf("Esperado string vazia para tipo não suportado, obtido '%s'", encoded)
		}
	})

	t.Run("Adicionar tipos não suportados diversos", func(t *testing.T) {
		req5, _ := http.NewRequest("GET", "http://example.com", nil)
		qp5 := NewQueryParams(req5)

		// Testar com diferentes tipos não suportados
		qp5.AddReflect("bool_param", true)
		qp5.AddReflect("slice_param", []string{"a", "b"})
		qp5.AddReflect("map_param", map[string]string{"key": "value"})

		encoded := qp5.Encode()
		if encoded != "" {
			t.Errorf("Esperado string vazia para tipos não suportados, obtido '%s'", encoded)
		}
	})

	t.Run("Misturar string e int usando reflection", func(t *testing.T) {
		req6, _ := http.NewRequest("GET", "http://example.com", nil)
		qp6 := NewQueryParams(req6)

		stringValue := "hello"
		intValue := 123

		qp6.AddReflect("str", stringValue)
		qp6.AddReflect("num", intValue)

		encoded := qp6.Encode()
		parsedValues, err := url.ParseQuery(encoded)
		if err != nil {
			t.Fatalf("Erro ao fazer parse do query string: %v", err)
		}

		if parsedValues.Get("str") != "hello" {
			t.Errorf("Esperado 'hello' para str, obtido '%s'", parsedValues.Get("str"))
		}

		if parsedValues.Get("num") != "123" {
			t.Errorf("Esperado '123' para num, obtido '%s'", parsedValues.Get("num"))
		}
	})

	t.Run("Adicionar int zero usando reflection", func(t *testing.T) {
		req7, _ := http.NewRequest("GET", "http://example.com", nil)
		qp7 := NewQueryParams(req7)

		zeroValue := 0
		qp7.AddReflect("zero_param", zeroValue)

		encoded := qp7.Encode()
		if encoded != "zero_param=0" {
			t.Errorf("Esperado 'zero_param=0', obtido '%s'", encoded)
		}
	})

	t.Run("Adicionar string vazia usando reflection", func(t *testing.T) {
		req8, _ := http.NewRequest("GET", "http://example.com", nil)
		qp8 := NewQueryParams(req8)

		emptyString := ""
		qp8.AddReflect("empty_string", emptyString)

		encoded := qp8.Encode()
		if encoded != "empty_string=" {
			t.Errorf("Esperado 'empty_string=', obtido '%s'", encoded)
		}
	})

	t.Run("Adicionar int negativo usando reflection", func(t *testing.T) {
		req9, _ := http.NewRequest("GET", "http://example.com", nil)
		qp9 := NewQueryParams(req9)

		negativeValue := -42
		qp9.AddReflect("negative_param", negativeValue)

		encoded := qp9.Encode()
		if encoded != "negative_param=-42" {
			t.Errorf("Esperado 'negative_param=-42', obtido '%s'", encoded)
		}
	})
}

func TestQueryParam_Encode(t *testing.T) {
	t.Run("Encode com parâmetros vazios", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		qp := NewQueryParams(req)

		encoded := qp.Encode()
		if encoded != "" {
			t.Errorf("Esperado string vazia, obtido '%s'", encoded)
		}
	})

	t.Run("Encode com parâmetros existentes na URL", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com?pre=existing&another=param", nil)
		qp := NewQueryParams(req)

		encoded := qp.Encode()
		parsedValues, err := url.ParseQuery(encoded)
		if err != nil {
			t.Fatalf("Erro ao fazer parse do query string: %v", err)
		}

		if parsedValues.Get("pre") != "existing" {
			t.Errorf("Esperado 'existing' para pre, obtido '%s'", parsedValues.Get("pre"))
		}

		if parsedValues.Get("another") != "param" {
			t.Errorf("Esperado 'param' para another, obtido '%s'", parsedValues.Get("another"))
		}
	})

	t.Run("Encode com caracteres especiais", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		qp := NewQueryParams(req)

		specialValue := "value with spaces & special chars"
		qp.Add("special", &specialValue)

		encoded := qp.Encode()
		// Verificar se foi codificado corretamente
		if !containsEncodedValue(encoded, "value+with+spaces") || !containsEncodedValue(encoded, "%26") {
			t.Errorf("Codificação especial não aplicada corretamente: %s", encoded)
		}
	})
}

func TestQueryParam_IntegrationTest(t *testing.T) {
	t.Run("Teste de integração completo", func(t *testing.T) {
		// Começar com URL que já tem parâmetros
		req, _ := http.NewRequest("GET", "http://example.com?existing=value", nil)
		qp := NewQueryParams(req)

		// Adicionar usando Add
		newValue := "added_value"
		qp.Add("added", &newValue)

		// Adicionar usando AddReflect com string
		qp.AddReflect("reflected_str", "reflected_string")

		// Adicionar usando AddReflect com int
		qp.AddReflect("reflected_int", 999)

		// Tentar adicionar nil (não deve adicionar nada)
		qp.Add("nil_test", nil)
		qp.AddReflect("nil_reflect", nil)

		encoded := qp.Encode()
		parsedValues, err := url.ParseQuery(encoded)
		if err != nil {
			t.Fatalf("Erro ao fazer parse do query string: %v", err)
		}

		// Verificar todos os valores esperados
		expectedValues := map[string]string{
			"existing":      "value",
			"added":         "added_value",
			"reflected_str": "reflected_string",
			"reflected_int": "999",
		}

		for key, expectedValue := range expectedValues {
			if parsedValues.Get(key) != expectedValue {
				t.Errorf("Para %s: esperado '%s', obtido '%s'", key, expectedValue, parsedValues.Get(key))
			}
		}

		// Verificar que valores nil não foram adicionados
		if parsedValues.Has("nil_test") {
			t.Error("nil_test não deveria estar presente")
		}

		if parsedValues.Has("nil_reflect") {
			t.Error("nil_reflect não deveria estar presente")
		}
	})
}

func TestQueryParam_EdgeCases(t *testing.T) {
	t.Run("Sobrescrever parâmetro existente", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com?param=original", nil)
		qp := NewQueryParams(req)

		newValue := "overwritten"
		qp.Add("param", &newValue)

		encoded := qp.Encode()
		if encoded != "param=overwritten" {
			t.Errorf("Esperado 'param=overwritten', obtido '%s'", encoded)
		}
	})

	t.Run("Misturar Add e AddReflect no mesmo parâmetro", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com", nil)
		qp := NewQueryParams(req)

		// Primeiro adicionar com Add
		firstValue := "first"
		qp.Add("param", &firstValue)

		// Depois sobrescrever com AddReflect
		qp.AddReflect("param", "second")

		encoded := qp.Encode()
		if encoded != "param=second" {
			t.Errorf("Esperado 'param=second', obtido '%s'", encoded)
		}
	})

	t.Run("URL com parâmetros malformados", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com?malformed=%ZZ", nil)
		qp := NewQueryParams(req)

		// Deve funcionar mesmo com parâmetros malformados na URL original
		newValue := "valid"
		qp.Add("new_param", &newValue)

		encoded := qp.Encode()
		parsedValues, err := url.ParseQuery(encoded)
		if err != nil {
			t.Fatalf("Erro ao fazer parse do query string: %v", err)
		}

		if parsedValues.Get("new_param") != "valid" {
			t.Errorf("Esperado 'valid' para new_param, obtido '%s'", parsedValues.Get("new_param"))
		}
	})
}

// Função auxiliar para verificar se uma string codificada contém um valor específico
func containsEncodedValue(encoded, value string) bool {
	decoded, err := url.QueryUnescape(encoded)
	if err != nil {
		return false
	}
	return containsSubstring(decoded, value) || containsSubstring(encoded, value)
}

// Função auxiliar para verificar substring
func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
