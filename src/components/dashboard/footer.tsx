export type FooterProps = {};

function Footer({}: FooterProps) {
  return (
    <footer className="mt-auto w-full flex flex-col items-center">
      <p className="text-foreground/30 py-3 text-xs">
        Copyright © 2025 EcoNest
      </p>
    </footer>
  );
}

export { Footer };
