import * as React from "react";
import { Check, ChevronsUpDown } from "lucide-react";
import { MdDelete } from "react-icons/md";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

export type MultipleSelectorItem = {
  label: string;
  value: string;
};

export type MultipleSelectorProps = {
  items: MultipleSelectorItem[];
  onChange?: (values: string[]) => void;
  search?: string;
  onSearchChange?: (value: string) => void;
  placeholder?: string;
};

export function MultipleSelector({
  items,
  onChange,
  search,
  onSearchChange,
  placeholder = "Select options...",
}: MultipleSelectorProps) {
  const [open, setOpen] = React.useState(false);
  const [selected, setSelected] = React.useState<MultipleSelectorItem[]>([]);

  const selectedVals = React.useMemo(() => {
    return selected.map((s) => s.value);
  }, [selected]);

  const handleSetValue = (val: string) => {
    if (selectedVals.includes(val)) {
      setSelected(selected.filter((item) => item.value !== val));
    } else {
      const item = items.find((s) => s.value === val);
      if (item) {
        setSelected([...selected, item]);
      }
    }
  };

  React.useEffect(() => {
    if (onChange) {
      onChange(selectedVals);
    }
  }, [selectedVals]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-[480px] justify-between min-h-fit group"
        >
          <div className="flex gap-2 justify-start flex-wrap cursor-default">
            {selected.length
              ? selected.map((val) => (
                  <div
                    key={val.value}
                    className="px-2 py-1 rounded-xl border text-xs font-medium bg-card text-card-foreground inline-flex gap-1 items-center"
                  >
                    {val.label}
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleSetValue(val.value);
                      }}
                      className="hover:text-destructive cursor-pointer"
                    >
                      <MdDelete size={14} />
                    </button>
                  </div>
                ))
              : placeholder}
          </div>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[480px] p-0">
        <Command>
          <CommandInput
            value={search}
            onValueChange={onSearchChange}
            placeholder="Search..."
          />
          <CommandEmpty>No results found.</CommandEmpty>
          <CommandGroup>
            <CommandList>
              {items.map((it) => (
                <CommandItem
                  key={it.value}
                  value={it.value}
                  onSelect={() => handleSetValue(it.value)}
                  className="cursor-pointer"
                >
                  <Check
                    className={cn(
                      "mr-2 h-4 w-4",
                      selectedVals.includes(it.value)
                        ? "opacity-100"
                        : "opacity-0",
                    )}
                  />
                  {it.label}
                </CommandItem>
              ))}
            </CommandList>
          </CommandGroup>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
