import type { ContentSection } from "../../api";
import TextRenderer from "./TextSection";
import TableRenderer from "./TableSection";
import CalloutRenderer from "./CalloutSection";
import ImageRenderer from "./ImageSection";
import CodeRenderer from "./CodeSection";
import DiagramRenderer from "./DiagramSection";
import UnknownSection from "./UnknownSection";

interface ContentRendererProps {
  sections: ContentSection[];
}

function renderSection(section: ContentSection, index: number) {
  switch (section.type) {
    case "text":
      return <TextRenderer key={index} section={section} />;
    case "table":
      return <TableRenderer key={index} section={section} />;
    case "callout":
      return <CalloutRenderer key={index} section={section} />;
    case "image":
      return <ImageRenderer key={index} section={section} />;
    case "code":
      return <CodeRenderer key={index} section={section} />;
    case "diagram":
      return <DiagramRenderer key={index} section={section} />;
    default:
      return (
        <UnknownSection
          key={index}
          type={(section as { type: string }).type}
        />
      );
  }
}

export default function ContentRenderer({ sections }: ContentRendererProps) {
  return (
    <div className="space-y-6">
      {sections.map((section, index) => renderSection(section, index))}
    </div>
  );
}
