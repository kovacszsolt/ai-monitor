import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded font-label-md text-label-md transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50 disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90 shadow-sm",
        /** Segmented / filter chips on dark surfaces — high contrast when selected */
        toggle:
          "bg-surface-bright text-on-surface shadow-sm ring-1 ring-primary/35 hover:bg-surface-container-highest",
        secondary:
          "bg-surface-container-high text-on-surface hover:bg-surface-bright shadow-sm",
        outline:
          "border border-outline-variant bg-transparent text-on-surface hover:bg-surface-container",
        ghost: "text-on-surface-variant hover:bg-surface-container hover:text-on-surface",
        destructive:
          "bg-error-container text-on-error-container hover:opacity-90 shadow-sm",
      },
      size: {
        default: "h-9 px-4 py-2",
        sm: "h-8 px-3 text-label-sm",
        lg: "h-10 px-space-lg",
        icon: "h-7 w-7",
      },
    },
    defaultVariants: { variant: "default", size: "default" },
  },
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

export function Button({ className, variant, size, asChild = false, ...props }: ButtonProps) {
  const Comp = asChild ? Slot : "button";
  return <Comp className={cn(buttonVariants({ variant, size }), className)} {...props} />;
}
