import { useState, useEffect, Fragment } from "react";
import { Link } from "react-router";
import { FaShoppingCart } from "react-icons/fa";
import { FaWallet } from "react-icons/fa6";

import { useIsMobile } from "@/hooks/use-mobile";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { SidebarTrigger } from "@/components/ui/sidebar";
import {
  Breadcrumb,
  BreadcrumbEllipsis,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Button } from "@/components/ui/button";

export type BreadcrumbLink = {
  text: string;
  href: string;
};

export type NavbarProps = {
  cartHref?: string;
  walletHref?: string;
  inCartCount?: number;
  breadcrumbLinks?: BreadcrumbLink[];
};

function Navbar({
  breadcrumbLinks,
  cartHref,
  inCartCount,
  walletHref,
}: NavbarProps) {
  const isMobile = useIsMobile();
  const [showEllipsis, setShowEllipsis] = useState<boolean>(false);
  const [visibleBreadcrumbLinks, setVisibleBreadcrumbLinks] = useState<
    BreadcrumbLink[]
  >([]);

  useEffect(() => {
    if (breadcrumbLinks == null || breadcrumbLinks?.length == 0) return;

    if (breadcrumbLinks && breadcrumbLinks.length > 3) {
      const items = [
        breadcrumbLinks[0],
        breadcrumbLinks[breadcrumbLinks.length - 2],
        breadcrumbLinks[breadcrumbLinks.length - 1],
      ];
      setVisibleBreadcrumbLinks(() => items);
      setShowEllipsis(() => true);
    } else {
      setVisibleBreadcrumbLinks(() => breadcrumbLinks);
      setShowEllipsis(() => false);
    }
  }, [breadcrumbLinks]);

  return (
    <>
      <nav className="flex items-center flex-row gap-2 w-full pr-3 pl-2 py-2">
        <SidebarTrigger />
        <Separator orientation="vertical" />
        {!isMobile ? (
          <Breadcrumb className="ml-2">
            <BreadcrumbList>
              {visibleBreadcrumbLinks.map((l, i) => (
                <Fragment key={l.href}>
                  <BreadcrumbItem>
                    <BreadcrumbLink href={l.href ?? ""}>
                      {l.text}
                    </BreadcrumbLink>
                  </BreadcrumbItem>
                  {i != visibleBreadcrumbLinks.length - 1 ? (
                    <BreadcrumbSeparator />
                  ) : null}
                  {showEllipsis && i == 0 ? (
                    <>
                      <BreadcrumbEllipsis />
                      <BreadcrumbSeparator />
                    </>
                  ) : null}
                </Fragment>
              ))}
            </BreadcrumbList>
          </Breadcrumb>
        ) : null}
        <Button
          variant="ghost"
          size="icon"
          className="rounded-full relative ml-auto"
          asChild
        >
          <Link to={cartHref ?? ""}>
            {inCartCount != null && inCartCount > 0 ? (
              <>
                <Badge
                  className="absolute top-0 right-0 text-[0.6rem] flex items-center justify-center w-3 h-3 rounded-full aspect-square p-2"
                  variant="destructive"
                >
                  {inCartCount}
                </Badge>
              </>
            ) : null}
            <FaShoppingCart />
          </Link>
        </Button>
        <Button variant="ghost" size="icon" className="rounded-full">
          <Link to={walletHref ?? ""}>
            <FaWallet />
          </Link>
        </Button>
      </nav>
    </>
  );
}

export { Navbar };
