import * as React from "react";
import { Slot } from "@radix-ui/react-slot";

import { cn } from "@/lib/utils";

function CategoryCard({
  className,
  asChild = false,
  image,
  text,
  ...props
}: React.ComponentProps<"button"> & {
  asChild?: boolean;
  text: string;
  image: string;
}) {
  const Comp = asChild ? Slot : "button";

  return (
    <Comp
      data-slot="button"
      className={cn(
        "flex flex-col items-center gap-y-6 cursor-pointer outline-none group",
        className,
      )}
      {...props}
    >
      <div className="bg-primary p-1 rounded-full w-28 h-28 group-hover:bg-primary/70 transition-all">
        <img src={image} alt={text} className="w-full h-full object-contain" />
      </div>

      <span className="font-normal text-sm text-center">{text}</span>
    </Comp>
  );
}

export { CategoryCard };
