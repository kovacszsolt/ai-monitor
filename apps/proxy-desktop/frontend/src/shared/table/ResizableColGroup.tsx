type Props = {
  columnIds: readonly string[];
  widths: Record<string, number>;
};

export function ResizableColGroup({ columnIds, widths }: Props) {
  return (
    <colgroup>
      {columnIds.map((id) => (
        <col key={id} style={{ width: widths[id] }} />
      ))}
    </colgroup>
  );
}
