import { test, expect } from "@playwright/test";

const BASE_URL = "http://localhost:5173";
const API_URL = "http://localhost:5173/api";

// Seed concept data for concept chip tests.
// The e2e-go-basics topic and its lessons are seeded by seed-e2e.sh.
// We add concepts and link them to e2e-les-1.
test.describe("PBI 4: Learning Progress & Notes", () => {
  test.beforeAll(async ({ request }) => {
    // Create two concepts (201 created, or 409 conflict if already seeded)
    for (const concept of [
      {
        id: "e2e-concept-goroutines",
        name: "Goroutines",
        definition: "Lightweight threads managed by the Go runtime.",
        status: "active",
      },
      {
        id: "e2e-concept-channels",
        name: "Channels",
        definition: "Typed conduits for communication between goroutines.",
        status: "active",
      },
    ]) {
      const res = await request.post(`${API_URL}/concepts`, {
        data: concept,
      });
      expect(
        res.ok() || res.status() === 409,
        `Failed to seed concept ${concept.id}: ${res.status()}`,
      ).toBeTruthy();
    }

    // Link concepts to e2e-les-1
    for (const conceptId of [
      "e2e-concept-goroutines",
      "e2e-concept-channels",
    ]) {
      const res = await request.post(
        `${API_URL}/concepts/${conceptId}/references`,
        { data: { lesson_id: "e2e-les-1" } },
      );
      expect(
        res.ok() || res.status() === 409,
        `Failed to link concept ${conceptId} to e2e-les-1: ${res.status()}`,
      ).toBeTruthy();
    }
  });

  // ---------------------------------------------------------------
  // AC1: PUT /api/progress/lessons/:id stores status and notes
  // ---------------------------------------------------------------
  test("AC1: PUT progress stores status and notes", async ({ request }) => {
    const res = await request.put(`${API_URL}/progress/lessons/e2e-les-1`, {
      data: { status: "completed", notes: "Great intro lesson" },
    });
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.lesson_id).toBe("e2e-les-1");
    expect(body.status).toBe("completed");
    expect(body.notes).toBe("Great intro lesson");
    expect(body.completed_at).toBeTruthy();
  });

  test("AC1: PUT progress rejects invalid status", async ({ request }) => {
    const res = await request.put(`${API_URL}/progress/lessons/e2e-les-1`, {
      data: { status: "invalid_status" },
    });
    expect(res.status()).toBe(400);
  });

  test("AC1: PUT progress returns 404 for non-existent lesson", async ({
    request,
  }) => {
    const res = await request.put(
      `${API_URL}/progress/lessons/non-existent-lesson`,
      { data: { status: "completed" } },
    );
    expect(res.status()).toBe(404);
  });

  // ---------------------------------------------------------------
  // AC2: GET /api/progress/topics/:id returns per-lesson status
  // ---------------------------------------------------------------
  test("AC2: GET topic progress returns per-lesson status", async ({
    request,
  }) => {
    // Explicitly set progress for known state (test isolation)
    await request.put(`${API_URL}/progress/lessons/e2e-les-1`, {
      data: { status: "completed" },
    });
    await request.put(`${API_URL}/progress/lessons/e2e-les-2`, {
      data: { status: "in_progress" },
    });

    const res = await request.get(
      `${API_URL}/progress/topics/e2e-go-basics`,
    );
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(body.topic_id).toBe("e2e-go-basics");
    expect(Array.isArray(body.lessons)).toBeTruthy();
    // Should contain all 3 lessons
    expect(body.lessons.length).toBe(3);

    const lessonIds = body.lessons.map(
      (l: { lesson_id: string }) => l.lesson_id,
    );
    expect(lessonIds).toContain("e2e-les-1");
    expect(lessonIds).toContain("e2e-les-2");
    expect(lessonIds).toContain("e2e-les-3");
  });

  // ---------------------------------------------------------------
  // AC3: GET /api/progress/summary returns completion % and active topics
  // ---------------------------------------------------------------
  test("AC3: GET progress summary returns completion percentage and active topics", async ({
    request,
  }) => {
    const res = await request.get(`${API_URL}/progress/summary`);
    expect(res.ok()).toBeTruthy();
    const body = await res.json();
    expect(typeof body.total_lessons).toBe("number");
    expect(typeof body.completed_lessons).toBe("number");
    expect(typeof body.completion_percentage).toBe("number");
    expect(typeof body.active_topics).toBe("number");
    expect(body.total_lessons).toBeGreaterThan(0);
    expect(body.active_topics).toBeGreaterThanOrEqual(1);
  });

  // ---------------------------------------------------------------
  // AC4: Mark Complete button updates status and shows confirmation
  // ---------------------------------------------------------------
  test("AC4: Mark Complete button updates status and shows visual confirmation", async ({
    page,
    request,
  }) => {
    // Reset lesson 3 to not_started so we can test the Mark Complete flow
    await request.put(`${API_URL}/progress/lessons/e2e-les-3`, {
      data: { status: "not_started" },
    });

    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    // Navigate to lesson 3 (Module 2)
    await page
      .locator("aside")
      .getByText("Data Types and Control Flow")
      .click();
    await page
      .getByRole("button", { name: "Variables and Types", exact: true })
      .click();
    await expect(
      page.locator("h1", { hasText: "Variables and Types" }),
    ).toBeVisible();

    // Should see Mark Complete button (not yet completed)
    const markBtn = page.getByRole("button", { name: "Mark Complete" });
    await expect(markBtn).toBeVisible();

    // Click Mark Complete
    await markBtn.click();

    // Wait for Completed indicator to appear
    await expect(page.getByText("Completed")).toBeVisible({ timeout: 5000 });
  });

  // ---------------------------------------------------------------
  // AC5: Module sidebar shows completion indicators
  // ---------------------------------------------------------------
  test("AC5: sidebar shows checkmark for completed lessons and dot for in-progress", async ({
    page,
    request,
  }) => {
    // Set up: les-1 completed, les-2 in_progress
    await request.put(`${API_URL}/progress/lessons/e2e-les-1`, {
      data: { status: "completed" },
    });
    await request.put(`${API_URL}/progress/lessons/e2e-les-2`, {
      data: { status: "in_progress" },
    });

    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("aside").getByText("Getting Started with Go"),
    ).toBeVisible();

    // Completed lesson should have a green checkmark SVG
    const completedBtn = page.getByRole("button", {
      name: "Introduction to Go",
      exact: true,
    });
    await expect(completedBtn).toBeVisible();
    await expect(completedBtn.locator("svg.text-green-500")).toBeVisible();

    // In-progress lesson should have a blue filled dot
    const inProgressBtn = page.getByRole("button", {
      name: "Your First Go Program",
      exact: true,
    });
    await expect(inProgressBtn).toBeVisible();
    await expect(inProgressBtn.locator("span.bg-blue-400")).toBeVisible();
  });

  // ---------------------------------------------------------------
  // AC6: Progress bar per module shows correct ratio
  // ---------------------------------------------------------------
  test("AC6: module progress bar shows correct completion ratio", async ({
    page,
    request,
  }) => {
    // Set up: les-1 completed, les-2 not started (Module 1 has 2 lessons)
    await request.put(`${API_URL}/progress/lessons/e2e-les-1`, {
      data: { status: "completed" },
    });
    await request.put(`${API_URL}/progress/lessons/e2e-les-2`, {
      data: { status: "not_started" },
    });

    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("aside").getByText("Getting Started with Go"),
    ).toBeVisible();

    // Module 1 should show "1/2" progress label
    const sidebar = page.locator("aside");
    await expect(sidebar.getByText("1/2")).toBeVisible();
  });

  // ---------------------------------------------------------------
  // AC7: Personal notes persist across page reloads
  // ---------------------------------------------------------------
  test("AC7: personal notes persist across page reloads", async ({
    page,
    request,
  }) => {
    // Clear any existing notes first
    await request.put(`${API_URL}/progress/lessons/e2e-les-1`, {
      data: { status: "completed", notes: "" },
    });

    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();

    // Type notes
    const textarea = page.getByPlaceholder(
      "Add your notes for this lesson...",
    );
    await expect(textarea).toBeVisible();
    await textarea.fill("These are my E2E test notes");

    // Click Save Notes
    await page.getByRole("button", { name: "Save Notes" }).click();

    // Wait for the "Saved" confirmation
    await expect(page.getByText("Saved")).toBeVisible({ timeout: 5000 });

    // Reload the page
    await page.reload();

    // Wait for lesson content to load again
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();

    // Notes should persist
    const reloadedTextarea = page.getByPlaceholder(
      "Add your notes for this lesson...",
    );
    await expect(reloadedTextarea).toHaveValue("These are my E2E test notes");
  });

  // ---------------------------------------------------------------
  // AC8: Concept chips render for concepts in the current lesson
  // ---------------------------------------------------------------
  test("AC8: concept chips render for lessons with concepts", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();

    // Concept chips should be visible for e2e-les-1
    await expect(
      page.getByRole("link", { name: "Goroutines" }),
    ).toBeVisible();
    await expect(
      page.getByRole("link", { name: "Channels" }),
    ).toBeVisible();
  });

  test("AC8: concept chips do not render for lessons without concepts", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    // Navigate to lesson 2 which has no concepts
    await page
      .getByRole("button", { name: "Your First Go Program", exact: true })
      .click();
    await expect(
      page.locator("h1", { hasText: "Your First Go Program" }),
    ).toBeVisible();

    // Should NOT see concept chip links
    await expect(
      page.getByRole("link", { name: "Goroutines" }),
    ).not.toBeVisible();
    await expect(
      page.getByRole("link", { name: "Channels" }),
    ).not.toBeVisible();
  });

  // ---------------------------------------------------------------
  // AC9: Clicking a concept chip navigates to /concepts/:id
  // ---------------------------------------------------------------
  test("AC9: clicking a concept chip navigates to concept detail route", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();

    // Click the Goroutines chip
    await page.getByRole("link", { name: "Goroutines" }).click();

    // Should navigate to concept detail page
    await expect(page).toHaveURL(/\/concepts\/e2e-concept-goroutines/);

    // Concept detail placeholder page should render
    await expect(page.getByText("Concept Detail")).toBeVisible();
    await expect(
      page.locator("span.font-mono", {
        hasText: "e2e-concept-goroutines",
      }),
    ).toBeVisible();
  });
});
