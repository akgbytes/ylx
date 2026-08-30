"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useMutation } from "@tanstack/react-query";
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

import { signUp } from "@/api/auth";
import { ZSignUpSchema, type SignUpValues } from "@/lib/validation/auth";

export function SignUpForm() {
  const router = useRouter();
  const signUpMutation = useMutation({
    mutationFn: signUp,
    onSuccess: (challenge, values) => {
      const params = new URLSearchParams({
        email: values.email,
        retryAt: challenge.retry_at,
      });
      router.push(`/verify-sign-up?${params.toString()}`);
    },
  });
  const form = useForm<SignUpValues>({
    resolver: zodResolver(ZSignUpSchema),
    defaultValues: { name: "", email: "", password: "" },
  });

  function submit(values: SignUpValues) {
    signUpMutation.reset();
    signUpMutation.mutate(values);
  }

  const nameError = form.formState.errors.name?.message;
  const emailError = form.formState.errors.email?.message;
  const passwordError = form.formState.errors.password?.message;

  return (
    <Card>
      <CardHeader className="text-center">
        <CardTitle className="text-xl">Create your account</CardTitle>
        <CardDescription>
          We will email you a code to verify your account
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={form.handleSubmit(submit)} noValidate>
          <FieldGroup>
            <Field data-invalid={!!nameError}>
              <FieldLabel htmlFor="name">Name</FieldLabel>
              <Input
                {...form.register("name")}
                id="name"
                autoComplete="name"
                placeholder="Your name"
                aria-invalid={!!nameError}
                required
              />
              {nameError && <FieldError>{nameError}</FieldError>}
            </Field>

            <Field data-invalid={!!emailError}>
              <FieldLabel htmlFor="email">Email</FieldLabel>
              <Input
                {...form.register("email")}
                id="email"
                type="email"
                autoComplete="email"
                spellCheck={false}
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
                autoComplete="new-password"
                placeholder="At least 6 characters"
                aria-invalid={!!passwordError}
                required
              />
              {passwordError && <FieldError>{passwordError}</FieldError>}
            </Field>

            <Field>
              <Button type="submit" disabled={signUpMutation.isPending}>
                {signUpMutation.isPending ? "Sending code…" : "Continue"}
              </Button>
              {signUpMutation.isError && (
                <FieldError>{signUpMutation.error.message}</FieldError>
              )}
              <FieldDescription className="text-center">
                Already have an account? <Link href="/sign-in">Sign in</Link>
              </FieldDescription>
            </Field>
          </FieldGroup>
        </form>
      </CardContent>
    </Card>
  );
}
