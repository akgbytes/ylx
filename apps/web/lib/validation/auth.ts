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

export const ZSignUpSchema = z.object({
  name: z.string().trim().min(1, "Name is required"),
  email: ZEmailSchema,
  password: z.string().min(6, "Password must contain at least 6 characters"),
});

export const ZVerificationCodeSchema = z.object({
  otp: z
    .string()
    .trim()
    .regex(/^\d{6}$/, "Enter the 6-digit verification code"),
});

export type SignInValues = z.infer<typeof ZSignInSchema>;
export type SignUpValues = z.infer<typeof ZSignUpSchema>;
export type VerificationCodeValues = z.infer<typeof ZVerificationCodeSchema>;
