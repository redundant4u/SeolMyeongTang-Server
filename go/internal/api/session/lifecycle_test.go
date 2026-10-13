package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"seolmyeong-tang-server/internal/config"
	"seolmyeong-tang-server/internal/pkg/k8s"
	"seolmyeong-tang-server/internal/pkg/validator"
)

func testKube(t *testing.T, serve http.HandlerFunc) *Kube {
	t.Helper()
	server := httptest.NewServer(serve)
	t.Cleanup(server.Close)
	client, err := kubernetes.NewForConfig(&rest.Config{Host: server.URL, ContentConfig: rest.ContentConfig{ContentType: "application/json", AcceptContentTypes: "application/json"}})
	if err != nil {
		t.Fatal(err)
	}
	old := config.Env
	t.Cleanup(func() { config.Env = old })
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("KUBE_SESSION_NAMESPACE=sessions\nCF_TUNNEL_ID=test-tunnel\nAWS_ACCESS_KEY=test\nAWS_SECRET_KEY=test\nDYNAMODB_TABLE=posts\n"), 0600); err != nil {
		t.Fatal(err)
	}
	config.InitEnv()
	return NewKube(&k8s.Client{Clientset: client}, "sessions")
}

func replyJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func apiFailure(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(metav1.Status{TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"}, Status: "Failure", Reason: metav1.StatusReasonForbidden, Code: 403, Message: "denied"})
}

func samplePod() corev1.Pod {
	p := buildSessionPodSpec(createPod{name: "demo", image: "ubuntu", clientId: "client01", sessionId: "session1", description: "work"}, "sessions", time.Now())
	p.Status.Phase = corev1.PodRunning
	return *p
}

func TestPodSpecificationAndResponseContract(t *testing.T) {
	p := samplePod()
	if p.Name != "session1" || p.Labels["client-id"] != "client01" || p.Annotations["description"] != "work" {
		t.Fatalf("incorrect metadata: %+v", p.ObjectMeta)
	}
	if *p.Spec.AutomountServiceAccountToken || *p.Spec.RuntimeClassName != runtimeClassName || len(p.Spec.Containers) != 2 {
		t.Fatal("isolation or sidecar missing")
	}
	vnc, tunnel := p.Spec.Containers[0], p.Spec.Containers[1]
	if vnc.Image != "vnc:ubuntu" || vnc.Ports[0].ContainerPort != 5901 || vnc.ImagePullPolicy != corev1.PullNever {
		t.Fatal("incorrect vnc container")
	}
	if !*tunnel.SecurityContext.RunAsNonRoot || *tunnel.SecurityContext.AllowPrivilegeEscalation {
		t.Fatal("tunnel security settings changed")
	}
	created, _ := time.Parse(time.RFC3339, p.Annotations["created-at"])
	expired, _ := time.Parse(time.RFC3339, p.Annotations["expired-at"])
	if expired.Sub(created) != 10*time.Minute {
		t.Fatal("incorrect session lifetime")
	}
	list, err := toGetSessionsResponse([]corev1.Pod{p})
	if err != nil || len(list) != 1 || list[0].SessionId != "session1" || list[0].TTL < 590 {
		t.Fatalf("response: %+v %v", list, err)
	}
	single, err := toCreateSessionResponse(&p, "session1")
	if err != nil || single.Image != "vnc:ubuntu" || single.Description != "work" {
		t.Fatalf("response: %+v %v", single, err)
	}
	delete(p.Annotations, "description")
	single, err = toCreateSessionResponse(&p, "session1")
	list, listErr := toGetSessionsResponse([]corev1.Pod{p})
	if err != nil || listErr != nil || single.Description != "" || list[0].Description != "" {
		t.Fatal("description must be optional")
	}
	for _, value := range []string{"", "invalid", time.Now().Add(-time.Minute).Format(time.RFC3339)} {
		p.Annotations["expired-at"] = value
		if extractTTL(&p) != 0 {
			t.Fatal("expired/invalid TTL must be zero")
		}
	}
	delete(p.Annotations, "expired-at")
	if extractTTL(&p) != 0 {
		t.Fatal("missing TTL must be zero")
	}
	delete(p.Labels, "name")
	if _, err := toGetSessionsResponse([]corev1.Pod{p}); err == nil {
		t.Fatal("missing name accepted")
	}
	if _, err := toCreateSessionResponse(&p, "session1"); err == nil {
		t.Fatal("missing name accepted")
	}
	if list, err := toGetSessionsResponse(nil); err != nil || len(list) != 0 {
		t.Fatal("empty list rejected")
	}
}

func TestSessionHTTP(t *testing.T) {
	cases := []struct {
		name, method, body, fail string
		pods                     int
		code                     int
	}{
		{"list", "GET", "", "", 1, 200},
		{"empty list", "GET", "", "", 0, 200},
		{"list failure", "GET", "", "list", 0, 400},
		{"list mapping failure", "GET", "", "map", 1, 400},
		{"create", "POST", `{"name":"demo","image":"ubuntu","description":"work"}`, "", 0, 201},
		{"missing name", "POST", `{"image":"ubuntu"}`, "", 0, 400},
		{"long name", "POST", `{"name":"` + strings.Repeat("a", 21) + `","image":"ubuntu"}`, "", 0, 400},
		{"long description", "POST", `{"name":"demo","image":"ubuntu","description":"` + strings.Repeat("a", 101) + `"}`, "", 0, 400},
		{"valid bounds", "POST", `{"name":"` + strings.Repeat("a", 20) + `","image":"ubuntu","description":"` + strings.Repeat("a", 100) + `"}`, "", 0, 201},
		{"capacity", "POST", `{"name":"demo","image":"ubuntu"}`, "", 4, 400},
		{"create list failure", "POST", `{"name":"demo","image":"ubuntu"}`, "list", 0, 400},
		{"configmap failure", "POST", `{"name":"demo","image":"ubuntu"}`, "configmaps", 0, 400},
		{"pod creation failure", "POST", `{"name":"demo","image":"ubuntu"}`, "create", 0, 400},
		{"create mapping failure", "POST", `{"name":"demo","image":"ubuntu"}`, "map", 0, 400},
		{"delete", "DELETE", `{"sessionId":"session1"}`, "", 1, 204},
		{"delete absent pod", "DELETE", `{"sessionId":"session1"}`, "", 0, 204},
		{"delete configmap failure", "DELETE", `{"sessionId":"session1"}`, "configmaps", 1, 400},
		{"delete list failure", "DELETE", `{"sessionId":"session1"}`, "list", 1, 400},
		{"delete pod failure", "DELETE", `{"sessionId":"session1"}`, "delete", 1, 400},
		{"delete invalid body", "DELETE", "not json", "", 0, 400},
		{"delete invalid id", "DELETE", `{"sessionId":"INVALID"}`, "", 0, 400},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			writes := 0
			kube := testKube(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					writes++
				}
				isCM := strings.Contains(r.URL.Path, "configmaps")
				if (tc.fail == "configmaps" && isCM) || (tc.fail == "list" && r.Method == "GET") || (tc.fail == "create" && !isCM && r.Method == "POST") || (tc.fail == "delete" && !isCM && r.Method == "DELETE") {
					apiFailure(w)
					return
				}
				switch r.Method {
				case "GET":
					if !strings.Contains(r.URL.Query().Get("labelSelector"), "client-id=client01") {
						t.Error("client filter missing")
					}
					pods := []corev1.Pod{}
					for i := 0; i < tc.pods; i++ {
						p := samplePod()
						if tc.fail == "map" {
							delete(p.Labels, "name")
						}
						pods = append(pods, p)
					}
					replyJSON(w, corev1.PodList{Items: pods})
				case "POST":
					if isCM {
						var cm corev1.ConfigMap
						_ = json.NewDecoder(r.Body).Decode(&cm)
						if !strings.Contains(cm.Data["config.yml"], "http://localhost:8080") {
							t.Error("invalid tunnel configuration")
						}
						replyJSON(w, cm)
					} else {
						var p corev1.Pod
						_ = json.NewDecoder(r.Body).Decode(&p)
						if tc.fail == "map" {
							delete(p.Labels, "name")
						}
						replyJSON(w, p)
					}
				default:
					replyJSON(w, metav1.Status{Status: "Success"})
				}
			})
			e := echo.New()
			e.Validator = validator.New()
			Init(e, kube)
			r := httptest.NewRequest(tc.method, "/session", strings.NewReader(tc.body))
			r.Header.Set("X-Client-Id", "client01")
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			e.ServeHTTP(w, r)
			if w.Code != tc.code {
				t.Fatalf("want %d got %d: %s", tc.code, w.Code, w.Body.String())
			}
			if tc.name == "create" {
				var res createPodResponse
				_ = json.Unmarshal(w.Body.Bytes(), &res)
				if len(res.SessionId) != 8 || res.Name != "demo" || writes != 2 {
					t.Fatalf("creation contract: %+v writes=%d", res, writes)
				}
			}
			if tc.code == 400 && tc.fail == "" && writes != 0 {
				t.Fatal("rejected request performed external writes")
			}
		})
	}
	e := echo.New()
	e.GET("/session", NewHandler(&Kube{}).getSessions)
	w := httptest.NewRecorder()
	e.ServeHTTP(w, httptest.NewRequest("GET", "/session", nil))
	if w.Code != 400 {
		t.Fatal("missing client accepted")
	}
}

func TestPodFiltering(t *testing.T) {
	kube := testKube(t, func(w http.ResponseWriter, r *http.Request) {
		running := samplePod()
		pending := samplePod()
		pending.Status.Phase = corev1.PodPending
		failed := samplePod()
		failed.Status.Phase = corev1.PodFailed
		deleting := samplePod()
		now := metav1.Now()
		deleting.DeletionTimestamp = &now
		replyJSON(w, corev1.PodList{Items: []corev1.Pod{running, pending, failed, deleting}})
	})
	pods, err := kube.getSessions(context.Background(), "client01")
	if err != nil || len(pods) != 2 {
		t.Fatalf("filter failed: %d %v", len(pods), err)
	}
}

func TestGC(t *testing.T) {
	for _, mode := range []string{"success", "list failure", "delete failure", "pagination"} {
		t.Run(mode, func(t *testing.T) {
			deletes := []string{}
			pages := 0
			kube := testKube(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "DELETE" {
					deletes = append(deletes, r.URL.Path)
					if mode == "delete failure" {
						apiFailure(w)
					} else {
						replyJSON(w, metav1.Status{Status: "Success"})
					}
					return
				}
				if mode == "list failure" {
					apiFailure(w)
					return
				}
				pages++
				if pages == 2 {
					if r.URL.Query().Get("continue") != "next" {
						t.Error("pagination cursor lost")
					}
					replyJSON(w, corev1.PodList{})
					return
				}
				expired := samplePod()
				expired.Annotations["expired-at"] = time.Now().Add(-time.Minute).Format(time.RFC3339)
				valid := samplePod()
				missing := samplePod()
				delete(missing.Annotations, "expired-at")
				invalid := samplePod()
				invalid.Annotations["expired-at"] = "bad"
				deleting := samplePod()
				now := metav1.Now()
				deleting.DeletionTimestamp = &now
				list := corev1.PodList{Items: []corev1.Pod{expired, valid, missing, invalid, deleting}}
				if mode == "pagination" {
					list.Continue = "next"
				}
				replyJSON(w, list)
			})
			kube.Gc.cleanup(context.Background())
			if mode == "list failure" {
				if len(deletes) != 0 {
					t.Fatal("deleted after list failure")
				}
			} else if len(deletes) != 1 || !strings.HasSuffix(deletes[0], "/session1") {
				t.Fatalf("only expired pod should be deleted: %v", deletes)
			}
		})
	}
}

func TestGCCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	g := newGC(nil, "sessions", 60)
	go func() { g.Run(ctx); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("GC ignores cancellation")
	}
}

func TestKubeErrors(t *testing.T) {
	denied := kerrors.NewForbidden(schema.GroupResource{Resource: "pods"}, "demo", errors.New("denied"))
	if kubeErrorCode(nil) != "" || kubeErrorCode(errors.New("network")) != "api_error" || kubeErrorCode(denied) != "Forbidden" {
		t.Fatal("error classification changed")
	}
	if kubeErrorCode(&kerrors.StatusError{}) != "api_error" {
		t.Fatal("empty status reason not normalized")
	}
	if unwrapKubeStatus(nil) != nil {
		t.Fatal("nil error changed")
	}
	original := errors.New("network")
	if unwrapKubeStatus(original) != original {
		t.Fatal("non-status error changed")
	}
	for _, lc := range []lifecycle{lifecycleNone, lifecyclePodCreate, lifecyclePodDelete} {
		if got := recordKubeFailure(context.Background(), "test", lc, denied); !strings.Contains(fmt.Sprint(got), "reason=Forbidden") {
			t.Fatal(got)
		}
	}
}
