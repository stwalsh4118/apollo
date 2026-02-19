import { test, expect } from "@playwright/test";

const BASE_URL = "http://localhost:5173";

// Helper: navigate to a specific lesson by clicking its sidebar entry
async function navigateToLesson(
  page: import("@playwright/test").Page,
  lessonTitle: string,
) {
  await page
    .getByRole("button", { name: lessonTitle, exact: true })
    .click();
  await expect(
    page.locator("h1", { hasText: lessonTitle }),
  ).toBeVisible();
}

test.describe("PBI 3: Frontend Foundation & Course View", () => {
  // AC1: React app builds and serves via Vite dev server
  test("AC1: app loads and serves via dev server", async ({ page }) => {
    await page.goto(BASE_URL);
    await expect(page).toHaveTitle("Apollo");
    await expect(page.locator("header")).toBeVisible();
  });

  // AC2: Topic list page shows topic cards
  test("AC2: topic list shows cards with title, difficulty, description, module count", async ({
    page,
  }) => {
    await page.goto(BASE_URL);
    const card = page.locator('a[href="/topics/e2e-go-basics"]');
    await expect(card).toBeVisible();
    await expect(card).toContainText("Go Fundamentals");
    await expect(card).toContainText("foundational");
    await expect(card).toContainText("comprehensive introduction");
    await expect(card).toContainText("2 modules");
  });

  // AC3: Course view renders module sidebar
  test("AC3: course view shows module sidebar with modules and lessons", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    // Wait for topic data to load — sidebar module titles
    await expect(
      page.locator("aside").getByText("Getting Started with Go"),
    ).toBeVisible();
    await expect(
      page.locator("aside").getByText("Data Types and Control Flow"),
    ).toBeVisible();
    // Lessons should appear in sidebar
    await expect(
      page.getByRole("button", { name: "Introduction to Go", exact: true }),
    ).toBeVisible();
  });

  // AC4: Text section renders markdown
  test("AC4a: text section renders markdown", async ({ page }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();
    // Markdown bold rendering
    await expect(
      page.locator("strong", { hasText: "excellent concurrency" }),
    ).toBeVisible();
  });

  // AC4: Callout sections render all 4 variants
  test("AC4b: callout sections render all 4 variants", async ({ page }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();
    // All 4 callout variants by their label text
    await expect(page.getByText("Info", { exact: true })).toBeVisible();
    await expect(page.getByText("Tip", { exact: true })).toBeVisible();
    await expect(page.getByText("Warning", { exact: true })).toBeVisible();
    await expect(
      page.getByText("Prerequisite", { exact: true }),
    ).toBeVisible();
  });

  // AC4: Table section renders
  test("AC4c: table section renders with headers and rows", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();
    await expect(page.locator("th", { hasText: "Feature" })).toBeVisible();
    await expect(page.locator("td", { hasText: "Goroutines" })).toBeVisible();
  });

  // AC4: Image section renders
  test("AC4d: image section renders with caption", async ({ page }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();
    const img = page.locator('img[alt="Go programming language logo"]');
    await expect(img).toBeVisible();
    await expect(
      page.getByText("The official Go gopher logo"),
    ).toBeVisible();
  });

  // AC4 + AC5: Code section with Shiki highlighting
  test("AC4e + AC5: code section renders with Shiki for all 6 languages", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await navigateToLesson(page, "Your First Go Program");

    // Shiki renders <pre> with class containing "shiki"
    await expect(page.locator(".shiki").first()).toBeVisible({
      timeout: 10000,
    });

    // Code title and explanation should render
    await expect(page.getByText("Hello World in Go")).toBeVisible();
    await expect(
      page.getByText(
        "Every Go program starts with a package declaration",
      ),
    ).toBeVisible();

    // Verify all 6 language badges are present in code section headers
    const codeHeaders = page.locator(".bg-gray-800");
    await expect(codeHeaders).toHaveCount(6);
  });

  // AC6: Mermaid diagrams render
  test("AC6: mermaid diagram renders as SVG", async ({ page }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    // Expand module 2 in sidebar
    await page
      .locator("aside")
      .getByText("Data Types and Control Flow")
      .click();
    await navigateToLesson(page, "Variables and Types");

    // Diagram title should render
    await expect(
      page.getByText("Variable Declaration Flow"),
    ).toBeVisible();

    // Mermaid should render an SVG inside the content area (not the nav icon)
    const diagramSvg = page.locator("article svg");
    await expect(diagramSvg.first()).toBeVisible({ timeout: 10000 });
  });

  // AC7: Exercise hints reveal progressively
  test("AC7: exercises render with progressive hint reveal", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await navigateToLesson(page, "Your First Go Program");

    // Exercises section should be visible
    await expect(page.getByText("Exercises")).toBeVisible();
    await expect(
      page.getByText("Exercise 1: Run Hello World"),
    ).toBeVisible();

    // Hints should be hidden initially
    await expect(page.getByText("Hint 1:")).not.toBeVisible();

    // Click to reveal first hint
    await page.getByText("Show hint 1 of 3").click();
    await expect(page.getByText("Hint 1:")).toBeVisible();
    await expect(page.getByText("Hint 2:")).not.toBeVisible();

    // Click to reveal second hint
    await page.getByText("Show hint 2 of 3").click();
    await expect(page.getByText("Hint 2:")).toBeVisible();
  });

  // AC8: Review questions in collapsible section
  test("AC8: review questions render in collapsible section", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();

    // Review questions section should exist but be collapsed
    const reviewSection = page.getByText("Review Questions (2)");
    await expect(reviewSection).toBeVisible();

    // Click to expand review section
    await reviewSection.click();

    // Questions should now be visible
    const questionLocator = page.locator("details summary", {
      hasText: "Who created the Go programming language?",
    });
    await expect(questionLocator).toBeVisible();

    // Click question to reveal answer
    await questionLocator.click();

    // Answer should be visible inside the review section
    const reviewAnswer = page.locator("details details div", {
      hasText: "Robert Griesemer",
    });
    await expect(reviewAnswer.first()).toBeVisible();
  });

  // AC9: Lesson navigation within and across modules
  test("AC9: prev/next navigation works within and across modules", async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/topics/e2e-go-basics`);
    // Should start on first lesson
    await expect(
      page.locator("h1", { hasText: "Introduction to Go" }),
    ).toBeVisible();

    const navBar = page.locator('nav[aria-label="Lesson navigation"]');

    // First lesson: should have Next but no Previous
    await expect(navBar.getByText("Next")).toBeVisible();
    const prevButtons = navBar.locator("button", { hasText: "Previous" });
    await expect(prevButtons).toHaveCount(0);

    // Navigate to next lesson via nav button
    await navBar.locator("button").last().click();
    await expect(
      page.locator("h1", { hasText: "Your First Go Program" }),
    ).toBeVisible();

    // Middle lesson: should have both Previous and Next
    await expect(navBar.getByText("Previous")).toBeVisible();
    await expect(navBar.getByText("Next")).toBeVisible();

    // Cross-module navigation: click Next to go to Module 2
    await navBar.locator("button", { hasText: "Next" }).click();
    await expect(
      page.locator("h1", { hasText: "Variables and Types" }),
    ).toBeVisible();

    // Last lesson: should have Previous but no Next
    await expect(navBar.getByText("Previous")).toBeVisible();
    const nextButtons = navBar.locator("button", { hasText: "Next" });
    await expect(nextButtons).toHaveCount(0);
  });

  // AC10: API proxy configured
  test("AC10: API proxy forwards /api requests to backend", async ({
    page,
  }) => {
    const response = await page.request.get(`${BASE_URL}/api/health`);
    expect(response.ok()).toBeTruthy();
    const body = await response.json();
    expect(body.status).toBe("ok");
  });
});
