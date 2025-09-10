import * as React from "react";
import { Filter } from "lucide-react";
import { Link } from "react-router";

import { SearchCard } from "@/components/product-list/search-card";
import { ProductCard } from "@/components/ui/product-card";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Skeleton } from "@/components/ui/skeleton";

import { IProduct, IProductSearchQuery, IProductTag } from "@/types/product";

import { ProductService } from "@/services/product";

function ProductListPage() {
  const [page, setPage] = React.useState<number>(1);
  const [pages, setPages] = React.useState<number>(1);
  const [prods, setProds] = React.useState<IProduct[]>([]);
  const [loading, setLoading] = React.useState<boolean>(true);
  const [error, setError] = React.useState<boolean>(false);
  const [searchQuery, setSearchQuery] = React.useState<{
    keyword: string;
    minPrice?: number;
    maxPrice?: number;
    tags: IProductTag[];
    showPopular: boolean;
    showOffers: boolean;
    showLowStock: boolean;
  }>({
    keyword: "",
    showLowStock: false,
    showOffers: false,
    showPopular: false,
    tags: [],
    maxPrice: undefined,
    minPrice: undefined,
  });

  const buildQuery = React.useCallback((): IProductSearchQuery => {
    return {
      k: searchQuery.keyword ?? undefined,
      avgscr: searchQuery.showPopular ? 3.5 : undefined,
      minq: searchQuery.showLowStock ? 1 : undefined,
      maxq: searchQuery.showLowStock ? 4 : undefined,
      offr: searchQuery.showOffers ? true : undefined,
      pmt: searchQuery.minPrice ?? undefined,
      plt: searchQuery.maxPrice ?? undefined,
      tags:
        searchQuery.tags.length > 0
          ? searchQuery.tags.map((t) => t.id).join(",")
          : undefined,
      p: page,
    };
  }, [searchQuery, page]);

  const handleAddToCart = (id: number) => {
    console.log("Added to cart:", id);
  };

  const onSetTags = (tags: IProductTag[]) => {
    setSearchQuery((prev) => ({ ...prev, tags }));
  };

  const onChangeKeyword = (value: string) => {
    setSearchQuery((prev) => ({ ...prev, keyword: value }));
  };

  const onChangeMinPrice = (value: string) => {
    if (value == null || value === "") {
      return setSearchQuery((prev) => ({ ...prev, minPrice: undefined }));
    }

    const val = parseFloat(value);
    setSearchQuery((prev) => ({ ...prev, minPrice: val }));
  };

  const onChangeMaxPrice = (value: string) => {
    if (value == null || value === "") {
      return setSearchQuery((prev) => ({ ...prev, maxPrice: undefined }));
    }

    const val = parseFloat(value);
    setSearchQuery((prev) => ({ ...prev, maxPrice: val }));
  };

  const onChangeShowPopular = (value: boolean) => {
    setSearchQuery((prev) => ({ ...prev, showPopular: value }));
  };

  const onChangeShowLowStock = (value: boolean) => {
    setSearchQuery((prev) => ({ ...prev, showLowStock: value }));
  };

  const onChangeShowOffers = (value: boolean) => {
    setSearchQuery((prev) => ({ ...prev, showOffers: value }));
  };

  const handlePageChange = (newPage: number) => {
    if (newPage >= 1 && newPage <= pages) {
      setPage(newPage);
    }
  };

  const updateProds = async () => {
    setLoading(() => true);
    setError(() => false);
    try {
      const query = buildQuery();
      const [res, totalPages] = await Promise.all([
        ProductService.getProducts(query),
        ProductService.getProductPages(query),
      ]);
      setProds(() => res);
      setPages(() => totalPages);
    } catch (err) {
      console.error(err);
      setError(() => true);
    } finally {
      setLoading(() => false);
    }
  };

  React.useEffect(() => {
    const timeout = setTimeout(() => {
      updateProds();
    }, 400);
    return () => {
      clearTimeout(timeout);
    };
  }, [searchQuery, page]);

  return (
    <section className="w-full px-5 py-3">
      {error ? (
        <>
          <div className="w-full flex flex-col items-center justify-center gap-6">
            <p className="text-3xl font-bold text-center">
              Something went wrong...
            </p>
            <Button asChild variant="destructive" size="lg">
              <Link to="/">Back To Home</Link>
            </Button>
          </div>
        </>
      ) : (
        <>
          <div className="flex flex-col md:flex-row gap-6">
            <div className="flex-1">
              <div className="mb-6 flex flex-col items-center gap-3 @container">
                <div className="w-full flex items-center justify-between gap-3">
                  <div>
                    <h2 className="text-xl font-semibold max-sm:text-lg">
                      Products
                    </h2>
                    <p className="text-sm text-muted-foreground max-sm:text-xs">
                      Handpicked items just for you
                    </p>
                  </div>

                  <div className="flex gap-3 flex-row items-center">
                    <div className="flex items-center gap-2">
                      <Input
                        className="bg-card hidden @lg:inline-flex"
                        placeholder="Search products..."
                        value={searchQuery.keyword}
                        onChange={(e) => onChangeKeyword(e.target.value)}
                      />

                      <Dialog>
                        <DialogTrigger
                          asChild
                          className="hidden max-lg:inline-flex"
                        >
                          <Button variant="ghost" size="icon">
                            <Filter />
                          </Button>
                        </DialogTrigger>
                        <DialogContent className="p-0">
                          <SearchCard
                            fullWidth
                            minPrice={searchQuery.minPrice}
                            maxPrice={searchQuery.maxPrice}
                            showLowStock={searchQuery.showLowStock}
                            showOffers={searchQuery.showOffers}
                            showPopular={searchQuery.showPopular}
                            onSetTags={onSetTags}
                            onChangeMinPrice={onChangeMinPrice}
                            onChangeMaxPrice={onChangeMaxPrice}
                            onChangeShowLowStock={onChangeShowLowStock}
                            onChangeShowOffers={onChangeShowOffers}
                            onChangeShowPopular={onChangeShowPopular}
                          />
                        </DialogContent>
                      </Dialog>

                      <div className="text-sm text-muted-foreground whitespace-nowrap max-sm:hidden">
                        Showing {prods.length} results
                      </div>
                    </div>
                  </div>
                </div>

                <Input
                  className="bg-card inline-flex @lg:hidden"
                  placeholder="Search products..."
                />
              </div>

              <div className="@container w-full">
                <div className="grid grid-cols-1 @[565px]:grid-cols-2 @[1000px]:grid-cols-3 @[1290px]:grid-cols-4 gap-6">
                  {loading ? (
                    <>
                      <Skeleton className="min-w-[240px] w-full h-full rounded-md aspect-[9/10] bg-muted/80" />
                      <Skeleton className="min-w-[240px] w-full h-full rounded-md aspect-[9/10] bg-muted/80" />
                      <Skeleton className="min-w-[240px] w-full h-full rounded-md aspect-[9/10] bg-muted/80" />
                      <Skeleton className="min-w-[240px] w-full h-full rounded-md aspect-[9/10] bg-muted/80" />
                      <Skeleton className="min-w-[240px] w-full h-full rounded-md aspect-[9/10] bg-muted/80" />
                      <Skeleton className="min-w-[240px] w-full h-full rounded-md aspect-[9/10] bg-muted/80" />
                    </>
                  ) : (
                    <>
                      {prods.map((p) => (
                        <ProductCard
                          key={p.id}
                          productId={p.id}
                          image={
                            p.mainImage
                              ? `http://localhost:5000/product/image/${p.mainImage.imageName}`
                              : ""
                          }
                          storeName={p.store.name}
                          title={p.name}
                          categoryName={p.subcategory.name}
                          price={p.price}
                          score={p.averageScore}
                          onAddToCart={handleAddToCart}
                          isLowStock={p.totalQuantity < 5}
                          isOutStock={p.totalQuantity === 0}
                          discount={p.offer != null ? p.offer.discount : 0}
                        />
                      ))}
                    </>
                  )}
                </div>
              </div>

              <div className="mt-8 flex items-center justify-center">
                <nav
                  className="inline-flex rounded-md gap-3"
                  aria-label="Pagination"
                >
                  <Button
                    variant="ghost"
                    disabled={page === 1}
                    onClick={() => handlePageChange(page - 1)}
                  >
                    Previous
                  </Button>

                  {Array.from({ length: pages }).map((_, i) => (
                    <Button
                      key={i + 1}
                      variant={page === i + 1 ? "default" : "ghost"}
                      onClick={() => handlePageChange(i + 1)}
                    >
                      {i + 1}
                    </Button>
                  ))}

                  <Button
                    variant="ghost"
                    disabled={page === pages}
                    onClick={() => handlePageChange(page + 1)}
                  >
                    Next
                  </Button>
                </nav>
              </div>
            </div>

            <aside className="w-63 flex-shrink-0 max-lg:hidden">
              <SearchCard
                minPrice={searchQuery.minPrice}
                maxPrice={searchQuery.maxPrice}
                showLowStock={searchQuery.showLowStock}
                showOffers={searchQuery.showOffers}
                showPopular={searchQuery.showPopular}
                onSetTags={onSetTags}
                onChangeMinPrice={onChangeMinPrice}
                onChangeMaxPrice={onChangeMaxPrice}
                onChangeShowLowStock={onChangeShowLowStock}
                onChangeShowOffers={onChangeShowOffers}
                onChangeShowPopular={onChangeShowPopular}
              />
            </aside>
          </div>
        </>
      )}
    </section>
  );
}

export default ProductListPage;
