package k8s

import (
	"os"
	"path/filepath"
	"testing"
)

func TestKubeConfig(t *testing.T) {
	dir := t.TempDir()
	if _, err := NewClient(filepath.Join(dir, "missing")); err == nil {
		t.Fatal("missing config accepted")
	}
	path := filepath.Join(dir, "config")
	data := `apiVersion: v1
kind: Config
clusters:
- name: local
  cluster:
    server: http://127.0.0.1:1
contexts:
- name: local
  context:
    cluster: local
current-context: local
`
	if err := os.WriteFile(path, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	client, err := NewClient(path)
	if err != nil || client.Clientset == nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	t.Setenv("KUBERNETES_SERVICE_HOST", "")
	t.Setenv("KUBERNETES_SERVICE_PORT", "")
	if _, err := NewClient(""); err == nil {
		t.Fatal("missing in-cluster config accepted")
	}
}
