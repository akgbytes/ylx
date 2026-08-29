"use client";

import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Button } from "@ylx/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@ylx/ui/components/card";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@ylx/ui/components/field";
import { Input } from "@ylx/ui/components/input";

import { ZSignInSchema, type SignInValues } from "@/lib/validation/auth";

export function SignInForm() {
  const form = useForm<SignInValues>({
    resolver: zodResolver(ZSignInSchema),
    defaultValues: { email: "", password: "" },
  });
  const emailError = form.formState.errors.email?.message;
  const passwordError = form.formState.errors.password?.message;

  function submit(values: SignInValues) {
    void values;
  }

  return (
    <Card>
      <CardHeader className="text-center">
        <CardTitle className="text-xl">Welcome back</CardTitle>
        <CardDescription>Sign in with your email and password</CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={form.handleSubmit(submit)} noValidate>
          <FieldGroup>
            <Field data-invalid={!!emailError}>
              <FieldLabel htmlFor="email">Email</FieldLabel>
              <Input
                {...form.register("email")}
                id="email"
                type="email"
                autoComplete="email"
                placeholder="you@example.com"
                aria-invalid={!!emailError}
                required
              />
              {emailError && <FieldError>{emailError}</FieldError>}
            </Field>

            <Field data-invalid={!!passwordError}>
              <FieldLabel htmlFor="password">Password</FieldLabel>
              <Input
                {...form.register("password")}
                id="password"
                type="password"
                autoComplete="current-password"
                placeholder="Enter your password"
                aria-invalid={!!passwordError}
                required
              />
              {passwordError && <FieldError>{passwordError}</FieldError>}
            </Field>

            <Field>
              <Button type="submit">Sign in</Button>
              <FieldDescription className="text-center">
                New to YLX? <Link href="/sign-up">Create an account</Link>
              </FieldDescription>
            </Field>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
}
