import type { TableSection as TableSectionType } from "../../api";

interface TableSectionProps {
  section: TableSectionType;
}

export default function TableSection({ section }: TableSectionProps) {
  return (
    <div className="overflow-x-auto">
      <table className="min-w-full divide-y divide-gray-200 border border-gray-200 text-sm">
        <thead className="bg-gray-50">
          <tr>
            {section.headers.map((header, i) => (
              <th
                key={i}
                className="px-4 py-2 text-left font-medium text-gray-700"
              >
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200">
          {section.rows.map((row, rowIdx) => (
            <tr
              key={rowIdx}
              className={rowIdx % 2 === 0 ? "bg-white" : "bg-gray-50"}
            >
              {row.map((cell, cellIdx) => (
                <td key={cellIdx} className="px-4 py-2 text-gray-600">
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
