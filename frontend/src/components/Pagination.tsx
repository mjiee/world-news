import { ActionIcon, Group, Text } from "@mantine/core";
import { IconChevronLeft, IconChevronRight } from "@tabler/icons-react";
import { useEffect, useState } from "react";

export const Pagination = ({ page = 0, total = 0, onChange = (newPage: number) => {}, size = "md" }) => {
  const [inputValue, setInputValue] = useState(page.toString());

  useEffect(() => {
    setInputValue(page.toString());
  }, [page]);

  const handlePrevPage = () => {
    if (page > 1) onChange(page - 1);
  };

  const handleNextPage = () => {
    if (page < total) onChange(page + 1);
  };

  const handleInputChange = () => {
    try {
      const newValue = parseInt(inputValue);
      if (!isNaN(newValue) && newValue >= 1 && newValue <= total) {
        onChange(newValue);
      }
    } catch (_) {
      setInputValue(page.toString());
    }
  };

  return (
    <Group gap="xs" align="center">
      <ActionIcon variant="light" size={size} disabled={page <= 1} onClick={handlePrevPage}>
        <IconChevronLeft />
      </ActionIcon>
      <input
        value={inputValue}
        onChange={(e) => setInputValue(e.target.value)}
        onBlur={handleInputChange}
        onKeyDown={(e) => {
          if (e.key === "Enter") handleInputChange();
        }}
        style={{
          border: "1px solid #ccc",
          borderRadius: "4px",
          textAlign: "center",
          width: "50px",
        }}
      />
      <Text size={size} c="dimmed">
        / {total}
      </Text>
      <ActionIcon variant="light" size={size} disabled={page >= total} onClick={handleNextPage}>
        <IconChevronRight />
      </ActionIcon>
    </Group>
  );
};
