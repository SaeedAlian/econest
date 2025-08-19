import { Link } from "react-router";
import { GiHamburgerMenu } from "react-icons/gi";

import { cn } from "@/lib/utils";
import { useIsMobile } from "@/hooks/use-mobile";
import { Button } from "@/components/ui/button";
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
  const isMobile = useIsMobile();

  return (
    <nav className="flex items-center flex-row gap-2 w-full justify-between px-5 py-2">
      <Link to="/">
        <img
          src={logo}
          alt="EcoNest"
          className={cn(
            "w-24 h-24 m-[-20px]",
            isMobile ? "w-20 h-20 m-[-25px]" : "",
          )}
        />
      </Link>
      {isMobile ? (
        <>
          <DropdownMenu>
            <DropdownMenuTrigger>
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
