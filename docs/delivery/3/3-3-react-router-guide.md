# 3-3 React Router v7 Guide

Date: 2026-02-18

## Package

- `react-router` (v7) -- single unified package, replaces both `react-router` and `react-router-dom` from v6

## Documentation

- https://reactrouter.com/
- https://reactrouter.com/start/data/installation
- https://reactrouter.com/upgrading/v6

---

## 1. Installation

```bash
pnpm add react-router
```

Requires React 18+ and Node 20+.

**v7 change:** The `react-router-dom` package is no longer needed. Everything is imported from `react-router`, with one exception: `RouterProvider` is imported from `react-router/dom`.

---

## 2. Router Setup with createBrowserRouter

Create a data router using `createBrowserRouter` and render it with `RouterProvider`.

```tsx
import React from "react";
import ReactDOM from "react-dom/client";
import { createBrowserRouter } from "react-router";
import { RouterProvider } from "react-router/dom";

const router = createBrowserRouter([
  {
    path: "/",
    Component: RootLayout,
    children: [
      { index: true, Component: Home },
      { path: "about", Component: About },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
```

**Key points:**

- `createBrowserRouter` accepts an array of route objects.
- `RouterProvider` is the only component that needs the `/dom` deep import.
- Routes use `Component` (capitalized) to reference the component directly, or `element` to pass JSX.

---

## 3. Defining Routes

### Nested Routes with Children

Parent paths are automatically prepended to children. Child routes render through `<Outlet />` in the parent component.

```tsx
import { createBrowserRouter } from "react-router";

const router = createBrowserRouter([
  {
    path: "/",
    Component: RootLayout,
    children: [
      { index: true, Component: Home },
      { path: "about", Component: About },
      {
        path: "courses",
        Component: CoursesLayout,
        children: [
          { index: true, Component: CourseList },
          { path: ":courseId", Component: CourseDetail },
        ],
      },
    ],
  },
]);
```

This produces the following URL structure:

| URL | Component |
|-----|-----------|
| `/` | `RootLayout` > `Home` |
| `/about` | `RootLayout` > `About` |
| `/courses` | `RootLayout` > `CoursesLayout` > `CourseList` |
| `/courses/42` | `RootLayout` > `CoursesLayout` > `CourseDetail` |

### Index Routes

An index route (`index: true`) renders at the parent's exact URL as the default child. It cannot have `children`.

```tsx
{ index: true, Component: Home }
```

### Layout Routes (No Path)

Omit `path` to create a layout wrapper that does not add a URL segment:

```tsx
{
  // No path -- purely a layout wrapper
  Component: AuthLayout,
  children: [
    { path: "login", Component: Login },
    { path: "register", Component: Register },
  ],
}
```

### Prefix Routes (No Component)

Omit `Component` to group routes under a shared path prefix without introducing a layout:

```tsx
{
  path: "settings",
  children: [
    { index: true, Component: SettingsHome },
    { path: "profile", Component: ProfileSettings },
    { path: "account", Component: AccountSettings },
  ],
}
```

---

## 4. Outlet for Layout Nesting

Parent route components use `<Outlet />` to indicate where child route content renders.

```tsx
import { Outlet } from "react-router";

function RootLayout() {
  return (
    <div>
      <header>
        <nav>{/* navigation links */}</nav>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  );
}
```

```tsx
import { Outlet } from "react-router";

function CoursesLayout() {
  return (
    <div>
      <h1>Courses</h1>
      <Outlet />
    </div>
  );
}
```

The `<Outlet />` renders the matched child route component. If no child matches, it renders nothing (or the index route if one exists).

---

## 5. Navigation: Link, NavLink, useParams

### Link

Client-side navigation without full page reload.

```tsx
import { Link } from "react-router";

function Navigation() {
  return (
    <nav>
      <Link to="/">Home</Link>
      <Link to="/courses">Courses</Link>
      <Link to="/about">About</Link>
    </nav>
  );
}
```

`to` accepts a string or an object:

```tsx
<Link to={{ pathname: "/courses", search: "?sort=title" }} />
```

### NavLink

Like `Link` but with built-in active/pending state awareness. Automatically applies `active`, `pending`, and `transitioning` CSS classes, and sets `aria-current="page"` when active.

```tsx
import { NavLink } from "react-router";

function Sidebar() {
  return (
    <nav>
      <NavLink to="/" end>
        Home
      </NavLink>
      <NavLink to="/courses">Courses</NavLink>
      <NavLink to="/about">About</NavLink>
    </nav>
  );
}
```

**CSS styling with default classes:**

```css
nav a.active {
  font-weight: bold;
  color: var(--color-primary);
}
nav a.pending {
  opacity: 0.6;
}
```

**Callback-based className:**

```tsx
<NavLink
  to="/courses"
  className={({ isActive, isPending }) =>
    isPending ? "pending" : isActive ? "active" : ""
  }
>
  Courses
</NavLink>
```

**The `end` prop:** By default, `NavLink` is active when the URL starts with its `to` value. Adding `end` restricts the match to exact URLs only. This is important for the root `/` link to prevent it being active on every page.

| Link | URL | isActive |
|------|-----|----------|
| `<NavLink to="/courses" />` | `/courses/42` | true |
| `<NavLink to="/courses" end />` | `/courses/42` | false |
| `<NavLink to="/courses" end />` | `/courses` | true |

### useParams

Access dynamic route parameters in components.

```tsx
import { useParams } from "react-router";

function CourseDetail() {
  const { courseId } = useParams();
  // courseId is string | undefined
  return <div>Course ID: {courseId}</div>;
}
```

For the route `path: ":courseId"`, the `:courseId` segment is parsed from the URL and available via `useParams()`.

### useNavigate

Programmatic navigation (redirects, post-action navigation):

```tsx
import { useNavigate } from "react-router";

function LoginForm() {
  const navigate = useNavigate();

  async function handleSubmit() {
    await login();
    navigate("/dashboard");
  }

  return <form onSubmit={handleSubmit}>{/* ... */}</form>;
}
```

---

## 6. Catch-All / 404 Routes

Use a splat route (`path: "*"`) to catch all unmatched URLs. Place it as the last child so it only matches when no other route does.

```tsx
const router = createBrowserRouter([
  {
    path: "/",
    Component: RootLayout,
    children: [
      { index: true, Component: Home },
      { path: "courses", Component: CourseList },
      { path: "courses/:courseId", Component: CourseDetail },
      { path: "*", Component: NotFound },
    ],
  },
]);
```

```tsx
import { Link } from "react-router";

function NotFound() {
  return (
    <div>
      <h1>404 - Page Not Found</h1>
      <p>The page you are looking for does not exist.</p>
      <Link to="/">Go Home</Link>
    </div>
  );
}
```

The splat value is accessible via `params["*"]`:

```tsx
import { useParams } from "react-router";

function NotFound() {
  const { "*": splat } = useParams();
  // splat contains the unmatched path segments
  return <p>No match for: /{splat}</p>;
}
```

---

## 7. Key v7 Changes from v6

| Change | v6 | v7 |
|--------|----|----|
| **Package** | `react-router-dom` | `react-router` (single package) |
| **RouterProvider import** | `from "react-router-dom"` | `from "react-router/dom"` |
| **All other imports** | `from "react-router-dom"` | `from "react-router"` |
| **json() helper** | `import { json } from "react-router-dom"` | Removed; use `Response.json()` or return plain objects from loaders |
| **defer() helper** | `import { defer } from "react-router-dom"` | Removed; return promises directly |
| **formMethod casing** | Lowercase (`"post"`, `"get"`) | Uppercase (`"POST"`, `"GET"`) |
| **Splat path matching** | Relative links resolve from route path | Relative links resolve from URL path (use `..` prefix) |
| **Router state updates** | `React.useState` | `React.useTransition` (concurrent mode) |
| **Bundle size** | Baseline | ~15% smaller |
| **Min Node** | 14+ | 20+ |
| **Min React** | 16.8+ | 18+ |

### Migration Path

If migrating from v6, enable future flags one at a time in v6 before upgrading:

- `v7_relativeSplatPath`
- `v7_startTransition`
- `v7_fetcherPersist`
- `v7_normalizeFormMethod`
- `v7_partialHydration`
- `v7_skipActionErrorRevalidation`

Once all flags are enabled and tests pass, upgrade to v7 with zero breaking changes.

---

## Apollo Usage Summary

For task 3-3, the key integration points are:

1. **Install**: `pnpm add react-router`
2. **Router**: Create with `createBrowserRouter` in a dedicated routes file
3. **Provider**: Render `<RouterProvider>` (from `react-router/dom`) in `main.tsx`
4. **Layouts**: Use `<Outlet />` in layout components for nested route rendering
5. **Navigation**: Use `<NavLink>` for sidebar/nav links (auto-active styling), `<Link>` elsewhere
6. **Params**: Use `useParams()` for dynamic route segments (e.g., `:courseId`)
7. **404**: Add `{ path: "*", Component: NotFound }` as the last child route
