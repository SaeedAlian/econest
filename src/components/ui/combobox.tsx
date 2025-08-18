import * as React from "react";
import { CheckIcon, ChevronsUpDownIcon } from "lucide-react";
import { cva, type VariantProps } from "class-variance-authority";

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

export interface ComboboxOption {
  value: string;
  label: string;
  disabled?: boolean;
}

const comboboxVariants = cva("flex items-center justify-between w-full", {
  variants: {
    variant: {
      default:
        "bg-background border border-input hover:bg-accent hover:text-accent-foreground",
      filled: "bg-accent border border-accent hover:bg-accent/80",
      ghost: "border-0 hover:bg-accent hover:text-accent-foreground",
    },
    size: {
      default: "h-9 px-4 py-2",
      sm: "h-8 px-3 text-sm",
      lg: "h-10 px-6",
    },
  },
  defaultVariants: {
    variant: "default",
    size: "default",
  },
});

export interface ComboboxProps
  extends React.ComponentProps<typeof PopoverTrigger>,
    VariantProps<typeof comboboxVariants> {
  options: ComboboxOption[];
  value?: string;
  defaultValue?: string;
  onValueChange?: (value: string) => void;
  placeholder?: string;
  searchPlaceholder?: string;
  emptyText?: string;
  disabled?: boolean;
  className?: string;
  buttonClassName?: string;
  contentClassName?: string;
}

const Combobox = React.forwardRef<HTMLButtonElement, ComboboxProps>(
  (
    {
      options,
      value: valueProp,
      defaultValue = "",
      onValueChange,
      placeholder = "Select an option...",
      searchPlaceholder = "Search...",
      emptyText = "No results found.",
      variant,
      size,
      disabled = false,
      className,
      buttonClassName,
      contentClassName,
      ...props
    },
    ref,
  ) => {
    const [open, setOpen] = React.useState(false);
    const [internalValue, setInternalValue] = React.useState(defaultValue);

    const value = valueProp !== undefined ? valueProp : internalValue;
    const setValue = (newValue: string) => {
      if (valueProp === undefined) {
        setInternalValue(newValue);
      }
      onValueChange?.(newValue);
    };

    const selectedOption = options.find((option) => option.value === value);

    return (
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            ref={ref}
            variant="ghost"
            role="combobox"
            aria-expanded={open}
            disabled={disabled}
            className={cn(
              comboboxVariants({ variant, size, className }),
              buttonClassName,
            )}
            {...props}
          >
            <span className="truncate">
              {selectedOption?.label || placeholder}
            </span>
            <ChevronsUpDownIcon className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent
          className={cn(
            "w-[var(--radix-popover-trigger-width)] p-0",
            contentClassName,
          )}
          align="start"
        >
          <Command>
            <CommandInput placeholder={searchPlaceholder} />
            <CommandEmpty>{emptyText}</CommandEmpty>
            <CommandGroup>
              <CommandList>
                {options.map((option) => (
                  <CommandItem
                    key={option.value}
                    value={option.value}
                    disabled={option.disabled}
                    onSelect={() => {
                      setValue(option.value === value ? "" : option.value);
                      setOpen(false);
                    }}
                  >
                    <CheckIcon
                      className={cn(
                        "mr-2 h-4 w-4",
                        value === option.value ? "opacity-100" : "opacity-0",
                      )}
                    />
                    {option.label}
                  </CommandItem>
                ))}
              </CommandList>
            </CommandGroup>
          </Command>
        </PopoverContent>
      </Popover>
    );
  },
);

Combobox.displayName = "Combobox";

export { Combobox, comboboxVariants };
