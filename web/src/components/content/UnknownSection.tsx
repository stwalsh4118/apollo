interface UnknownSectionProps {
  type: string;
}

export default function UnknownSection({ type }: UnknownSectionProps) {
  return (
    <div className="rounded-md border border-dashed border-gray-300 bg-gray-50 p-4 text-sm text-gray-500">
      Unsupported content section: {type}
    </div>
  );
}
