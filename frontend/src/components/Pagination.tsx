import { useState } from "react";
import { Group, Text, ActionIcon } from "@mantine/core";
import IconChevronLeft from "@/assets/icons/IconChevronLeft.svg?react";
import IconChevronRight from "@/assets/icons/IconChevronRight.svg?react";

export const Pagination = ({ value = 0, total = 0, onChange = (value: number) => {}, size = "md" }) => {
  const [inputValue, setInputValue] = useState(value.toString());

  const handlePrevPage = () => {
    if (value > 1) onChange(value - 1);
  };

  const handleNextPage = () => {
    if (value < total) onChange(value + 1);
  };

  const handleInputChange = () => {
    try {
      const newValue = parseInt(inputValue);
      if (!isNaN(newValue) && newValue >= 1 && newValue <= total) {
        onChange(newValue);
      }
    } catch (_) {
      setInputValue(value.toString());
    }
  };

  return (
    <Group gap="xs" align="center">
      <ActionIcon variant="light" size={size} disabled={value <= 1} onClick={handlePrevPage}>
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
      <ActionIcon variant="light" size={size} disabled={value >= total} onClick={handleNextPage}>
        <IconChevronRight />
      </ActionIcon>
    </Group>
  );
};
