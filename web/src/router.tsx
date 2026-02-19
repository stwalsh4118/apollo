import { createBrowserRouter } from "react-router";
import AppLayout from "./layouts/AppLayout";
import TopicListPage from "./pages/TopicListPage";
import CourseViewPage from "./pages/CourseViewPage";
import NotFoundPage from "./pages/NotFoundPage";

export const router = createBrowserRouter([
  {
    element: <AppLayout />,
    children: [
      { index: true, element: <TopicListPage /> },
      { path: "topics/:id", element: <CourseViewPage /> },
      { path: "*", element: <NotFoundPage /> },
    ],
  },
]);
