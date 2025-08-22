import * as React from "react";
import { FaRegStar } from "react-icons/fa";
import { BsShop } from "react-icons/bs";
import {
  MdOutlineAddShoppingCart,
  MdOutlineShoppingCartCheckout,
  MdLocalOffer,
} from "react-icons/md";

import { cn } from "@/lib/utils";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";

function ProductCard({
  className,
  productId,
  image,
  storeName,
  title,
  categoryName,
  isLowStock,
  isOutStock,
  score,
  price,
  discount,
  inCart,
  onAddToCart,
  ...props
}: React.ComponentProps<"div"> & {
  productId: number;
  image: string;
  storeName: string;
  title: string;
  categoryName: string;
  isLowStock?: boolean;
  isOutStock?: boolean;
  score: number;
  price: number;
  discount?: number;
  inCart?: boolean;
  onAddToCart: (id: number) => void;
}) {
  return (
    <Card className="min-w-[280px] py-4" {...props}>
      <CardHeader
        className="px-4 relative group cursor-pointer !flex flex-col items-center"
        onClick={() => onAddToCart(productId)}
      >
        <img src={image} alt={title} className="w-40 h-40 object-contain" />
        <div className="absolute left-4 bottom-0 flex items-center">
          <div
            className={cn(
              "flex items-center rounded-full px-3 py-1 transition-all duration-300 ease-in-out",
              inCart
                ? "bg-muted text-muted-foreground"
                : "bg-primary text-primary-foreground",
            )}
          >
            <div className="flex items-center justify-center rounded-full">
              {inCart ? (
                <MdOutlineShoppingCartCheckout className="!w-4 !h-4" />
              ) : (
                <MdOutlineAddShoppingCart className="!w-4 !h-4" />
              )}
            </div>
            <span className="ml-0 opacity-0 max-w-0 group-hover:max-w-xs group-hover:ml-2 group-hover:opacity-100 transition-all duration-300 whitespace-nowrap overflow-hidden text-xs">
              {inCart ? "In Cart" : "Add to Cart"}
            </span>
          </div>
        </div>
        <div className="absolute right-4 bottom-0 flex flex-col gap-y-1 items-end">
          {discount != null && discount > 0 ? (
            <div className="flex items-center bg-accent text-accent-foreground text-xs font-bold px-2 py-1 rounded-full">
              <MdLocalOffer />
            </div>
          ) : null}
        </div>
      </CardHeader>
      <CardContent className="px-4 flex flex-col gap-1">
        <div className="flex items-center gap-1">
          <BsShop className="w-3 h-3 text-accent" />
          <p className="text-xs mt-[0.18rem]">{storeName}</p>
        </div>
        <p className="text-sm font-bold">{title}</p>
        <p className="text-xs font-normal text-card-foreground/50">
          {categoryName}
        </p>
      </CardContent>
      <CardFooter className="px-4 flex items-end gap-1 mt-auto">
        <div>
          {isLowStock || isOutStock ? (
            <div className="text-destructive text-xs mb-1 font-bold">
              {isOutStock ? "Out of stock!" : "Low on stock!"}
            </div>
          ) : null}
          <div className="flex items-center gap-1">
            <FaRegStar className="text-yellow-600 w-4 h-4" />
            <span className="leading-3 mt-[0.18rem] text-sm">{score}</span>
          </div>
        </div>

        <div className="ml-auto flex flex-col items-end">
          {discount != null && discount > 0 ? (
            <p className="text-sm font-extrabold line-through text-card-foreground/20 mb-1">
              ${price.toLocaleString("en")}
            </p>
          ) : null}
          <p className="text-sm font-extrabold leading-3">
            $
            {discount != null && discount > 0
              ? (price * discount).toLocaleString("en")
              : price.toLocaleString("en")}
          </p>
        </div>
      </CardFooter>
    </Card>
  );
}

export { ProductCard };
