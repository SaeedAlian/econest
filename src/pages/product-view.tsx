import * as React from "react";
import { Star } from "lucide-react";
import { Link, useParams } from "react-router";

import { cn } from "@/lib/utils";

import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";

import {
  IProductAttribute,
  IProductAttributeOption,
  IProductCommentWithUser,
  IProductExtended,
  IProductImage,
  IProductVariantSelectedAttributeOption,
  IProductVariantWithAttributeSet,
} from "@/types/product";

import { ProductService } from "@/services/product";

function ProductView() {
  const params = useParams();

  const [product, setProduct] = React.useState<IProductExtended | null>(null);
  const [mainImage, setMainImage] = React.useState<IProductImage | null>(null);
  const [secondaryImages, setSecondaryImages] = React.useState<IProductImage[]>(
    [],
  );
  const [loading, setLoading] = React.useState<boolean>(true);
  const [error, setError] = React.useState<boolean>(false);

  const [comments, setComments] = React.useState<IProductCommentWithUser[]>([]);
  const [commentPage, setCommentPage] = React.useState<number>(1);
  const [commentPages, setCommentPages] = React.useState<number>(1);
  const [loadingComments, setLoadingComments] = React.useState<boolean>(true);

  const [selectedAttributes, setSelectedAttributes] = React.useState<
    IProductVariantSelectedAttributeOption[]
  >([]);

  const [selectedVariant, setSelectedVariant] =
    React.useState<IProductVariantWithAttributeSet | null>(null);

  const observerTargetRef = React.useRef<HTMLDivElement | null>(null);

  const onChangeSelectedAttribute = (
    attribute: IProductAttribute,
    selectedOption: IProductAttributeOption,
  ) => {
    let found = false;
    let attrSet = selectedAttributes.map((a) => {
      if (a.id == attribute.id) {
        found = true;
        return { ...a, selectedOption };
      } else {
        return a;
      }
    });

    if (!found) {
      attrSet = [...attrSet, { ...attribute, selectedOption }];
    }

    setSelectedAttributes(attrSet);
  };

  const onRemoveSelectedAttribute = (attribute: IProductAttribute) => {
    const attrSet = selectedAttributes.filter((a) => {
      return a.id !== attribute.id;
    });

    setSelectedAttributes(attrSet);
  };

  React.useEffect(() => {
    window.scrollTo(0, 0);

    async function run() {
      setLoading(true);
      setError(false);
      try {
        const productId = params.productId;
        if (productId == null) {
          setError(true);
          return;
        }
        const prod = await ProductService.getProductExtended(
          parseInt(productId),
        );
        setProduct(prod);
      } catch (err) {
        console.error(err);
        setError(true);
      } finally {
        setLoading(false);
      }
    }
    run();
  }, []);

  React.useEffect(() => {
    if (product == null) return;

    const mainImg = product.images.find((i) => i.isMain);
    const secondaryImgs = product.images.filter((i) => !i.isMain);
    if (mainImg) {
      setMainImage(mainImg);
    }
    setSecondaryImages(secondaryImgs);

    const selectedVariant = product.variants[0];
    setSelectedVariant(selectedVariant);
    setSelectedAttributes(selectedVariant.attributeSet);

    async function loadComments() {
      if (product == null) return;

      setLoadingComments(true);
      try {
        const pages = await ProductService.getProductCommentsPages(
          product.id,
          {},
        );
        setCommentPages(pages);
        const initialComments = await ProductService.getProductComments(
          product.id,
          { p: 1 },
        );
        setComments(initialComments);
      } catch (err) {
        console.error(err);
      } finally {
        setLoadingComments(false);
      }
    }
    loadComments();
  }, [product]);

  React.useEffect(() => {
    if (observerTargetRef.current == null) return;
    if (commentPage >= commentPages) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          if (observerTargetRef.current == null) return;
          observerTargetRef.current.hidden = true;
          console.log(commentPage);
          setCommentPage((prev) => prev + 1);
        }
      },
      { threshold: 1 },
    );

    observer.observe(observerTargetRef.current);

    return () => {
      if (observerTargetRef.current) {
        observer.unobserve(observerTargetRef.current);
      }
    };
  }, [comments, commentPages, commentPage, observerTargetRef]);

  React.useEffect(() => {
    if (product == null) return;
    if (commentPage === 1) return;
    if (commentPage > commentPages) return;

    async function fetchMore() {
      if (product == null) return;

      setLoadingComments(true);
      try {
        const more = await ProductService.getProductComments(
          product.id,
          { p: commentPage },
          undefined,
        );
        setComments((prev) => [...prev, ...more]);

        if (observerTargetRef.current == null) return;
        observerTargetRef.current.hidden = false;
      } catch (err) {
        console.error(err);
      } finally {
        setLoadingComments(false);
      }
    }
    fetchMore();
  }, [commentPage, product, commentPages]);

  return (
    <div className="container mx-auto px-6 py-8">
      {error ? (
        <>
          <div className="w-full flex flex-col items-center justify-center gap-6">
            <p className="text-3xl font-bold text-center">
              Something went wrong...
            </p>
            <Button asChild variant="destructive" size="lg">
              <Link to="/product">Back To List</Link>
            </Button>
          </div>
        </>
      ) : (
        <>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">
            {loading || product == null ? (
              <>
                <div>
                  <Skeleton className="rounded-xl w-full aspect-square bg-muted/80" />
                  <div className="flex gap-2 mt-4 flex-wrap justify-center">
                    {[1, 2, 3, 4].map((i) => (
                      <Skeleton
                        key={i}
                        className="rounded-md w-20 h-20 aspect-square bg-muted/80"
                      />
                    ))}
                  </div>
                </div>

                <div className="space-y-6">
                  <Skeleton className="rounded-md max-w-60 w-full h-8 aspect-square bg-muted/80" />
                  <div className="flex items-center gap-3">
                    <Skeleton className="rounded-md max-w-16 w-full h-8 aspect-square bg-muted/80" />
                  </div>

                  <div className="flex gap-2 flex-wrap">
                    <Skeleton className="rounded-md w-16 h-5 aspect-square bg-muted/80" />
                    <Skeleton className="rounded-md w-20 h-5 aspect-square bg-muted/80" />
                    <Skeleton className="rounded-md w-14 h-5 aspect-square bg-muted/80" />
                    <Skeleton className="rounded-md w-24 h-5 aspect-square bg-muted/80" />
                  </div>

                  <Skeleton className="rounded-md w-full h-5 aspect-square bg-muted/80" />
                  <Skeleton className="rounded-md w-full h-5 aspect-square bg-muted/80" />
                  <Skeleton className="rounded-md max-w-48 w-full h-5 aspect-square bg-muted/80" />

                  <Skeleton className="rounded-md w-full h-72 aspect-square bg-muted/80" />

                  <Skeleton className="rounded-md w-full h-64 aspect-square bg-muted/80" />
                  <Skeleton className="rounded-md w-full h-64 aspect-square bg-muted/80" />

                  <Skeleton className="rounded-md w-full h-10 aspect-square bg-muted/80" />
                </div>
              </>
            ) : (
              <>
                <div>
                  <img
                    src={
                      mainImage
                        ? `${ProductService.baseURL}/product/image/${mainImage.imageName}`
                        : `https://placehold.co/400x250?text=${product.name.replace(" ", "+")}`
                    }
                    alt="Product"
                    className="w-full aspect-square object-cover rounded-xl border p-3"
                  />
                  <div className="flex gap-2 mt-4 flex-wrap justify-center">
                    {secondaryImages.map((i) => (
                      <img
                        key={i.id}
                        src={`${ProductService.baseURL}/product/image/${i.imageName}`}
                        alt={`Thumbnail ${product.name}`}
                        className="w-20 h-20 object-cover rounded-md border cursor-pointer"
                      />
                    ))}
                  </div>
                </div>

                <div className="space-y-4 mb-3">
                  <h1 className="text-3xl font-bold">{product.name}</h1>
                  <div className="flex items-center gap-3">
                    <p className="text-2xl font-semibold">
                      $
                      {product.offer != null && product.offer.discount > 0
                        ? (
                            product.price *
                            (1 - product.offer.discount)
                          ).toLocaleString("en", { maximumFractionDigits: 2 })
                        : product.price.toLocaleString("en", {
                            maximumFractionDigits: 2,
                          })}
                    </p>
                    {product.offer != null ? (
                      <Badge className="bg-green-600 text-white">
                        {product.offer.discount * 100}% OFF
                      </Badge>
                    ) : null}
                  </div>

                  <div className="flex gap-2 flex-wrap mb-7">
                    {product.tags.map((tag) => (
                      <Badge key={tag.id} variant="outline">
                        {tag.name}
                      </Badge>
                    ))}
                  </div>

                  <p className="text-muted-foreground mb-8">
                    {product.description}
                  </p>

                  <Card>
                    <CardHeader>
                      <CardTitle>Specifications</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <table className="w-full text-sm">
                        <tbody>
                          {product.specs.map((s) => (
                            <tr key={s.id}>
                              <td className="py-1 font-medium">{s.label}</td>
                              <td className="py-1">{s.value}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </CardContent>
                  </Card>

                  <Card>
                    <CardHeader>
                      <CardTitle className="flex justify-between items-center">
                        <span>Variants</span>
                        {selectedVariant != null ? (
                          <span>Quantity: {selectedVariant.quantity}</span>
                        ) : null}
                      </CardTitle>
                    </CardHeader>
                    <CardContent className="flex flex-col gap-3">
                      {product.attributes.length > 0 ? (
                        product.attributes.map((v) => (
                          <div key={v.id}>
                            <p className="font-medium">{v.label}</p>
                            <div className="flex flex-wrap gap-2 mt-1">
                              {v.options.map((o) => (
                                <Button
                                  key={`${v.id}-${o.id}`}
                                  variant="outline"
                                  size="sm"
                                  onClick={() =>
                                    selectedAttributes.find(
                                      (a) =>
                                        a.id === v.id &&
                                        a.selectedOption.id === o.id,
                                    )
                                      ? onRemoveSelectedAttribute(v)
                                      : onChangeSelectedAttribute(v, o)
                                  }
                                  className={cn(
                                    selectedAttributes.find(
                                      (a) =>
                                        a.id === v.id &&
                                        a.selectedOption.id === o.id,
                                    )
                                      ? "bg-foreground text-background hover:bg-foreground/60 hover:text-background"
                                      : "",
                                  )}
                                >
                                  {o.value}
                                </Button>
                              ))}
                            </div>
                          </div>
                        ))
                      ) : (
                        <div>No variant found for this product</div>
                      )}
                    </CardContent>
                  </Card>

                  <Card>
                    <CardHeader>
                      <CardTitle>Store Information</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="font-medium">{product.store.name}</p>
                      <p className="text-sm text-muted-foreground">
                        {product.store.description}
                      </p>
                      <Button variant="link" className="mt-2 px-0">
                        Visit Store
                      </Button>
                    </CardContent>
                  </Card>

                  <Button size="lg" className="w-full">
                    Add to Cart
                  </Button>
                </div>
              </>
            )}
          </div>

          <Separator className="my-10" />

          <div>
            <h2 className="text-2xl font-semibold mb-6">Customer Reviews</h2>
            {loadingComments && comments.length === 0 ? (
              <div className="flex flex-col gap-4">
                <Skeleton className="rounded-md w-full h-52 bg-muted/80" />
                <Skeleton className="rounded-md w-full h-52 bg-muted/80" />
              </div>
            ) : (
              <div className="space-y-6">
                {comments.map((c) => (
                  <Card key={c.id}>
                    <CardContent>
                      <div className="mb-2">
                        {c.user.fullName ?? "Anonymous"}
                      </div>
                      <div className="flex items-center gap-2 mb-5">
                        {[...Array(5)].map((_, i) => (
                          <Star
                            key={i}
                            className={`w-4 h-4 ${
                              i < c.scoring
                                ? "fill-yellow-400 text-yellow-400"
                                : "text-muted-foreground"
                            }`}
                          />
                        ))}
                      </div>
                      <p>{c.comment}</p>
                      <p className="text-xs text-muted-foreground mt-2">
                        Posted on {c.createdAt.toLocaleDateString()}
                      </p>
                    </CardContent>
                  </Card>
                ))}
                <div ref={observerTargetRef}></div>
                {loadingComments && comments.length > 0 && (
                  <div className="text-center text-muted-foreground">
                    Loading more...
                  </div>
                )}
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}

export default ProductView;
