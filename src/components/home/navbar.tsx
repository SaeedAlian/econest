import { Link, useNavigate } from "react-router";
import { GiHamburgerMenu } from "react-icons/gi";

import { cn } from "@/lib/utils";

import { useIsMobile } from "@/hooks/use-mobile";

import { Button } from "@/components/ui/button";
import { ProductSearchBar } from "@/components/ui/product-searchbar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import logo from "@/assets/logo-dark.svg";

export type NavbarLink = {
  text: string;
  href: string;
  target?: string;
  bordered?: boolean;
};

export type NavbarProps = {
  links: NavbarLink[];
};

function Navbar({ links }: NavbarProps) {
  const nav = useNavigate();
  const isMobile = useIsMobile();

  return (
    <nav className="flex items-center flex-row gap-2 w-full px-5 py-2 sticky top-0 bg-background/80 backdrop-blur z-20">
      <Link
        to="/"
        className={cn(
          "flex items-center justify-center m-[-20px] mr-auto",
          isMobile ? "m-[-25px]" : "",
        )}
      >
        <img
          src={logo}
          alt="EcoNest"
          className={cn("w-24 h-24", isMobile ? "w-20 h-20" : "")}
        />
      </Link>
      <div className="w-full max-w-[390px] mr-7 max-[1120px]:hidden">
        <ProductSearchBar
          onSelect={(p) => {
            nav(`/product/${p.id}`);
          }}
        />
      </div>
      {isMobile ? (
        <>
          <DropdownMenu>
            <DropdownMenuTrigger className="ml-auto">
              <GiHamburgerMenu />
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              {links.map((l, i) => (
                <DropdownMenuItem key={`${l.text}-${i}`}>
                  <Link
                    to={l.href}
                    target={l.target ?? "_self"}
                    className="w-full h-full"
                  >
                    {l.text}
                  </Link>
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </>
      ) : (
        <div className="flex flex-row items-center flex-wrap justify-center gap-5">
          {links.map((l, i) => (
            <Button
              key={`${l}-${i}`}
              asChild
              variant={l.bordered ? "outline" : "link"}
              className="text-foreground text-xs"
              size="sm"
            >
              <Link to={l.href} target={l.target ?? "_self"}>
                {l.text}
              </Link>
            </Button>
          ))}
        </div>
      )}
    </nav>
  );
}

export { Navbar };
