import Markdown from "react-markdown";
import rehypeRaw from "rehype-raw";
import type { TextSection as TextSectionType } from "../../api";

interface TextSectionProps {
  section: TextSectionType;
}

export default function TextSection({ section }: TextSectionProps) {
  return (
    <div className="prose prose-gray max-w-none">
      <Markdown rehypePlugins={[rehypeRaw]}>{section.body}</Markdown>
    </div>
  );
}
