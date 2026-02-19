import type {
  CalloutSection as CalloutSectionType,
  CalloutVariant,
} from "../../api";

const VARIANT_STYLES: Record<CalloutVariant, { bg: string; border: string; icon: string }> = {
  prerequisite: {
    bg: "bg-purple-50",
    border: "border-purple-300",
    icon: "M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253",
  },
  warning: {
    bg: "bg-amber-50",
    border: "border-amber-300",
    icon: "M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z",
  },
  tip: {
    bg: "bg-green-50",
    border: "border-green-300",
    icon: "M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z",
  },
  info: {
    bg: "bg-blue-50",
    border: "border-blue-300",
    icon: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
  },
};

const VARIANT_LABELS: Record<CalloutVariant, string> = {
  prerequisite: "Prerequisite",
  warning: "Warning",
  tip: "Tip",
  info: "Info",
};

interface CalloutSectionProps {
  section: CalloutSectionType;
}

export default function CalloutSection({ section }: CalloutSectionProps) {
  const style = VARIANT_STYLES[section.variant] ?? VARIANT_STYLES.info;
  const label = VARIANT_LABELS[section.variant] ?? section.variant;

  return (
    <div
      className={`rounded-lg border-l-4 p-4 ${style.bg} ${style.border}`}
    >
      <div className="flex items-start gap-3">
        <svg
          className="mt-0.5 size-5 shrink-0 opacity-70"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={1.5}
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d={style.icon}
          />
        </svg>
        <div>
          <p className="text-sm font-semibold">{label}</p>
          <p className="mt-1 text-sm">{section.body}</p>
        </div>
      </div>
    </div>
  );
}
