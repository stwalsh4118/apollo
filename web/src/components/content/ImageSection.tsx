import type { ImageSection as ImageSectionType } from "../../api";

interface ImageSectionProps {
  section: ImageSectionType;
}

export default function ImageSection({ section }: ImageSectionProps) {
  return (
    <figure>
      <img
        src={section.url}
        alt={section.alt}
        className="rounded-lg"
        loading="lazy"
      />
      {section.caption && (
        <figcaption className="mt-2 text-center text-sm text-gray-500">
          {section.caption}
        </figcaption>
      )}
    </figure>
  );
}
