import { toast } from "sonner";

export function useToast() {
  return {
    toast: (props: {
      title?: string;
      description?: string;
      variant?: "default" | "destructive";
    }) => {
      const { title, description, variant = "default" } = props;

      if (variant === "destructive") {
        toast.error(title, {
          description,
        });
      } else {
        toast(title, {
          description,
        });
      }
    },
  };
} 