import * as React from "react";

type DivProps = React.HTMLAttributes<HTMLDivElement>;

export function Card({ className = "", ...props }: DivProps) {
  return (
    <div
      className={`rounded-xl border border-white/20 bg-white/10 p-4 shadow-lg backdrop-blur ${className}`}
      {...props}
    />
  );
}

export function CardHeader({ className = "", ...props }: DivProps) {
  return <div className={`mb-2 ${className}`} {...props} />;
}

export function CardTitle({ className = "", ...props }: DivProps) {
  return <h3 className={`text-base font-semibold ${className}`} {...props} />;
}

export function CardContent({ className = "", ...props }: DivProps) {
  return <div className={className} {...props} />;
}
