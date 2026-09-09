package schemas

import (
	"encoding/json"
	"os"
	"testing"
)

func TestSchemaAuthAndProviderGuards(t *testing.T) {
	content, err := os.ReadFile("csl-config.schema.json")
	if err != nil {
		t.Fatal(err)
	}

	var schema map[string]any
	if err := json.Unmarshal(content, &schema); err != nil {
		t.Fatal(err)
	}

	llm := schema["properties"].(map[string]any)["llm"].(map[string]any)
	llmProperties := llm["properties"].(map[string]any)
	authProperties := llmProperties["auth"].(map[string]any)["properties"].(map[string]any)
	awsProperties := authProperties["aws"].(map[string]any)["properties"].(map[string]any)
	if _, ok := authProperties["google"]; ok {
		t.Fatal("auth.google should not be accepted")
	}
	if _, ok := authProperties["method"]; !ok {
		t.Fatal("auth.method missing")
	}
	for _, property := range []string{"aws", "gcloud", "azure"} {
		if _, ok := authProperties[property]; !ok {
			t.Fatalf("auth.%s missing", property)
		}
	}
	for _, property := range []string{"region", "profile"} {
		if _, ok := awsProperties[property]; !ok {
			t.Fatalf("auth.aws.%s missing", property)
		}
	}
	for _, property := range []string{"base_url", "host"} {
		if llmProperties[property].(map[string]any)["pattern"] != "^https?://.+" {
			t.Fatalf("%s http pattern missing", property)
		}
	}

	for _, rule := range llm["allOf"].([]any) {
		ruleObject := rule.(map[string]any)
		required, ok := ruleObject["if"].(map[string]any)["required"].([]any)
		if !ok || len(required) != 1 || required[0] != "claude_code_session" {
			continue
		}
		provider := ruleObject["then"].(map[string]any)["properties"].(map[string]any)["provider"].(map[string]any)
		if provider["const"] != "claude-code" {
			t.Fatalf("claude_code_session provider const = %v", provider["const"])
		}
		return
	}

	t.Fatal("claude_code_session provider condition missing")
}

func TestSchemaPRAndHintSections(t *testing.T) {
	content, err := os.ReadFile("csl-config.schema.json")
	if err != nil {
		t.Fatal(err)
	}

	var schema map[string]any
	if err := json.Unmarshal(content, &schema); err != nil {
		t.Fatal(err)
	}

	properties := schema["properties"].(map[string]any)

	pr, ok := properties["pr"].(map[string]any)
	if !ok {
		t.Fatal("pr section missing")
	}
	prProperties := pr["properties"].(map[string]any)
	assertEnum(t, prProperties, "merge", []string{"squash", "rebase", "merge"})
	if _, ok := prProperties["title_max_length"]; !ok {
		t.Fatal("pr.title_max_length missing")
	}

	hint, ok := properties["hint"].(map[string]any)
	if !ok {
		t.Fatal("hint section missing")
	}
	assertEnum(t, hint["properties"].(map[string]any), "program", []string{"url", "gh", "none"})
}

func TestSchemaDiffSection(t *testing.T) {
	content, err := os.ReadFile("csl-config.schema.json")
	if err != nil {
		t.Fatal(err)
	}

	var schema map[string]any
	if err := json.Unmarshal(content, &schema); err != nil {
		t.Fatal(err)
	}

	properties := schema["properties"].(map[string]any)

	diff, ok := properties["diff"].(map[string]any)
	if !ok {
		t.Fatal("diff section missing")
	}
	diffProperties := diff["properties"].(map[string]any)

	for _, key := range []string{"exclude", "exclude_regex"} {
		property, ok := diffProperties[key].(map[string]any)
		if !ok {
			t.Fatalf("diff.%s missing", key)
		}
		if property["type"] != "array" {
			t.Fatalf("diff.%s type = %v, want array", key, property["type"])
		}
		items := property["items"].(map[string]any)
		if items["type"] != "string" {
			t.Fatalf("diff.%s items type = %v, want string", key, items["type"])
		}
	}
}

func assertEnum(t *testing.T, properties map[string]any, key string, want []string) {
	t.Helper()
	property, ok := properties[key].(map[string]any)
	if !ok {
		t.Fatalf("%s missing", key)
	}
	values := property["enum"].([]any)
	if len(values) != len(want) {
		t.Fatalf("%s enum = %v, want %v", key, values, want)
	}
	for i, value := range want {
		if values[i] != value {
			t.Fatalf("%s enum[%d] = %v, want %s", key, i, values[i], value)
		}
	}
}
