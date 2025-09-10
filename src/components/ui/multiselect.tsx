import * as React from "react";
import { Check, ChevronsUpDown, Search } from "lucide-react";
import { MdDelete } from "react-icons/md";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Input } from "@/components/ui/input";

export type MultipleSelectorItem = {
  label: string;
  value: string;
};

const multipleSelectorVariants = cva("w-full justify-between min-h-fit group", {
  variants: {
    variant: {
      default:
        "border border-foreground bg-background shadow-xs hover:bg-muted hover:text-muted-foreground",
      outline:
        "border border-border bg-background hover:bg-muted hover:text-muted-foreground",
      ghost: "bg-transparent hover:bg-muted",
    },
    size: {
      default: "h-8 px-3 py-2",
      sm: "h-7 px-2 text-sm rounded-md",
      lg: "h-9 px-5 rounded-md text-base",
    },
  },
  defaultVariants: {
    variant: "outline",
    size: "default",
  },
});

export type MultipleSelectorProps = {
  items: MultipleSelectorItem[];
  onChange?: (values: string[]) => void;
  search?: string;
  onSearchChange?: (value: string) => void;
  placeholder?: string;
  variant?: VariantProps<typeof multipleSelectorVariants>["variant"];
  size?: VariantProps<typeof multipleSelectorVariants>["size"];
  fullWidth?: boolean;
};

export function MultipleSelector({
  items,
  onChange,
  search,
  onSearchChange,
  placeholder = "Select options...",
  variant,
  size,
  fullWidth,
}: MultipleSelectorProps) {
  const [open, setOpen] = React.useState(false);
  const [selected, setSelected] = React.useState<MultipleSelectorItem[]>([]);

  const selectedVals = React.useMemo(
    () => selected.map((s) => s.value),
    [selected],
  );

  const handleSetValue = (val: string) => {
    if (selectedVals.includes(val)) {
      setSelected((prev) => prev.filter((item) => item.value !== val));
    } else {
      const item = items.find((s) => s.value === val);
      if (item) setSelected((prev) => [...prev, item]);
    }
  };

  React.useEffect(() => {
    onChange?.(selected.map((s) => s.value));
  }, [selected]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          className={cn(multipleSelectorVariants({ variant, size }))}
          role="combobox"
          aria-expanded={open}
        >
          <div className="flex gap-2 justify-start flex-wrap cursor-default">
            {selected.length
              ? selected.map((val) => (
                  <div
                    key={val.value}
                    className="px-2 py-1 rounded-xl border text-xs font-medium bg-card text-card-foreground inline-flex gap-1 items-center"
                  >
                    {val.label}
                    <div
                      onClick={(e) => {
                        e.stopPropagation();
                        handleSetValue(val.value);
                      }}
                      className="hover:text-destructive cursor-pointer"
                    >
                      <MdDelete size={14} />
                    </div>
                  </div>
                ))
              : placeholder}
          </div>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className={cn("p-0", fullWidth ? "max-w-80 w-[calc(100vw-200px)]" : "")}
      >
        <div className="flex flex-col px-3 py-3 gap-1">
          <div className="bg-muted flex items-center px-3 rounded-xl mb-3">
            <Search />
            <Input
              value={search}
              onChange={(e) => onSearchChange?.(e.target.value)}
              placeholder="Search..."
              className="bg-muted border-0 !outline-none !shadow-none !ring-0"
            />
          </div>
          {items.map((it) => (
            <div
              key={it.value}
              onClick={() => handleSetValue(it.value)}
              className="cursor-pointer flex items-center px-4 hover:bg-accent rounded-xl"
            >
              <Check
                className={cn(
                  "mr-2 h-4 w-4",
                  selectedVals.includes(it.value) ? "opacity-100" : "opacity-0",
                )}
              />
              {it.label}
            </div>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  );
}
