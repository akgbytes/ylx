import { z } from "zod";

const ZEmailSchema = z
  .string()
  .trim()
  .toLowerCase()
  .pipe(z.email("Enter a valid email address"));

export const ZSignInSchema = z.object({
  email: ZEmailSchema,
  password: z.string().min(1, "Password is required"),
});

export type SignInValues = z.infer<typeof ZSignInSchema>;
