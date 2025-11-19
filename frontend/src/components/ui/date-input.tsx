import * as React from "react";

import { cn } from "../../lib/utils";
import { Input } from "./input";

export interface DateInputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {}

const DateInput = React.forwardRef<HTMLInputElement, DateInputProps>(
  ({ className, ...props }, ref) => {
    return (
      <Input
        type="date"
        className={cn(
          "[color-scheme:light] dark:[color-scheme:light]",
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
DateInput.displayName = "DateInput";

export { DateInput };
