import * as React from "react";
import axios from "axios";

import {
  Command,
  CommandGroup,
  CommandItem,
  CommandList,
  CommandEmpty,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

import { IProduct } from "@/types/product";

import { ProductService } from "@/services/product";

type SearchBarProps = {
  onSelect?: (product: IProduct) => void;
  placeholder?: string;
  debounceMs?: number;
};

function ProductSearchBar({
  onSelect,
  placeholder = "Search products...",
  debounceMs = 400,
}: SearchBarProps) {
  const inputRef = React.useRef<HTMLInputElement | null>(null);
  const [open, setOpen] = React.useState(false);
  const [query, setQuery] = React.useState("");
  const [loading, setLoading] = React.useState(false);
  const [results, setResults] = React.useState<IProduct[]>([]);

  React.useEffect(() => {
    if (open) {
      const t = setTimeout(() => inputRef.current?.focus(), 0);
      return () => clearTimeout(t);
    }
  }, [open]);

  React.useEffect(() => {
    if (!query) {
      setResults([]);
      setLoading(false);
      return;
    }

    const source = axios.CancelToken.source();
    const timer = setTimeout(async () => {
      try {
        setLoading(() => true);
        const res = await ProductService.getProducts(
          { k: query },
          source.token,
        );
        setResults(res || []);
      } catch (err: any) {
        if (axios.isCancel(err)) {
          setOpen(() => false);
        } else {
          console.error("Search error:", err);
        }
      } finally {
        setLoading(() => false);
      }
    }, debounceMs);

    return () => {
      clearTimeout(timer);
      source.cancel("canceled by new request");
    };
  }, [query, debounceMs]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <div className="w-full max-w-md">
          <input
            ref={inputRef}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={placeholder}
            className="w-full rounded-lg border px-3 py-2 text-sm outline-none"
          />
        </div>
      </PopoverTrigger>

      <PopoverContent className="max-w-[360px] p-0">
        <Command>
          <CommandList>
            {loading ? (
              <CommandEmpty className="p-4">Loading...</CommandEmpty>
            ) : results.length === 0 ? (
              <CommandEmpty className="p-4">No results found.</CommandEmpty>
            ) : (
              <CommandGroup heading="Products">
                {results.map((product) => (
                  <CommandItem
                    key={product.id}
                    value={product.name}
                    className="cursor-pointer"
                    onSelect={() => {
                      onSelect?.(product);
                      setOpen(false);
                    }}
                  >
                    <div className="flex items-center gap-3">
                      {product.mainImage?.imageName ? (
                        <img
                          src={`http://localhost:5000/product/image/${product.mainImage.imageName}`}
                          alt={product.name}
                          className="w-10 h-10 object-cover rounded"
                        />
                      ) : (
                        <div className="w-10 h-10 bg-muted rounded" />
                      )}
                      <div className="min-w-0">
                        <p className="font-medium truncate">{product.name}</p>
                        <p className="text-xs text-muted-foreground">
                          ${product.price.toFixed(2)}
                        </p>
                      </div>
                    </div>
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}

export { ProductSearchBar };
