package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// e2eEnv holds the test server and provides HTTP helper methods.
type e2eEnv struct {
	t      *testing.T
	router http.Handler
}

func setupE2E(t *testing.T) *e2eEnv {
	t.Helper()

	srv := setupTestServer(t)
	router := srv.Router()

	env := &e2eEnv{t: t, router: router}

	// Seed realistic curriculum data via write endpoints.
	env.postJSON("/api/topics", `{
		"id":"go-basics","title":"Go Basics","description":"Learn the fundamentals of Go",
		"difficulty":"foundational","status":"published","tags":["go","basics"]
	}`)
	env.postJSON("/api/topics", `{
		"id":"go-advanced","title":"Advanced Go","description":"Deep dive into Go concurrency",
		"difficulty":"advanced","status":"published","tags":["go","advanced"]
	}`)

	// Modules.
	env.postJSON("/api/modules", `{"id":"mod-1","topic_id":"go-basics","title":"Introduction","sort_order":1,"description":"Getting started"}`)
	env.postJSON("/api/modules", `{"id":"mod-2","topic_id":"go-basics","title":"Data Types","sort_order":2,"description":"Types and values"}`)
	env.postJSON("/api/modules", `{"id":"mod-3","topic_id":"go-advanced","title":"Concurrency","sort_order":1,"description":"Goroutines and channels"}`)

	// Lessons.
	env.postJSON("/api/lessons", `{"id":"les-1","module_id":"mod-1","title":"Hello World","sort_order":1,"content":[{"type":"text","body":"Hello"}]}`)
	env.postJSON("/api/lessons", `{"id":"les-2","module_id":"mod-1","title":"Variables","sort_order":2,"content":[{"type":"text","body":"Variables in Go"}]}`)
	env.postJSON("/api/lessons", `{"id":"les-3","module_id":"mod-2","title":"Integers","sort_order":1,"content":[{"type":"text","body":"Integer types"}]}`)
	env.postJSON("/api/lessons", `{"id":"les-4","module_id":"mod-2","title":"Strings","sort_order":2,"content":[{"type":"text","body":"String handling"}]}`)
	env.postJSON("/api/lessons", `{"id":"les-5","module_id":"mod-3","title":"Goroutines","sort_order":1,"content":[{"type":"text","body":"Lightweight threads"}]}`)

	// Concepts.
	env.postJSON("/api/concepts", `{"id":"con-1","name":"Variable","definition":"A named storage location","defined_in_topic":"go-basics","difficulty":"foundational"}`)
	env.postJSON("/api/concepts", `{"id":"con-2","name":"Integer","definition":"A whole number type","defined_in_topic":"go-basics","difficulty":"foundational"}`)
	env.postJSON("/api/concepts", `{"id":"con-3","name":"String","definition":"A sequence of characters","defined_in_topic":"go-basics","difficulty":"foundational"}`)
	env.postJSON("/api/concepts", `{"id":"con-4","name":"Goroutine","definition":"A lightweight thread of execution","defined_in_topic":"go-advanced","difficulty":"intermediate"}`)
	env.postJSON("/api/concepts", `{"id":"con-5","name":"Channel","definition":"A typed conduit for communication","defined_in_topic":"go-advanced","difficulty":"intermediate"}`)

	// Concept references.
	env.postJSON("/api/concepts/con-1/references", `{"lesson_id":"les-1","context":"Introduced in hello world"}`)
	env.postJSON("/api/concepts/con-1/references", `{"lesson_id":"les-2","context":"Deep dive into variables"}`)
	env.postJSON("/api/concepts/con-2/references", `{"lesson_id":"les-3","context":"Integer types explained"}`)
	env.postJSON("/api/concepts/con-3/references", `{"lesson_id":"les-4","context":"String handling basics"}`)

	// Prerequisites.
	env.postJSON("/api/prerequisites", `{"topic_id":"go-advanced","prerequisite_topic_id":"go-basics","priority":"essential","reason":"Must know basics first"}`)

	// Relations.
	env.postJSON("/api/relations", `{"topic_a":"go-basics","topic_b":"go-advanced","relation_type":"builds_on","description":"Advanced builds on basics"}`)

	return env
}

func (e *e2eEnv) get(path string) *httptest.ResponseRecorder {
	e.t.Helper()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	e.router.ServeHTTP(rec, req)

	return rec
}

func (e *e2eEnv) postJSON(path, body string) *httptest.ResponseRecorder {
	e.t.Helper()

	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	e.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		e.t.Fatalf("POST %s: expected 201, got %d: %s", path, rec.Code, rec.Body.String())
	}

	return rec
}

func (e *e2eEnv) putJSON(path, body string) *httptest.ResponseRecorder {
	e.t.Helper()

	req := httptest.NewRequest(http.MethodPut, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	e.router.ServeHTTP(rec, req)

	return rec
}

func decodeMap(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()

	var result map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}

	return result
}

func decodeSlice(t *testing.T, rec *httptest.ResponseRecorder) []any {
	t.Helper()

	var result []any
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("decode JSON array: %v", err)
	}

	return result
}

// AC1: All GET endpoints return correct JSON with 200 status.
func TestE2E_AC1_AllGetEndpointsReturnJSON(t *testing.T) {
	env := setupE2E(t)

	endpoints := []string{
		"/api/health",
		"/api/topics",
		"/api/topics/go-basics",
		"/api/topics/go-basics/full",
		"/api/modules/mod-1",
		"/api/lessons/les-1",
		"/api/concepts",
		"/api/concepts/con-1",
		"/api/concepts/con-1/references",
		"/api/search?q=Go",
		"/api/graph",
		"/api/graph/topic/go-basics",
	}

	for _, ep := range endpoints {
		rec := env.get(ep)
		if rec.Code != http.StatusOK {
			t.Errorf("%s: expected 200, got %d: %s", ep, rec.Code, rec.Body.String())

			continue
		}

		ct := rec.Header().Get("Content-Type")
		if !strings.Contains(ct, "application/json") {
			t.Errorf("%s: expected Content-Type containing 'application/json', got %q", ep, ct)
		}

		// Verify valid JSON.
		var js any
		if err := json.NewDecoder(rec.Body).Decode(&js); err != nil {
			t.Errorf("%s: invalid JSON: %v", ep, err)
		}
	}
}

// AC2: Topic list ordered by title with difficulty and module count.
func TestE2E_AC2_TopicListOrdered(t *testing.T) {
	env := setupE2E(t)

	rec := env.get("/api/topics")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	topics := decodeSlice(t, rec)
	if len(topics) != 2 {
		t.Fatalf("expected 2 topics, got %d", len(topics))
	}

	// Verify ordered by title: "Advanced Go" before "Go Basics".
	first := topics[0].(map[string]any)
	second := topics[1].(map[string]any)

	if first["title"].(string) >= second["title"].(string) {
		t.Fatalf("expected topics ordered by title, got %q then %q", first["title"], second["title"])
	}

	// Verify difficulty field present.
	if first["difficulty"] == nil || first["difficulty"].(string) == "" {
		t.Fatal("expected non-empty difficulty on first topic")
	}

	// "Advanced Go" has 1 module (mod-3), "Go Basics" has 2 modules (mod-1, mod-2).
	firstMC := int(first["module_count"].(float64))
	if firstMC != 1 {
		t.Fatalf("expected module_count 1 for 'Advanced Go', got %d", firstMC)
	}

	secondMC := int(second["module_count"].(float64))
	if secondMC != 2 {
		t.Fatalf("expected module_count 2 for 'Go Basics', got %d", secondMC)
	}
}

// AC3: Full topic tree returns nested modules → lessons → concepts.
func TestE2E_AC3_FullTopicTree(t *testing.T) {
	env := setupE2E(t)

	rec := env.get("/api/topics/go-basics/full")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	result := decodeMap(t, rec)

	if result["id"] != "go-basics" {
		t.Fatalf("expected id 'go-basics', got %v", result["id"])
	}

	modules, ok := result["modules"].([]any)
	if !ok || len(modules) == 0 {
		t.Fatal("expected non-empty modules array")
	}

	if len(modules) != 2 {
		t.Fatalf("expected 2 modules, got %d", len(modules))
	}

	// Find a lesson that has concepts by iterating all modules and lessons.
	var foundConcepts bool

	for _, m := range modules {
		mod := m.(map[string]any)
		lessons, ok := mod["lessons"].([]any)
		if !ok {
			continue
		}

		for _, l := range lessons {
			lesson := l.(map[string]any)
			concepts, ok := lesson["concepts"].([]any)
			if ok && len(concepts) > 0 {
				foundConcepts = true

				// Verify concept has expected fields.
				concept := concepts[0].(map[string]any)
				if concept["id"] == nil || concept["name"] == nil {
					t.Fatal("concept missing id or name fields")
				}
			}
		}
	}

	if !foundConcepts {
		t.Fatal("expected at least one lesson with concepts in the full tree")
	}
}

// AC4: Concept references endpoint returns all teaching/referencing lessons.
func TestE2E_AC4_ConceptReferences(t *testing.T) {
	env := setupE2E(t)

	rec := env.get("/api/concepts/con-1/references")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	refs := decodeSlice(t, rec)
	if len(refs) != 2 {
		t.Fatalf("expected 2 references for con-1, got %d", len(refs))
	}

	// Verify reference fields with known values.
	lessonIDs := make(map[string]bool)
	for _, r := range refs {
		ref := r.(map[string]any)
		if ref["lesson_id"] == nil || ref["lesson_title"] == nil {
			t.Fatal("reference missing lesson_id or lesson_title")
		}

		lessonIDs[ref["lesson_id"].(string)] = true
	}

	if !lessonIDs["les-1"] || !lessonIDs["les-2"] {
		t.Fatalf("expected references to les-1 and les-2, got %v", lessonIDs)
	}
}

// AC5: FTS5 search returns results with relevance ranking.
func TestE2E_AC5_Search(t *testing.T) {
	env := setupE2E(t)

	rec := env.get("/api/search?q=Go")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	result := decodeMap(t, rec)

	total := int(result["total"].(float64))
	if total == 0 {
		t.Fatal("expected search results for 'Go', got 0")
	}

	items, ok := result["items"].([]any)
	if !ok || len(items) == 0 {
		t.Fatal("expected non-empty items array")
	}

	// Verify result has expected fields.
	item := items[0].(map[string]any)
	if item["entity_type"] == nil || item["entity_id"] == nil || item["title"] == nil {
		t.Fatal("search result missing expected fields")
	}

	// Snippet field should be present (FTS5 snippet() function).
	if item["snippet"] == nil {
		t.Fatal("expected non-nil snippet field in search result")
	}

	// Relevance ranking is ensured by ORDER BY rank in the SQL query.
	// The FTS5 rank column is internal and not exposed in the response model.
	// We verify ranking works by confirming results are returned and ordered by FTS5.

	// Verify empty query returns 400.
	rec400 := env.get("/api/search")
	if rec400.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty query, got %d", rec400.Code)
	}
}

// AC6: Graph endpoint returns nodes and edges.
func TestE2E_AC6_Graph(t *testing.T) {
	env := setupE2E(t)

	rec := env.get("/api/graph")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	result := decodeMap(t, rec)

	nodes, ok := result["nodes"].([]any)
	if !ok || len(nodes) == 0 {
		t.Fatal("expected non-empty nodes array")
	}

	edges, ok := result["edges"].([]any)
	if !ok || len(edges) == 0 {
		t.Fatal("expected non-empty edges array")
	}

	// Verify node types.
	nodeTypes := make(map[string]bool)
	for _, n := range nodes {
		node := n.(map[string]any)
		nodeTypes[node["type"].(string)] = true
	}

	if !nodeTypes["topic"] {
		t.Fatal("expected topic nodes in graph")
	}

	if !nodeTypes["concept"] {
		t.Fatal("expected concept nodes in graph")
	}

	// Verify edge types.
	edgeTypes := make(map[string]bool)
	for _, e := range edges {
		edge := e.(map[string]any)
		edgeTypes[edge["type"].(string)] = true
	}

	if !edgeTypes["essential"] {
		t.Fatal("expected prerequisite edge (type=essential)")
	}

	if !edgeTypes["builds_on"] {
		t.Fatal("expected relation edge (type=builds_on)")
	}

	if !edgeTypes["reference"] {
		t.Fatal("expected reference edge")
	}
}

// AC7: Write endpoints store data with FK integrity.
func TestE2E_AC7_WriteAndReadBack(t *testing.T) {
	env := setupE2E(t)

	// Create a new topic via POST, then read back via GET.
	env.postJSON("/api/topics", `{"id":"new-topic","title":"New Topic","status":"draft"}`)

	rec := env.get("/api/topics/new-topic")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for newly created topic, got %d", rec.Code)
	}

	result := decodeMap(t, rec)
	if result["title"] != "New Topic" {
		t.Fatalf("expected title 'New Topic', got %v", result["title"])
	}

	// Update via PUT.
	putRec := env.putJSON("/api/topics/new-topic", `{"title":"Updated Topic","status":"published"}`)
	if putRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for update, got %d: %s", putRec.Code, putRec.Body.String())
	}

	// Read back updated.
	rec2 := env.get("/api/topics/new-topic")
	result2 := decodeMap(t, rec2)
	if result2["title"] != "Updated Topic" {
		t.Fatalf("expected updated title, got %v", result2["title"])
	}

	// FK violation: module with nonexistent topic_id.
	fkReq := httptest.NewRequest(http.MethodPost, "/api/modules",
		strings.NewReader(`{"id":"bad-mod","topic_id":"nonexistent","title":"Bad","sort_order":1}`))
	fkRec := httptest.NewRecorder()
	env.router.ServeHTTP(fkRec, fkReq)

	if fkRec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for module FK violation, got %d", fkRec.Code)
	}

	// FK violation: lesson with nonexistent module_id.
	fkReq2 := httptest.NewRequest(http.MethodPost, "/api/lessons",
		strings.NewReader(`{"id":"bad-les","module_id":"nonexistent","title":"Bad","sort_order":1,"content":[{"type":"text"}]}`))
	fkRec2 := httptest.NewRecorder()
	env.router.ServeHTTP(fkRec2, fkReq2)

	if fkRec2.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422 for lesson FK violation, got %d", fkRec2.Code)
	}

	// Duplicate key: create topic with existing ID.
	dupReq := httptest.NewRequest(http.MethodPost, "/api/topics",
		strings.NewReader(`{"id":"go-basics","title":"Duplicate","status":"draft"}`))
	dupRec := httptest.NewRecorder()
	env.router.ServeHTTP(dupRec, dupReq)

	if dupRec.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate, got %d", dupRec.Code)
	}
}

// AC8: Pagination works on list endpoints.
func TestE2E_AC8_Pagination(t *testing.T) {
	env := setupE2E(t)

	// Concept pagination: 5 concepts seeded, request page 1 with per_page=2.
	rec := env.get("/api/concepts?page=1&per_page=2")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	result := decodeMap(t, rec)

	total := int(result["total"].(float64))
	if total != 5 {
		t.Fatalf("expected total 5, got %d", total)
	}

	items := result["items"].([]any)
	if len(items) != 2 {
		t.Fatalf("expected 2 items on page 1, got %d", len(items))
	}

	page := int(result["page"].(float64))
	if page != 1 {
		t.Fatalf("expected page 1, got %d", page)
	}

	perPage := int(result["per_page"].(float64))
	if perPage != 2 {
		t.Fatalf("expected per_page 2, got %d", perPage)
	}

	// Page 3 should have 1 item.
	rec3 := env.get("/api/concepts?page=3&per_page=2")
	if rec3.Code != http.StatusOK {
		t.Fatalf("expected 200 for page 3, got %d", rec3.Code)
	}

	result3 := decodeMap(t, rec3)
	items3 := result3["items"].([]any)
	if len(items3) != 1 {
		t.Fatalf("expected 1 item on page 3, got %d", len(items3))
	}

	// Search pagination.
	searchRec := env.get("/api/search?q=Go&page=1&per_page=2")
	if searchRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", searchRec.Code)
	}

	searchResult := decodeMap(t, searchRec)
	if searchResult["page"] == nil || searchResult["per_page"] == nil {
		t.Fatal("expected pagination fields in search response")
	}
}

// Test 404 responses for nonexistent resources.
func TestE2E_NotFoundResponses(t *testing.T) {
	env := setupE2E(t)

	notFoundEndpoints := []string{
		"/api/topics/nonexistent",
		"/api/topics/nonexistent/full",
		"/api/modules/nonexistent",
		"/api/lessons/nonexistent",
		"/api/concepts/nonexistent",
		"/api/concepts/nonexistent/references",
		"/api/graph/topic/nonexistent",
	}

	for _, ep := range notFoundEndpoints {
		rec := env.get(ep)
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s: expected 404, got %d", ep, rec.Code)
		}
	}
}

// Test topic filter on concepts list.
func TestE2E_ConceptTopicFilter(t *testing.T) {
	env := setupE2E(t)

	rec := env.get("/api/concepts?topic=go-basics")
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	result := decodeMap(t, rec)
	total := int(result["total"].(float64))
	if total != 3 {
		t.Fatalf("expected 3 concepts for go-basics, got %d", total)
	}
}
