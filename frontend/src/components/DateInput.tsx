import { DateInput as MantineDateInput } from "@mantine/dates";
import { useTranslation } from "react-i18next";
import "dayjs/locale/en";
import "dayjs/locale/zh";
import { useState } from "react";

export interface DateInputProps {
  label?: string;
  placeholder?: string;
  onChange: (date: string | null) => void;
}

export function DateInput(props: DateInputProps) {
  const [value, setValue] = useState<string | null>(null);
  const { t, i18n } = useTranslation();

  return (
    <MantineDateInput
      maxDate={new Date()}
      locale={i18n.language}
      label={props.label}
      placeholder={props.placeholder}
      valueFormat="YYYY-MM-DD"
      value={value}
      onChange={(date) => {
        setValue(date);
        props.onChange(date);
      }}
    />
  );
}
