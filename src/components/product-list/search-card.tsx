import * as React from "react";

import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { MultipleSelector } from "@/components/ui/multiselect";
import { Checkbox } from "@/components/ui/checkbox";

import { IProductTag } from "@/types/product";

import { ProductService } from "@/services/product";

function SearchCard({
  className,
  fullWidth,
  minPrice,
  maxPrice,
  showLowStock,
  showOffers,
  showPopular,
  onSetTags,
  onChangeMinPrice,
  onChangeMaxPrice,
  onChangeShowLowStock,
  onChangeShowOffers,
  onChangeShowPopular,
  ...props
}: React.ComponentProps<"div"> & {
  fullWidth?: boolean;
  minPrice: number | undefined;
  maxPrice: number | undefined;
  showLowStock: boolean;
  showOffers: boolean;
  showPopular: boolean;
  onSetTags: (values: IProductTag[]) => void;
  onChangeMinPrice: (value: string) => void;
  onChangeMaxPrice: (value: string) => void;
  onChangeShowLowStock: (value: boolean) => void;
  onChangeShowOffers: (value: boolean) => void;
  onChangeShowPopular: (value: boolean) => void;
}) {
  const [tagsResult, setTagsResult] = React.useState<IProductTag[]>([]);
  const [tagsSearchVal, setTagsSearchVal] = React.useState<string>("");

  const getTags = async () => {
    const tags = await ProductService.getProductTags({
      lim: 10,
      offst: 0,
      name: tagsSearchVal,
    });

    setTagsResult(() => tags);
  };

  const handleChangeTags = (values: string[]) => {
    const ids = values.map((v) => parseInt(v));
    const notFoundIds: number[] = [];
    const selectedTags: IProductTag[] = [];

    ids.forEach((i) => {
      const tag = tagsResult.find((t) => i === t.id);
      if (tag) {
        selectedTags.push(tag);
      } else {
        notFoundIds.push(i);
      }
    });

    notFoundIds.forEach(async (i) => {
      try {
        const tag = await ProductService.getProductTag(i);
        selectedTags.push(tag);
      } catch {}
    });

    onSetTags(selectedTags);
  };

  React.useEffect(() => {
    const timeout = setTimeout(() => {
      getTags();
    }, 400);

    return () => {
      clearTimeout(timeout);
    };
  }, [tagsSearchVal]);

  return (
    <Card className="py-4 px-1 w-full" {...props}>
      <CardHeader>
        <CardTitle className="text-lg">Filters</CardTitle>
        <CardDescription className="text-sm text-muted-foreground">
          Refine your search
        </CardDescription>
      </CardHeader>
      <CardContent className="w-full">
        <div className="space-y-4 w-full">
          <div className="w-full">
            <div className="text-sm mb-2 font-medium">Tags</div>
            <div className="flex flex-wrap gap-2 w-full">
              <MultipleSelector
                variant="outline"
                fullWidth={fullWidth}
                onChange={handleChangeTags}
                search={tagsSearchVal}
                onSearchChange={(v) => setTagsSearchVal(v)}
                items={tagsResult.map((t) => ({
                  label: t.name,
                  value: t.id.toString(),
                }))}
              />
            </div>
          </div>

          <div>
            <div className="text-sm mb-2 font-medium">Price</div>
            <div className="flex gap-2">
              <Input
                type="number"
                value={minPrice}
                onChange={(e) => onChangeMinPrice(e.target.value)}
                placeholder="min"
                className="w-1/2"
              />
              <Input
                value={maxPrice}
                onChange={(e) => onChangeMaxPrice(e.target.value)}
                placeholder="max"
                className="w-1/2"
              />
            </div>
          </div>

          <div>
            <div className="flex flex-col gap-2">
              <div className="flex items-center gap-2">
                <Checkbox
                  onCheckedChange={onChangeShowLowStock}
                  checked={showLowStock}
                  className="bg-muted"
                />
                <span>Low stock</span>
              </div>

              <div className="flex items-center gap-2">
                <Checkbox
                  onCheckedChange={onChangeShowOffers}
                  checked={showOffers}
                  className="bg-muted"
                />
                <span>Offers</span>
              </div>

              <div className="flex items-center gap-2">
                <Checkbox
                  onCheckedChange={onChangeShowPopular}
                  checked={showPopular}
                  className="bg-muted"
                />
                <span>Popular</span>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export { SearchCard };
