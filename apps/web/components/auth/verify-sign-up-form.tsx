"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Button, buttonVariants } from "@ylx/ui/components/button";
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

import { authQueryKeys, resendSignUp, verifySignUp } from "@/api/auth";
import { ApiError } from "@/lib/api-client";
import {
  ZVerificationCodeSchema,
  type VerificationCodeValues,
} from "@/lib/validation/auth";

export function VerifySignUpForm() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const queryClient = useQueryClient();
  const email = searchParams.get("email") ?? "";
  const [retryAt, setRetryAt] = useState(searchParams.get("retryAt"));
  const secondsRemaining = useCooldown(retryAt);
  const verifyMutation = useMutation({
    mutationFn: verifySignUp,
    onSuccess: (user) => {
      queryClient.setQueryData(authQueryKeys.me(), user);
      router.replace("/");
    },
  });
  const resendMutation = useMutation({
    mutationFn: resendSignUp,
    onSuccess: (challenge) => {
      setRetryAt(challenge.retry_at);
      router.replace(verificationPath(email, challenge.retry_at));
      form.reset();
    },
    onError: (error) => {
      const nextRetryAt = getRetryAt(error);
      if (nextRetryAt) {
        setRetryAt(nextRetryAt);
        router.replace(verificationPath(email, nextRetryAt));
      }
    },
  });
  const form = useForm<VerificationCodeValues>({
    resolver: zodResolver(ZVerificationCodeSchema),
    defaultValues: { otp: "" },
  });
  const otpError = form.formState.errors.otp?.message;

  if (!email) {
    return (
      <Card>
        <CardHeader className="text-center">
          <CardTitle className="text-xl">Verification link expired</CardTitle>
          <CardDescription>
            Start signup again to receive a new code.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Link href="/sign-up" className={`${buttonVariants()} w-full`}>
            Start again
          </Link>
        </CardContent>
      </Card>
    );
  }

  function submit(values: VerificationCodeValues) {
    verifyMutation.reset();
    verifyMutation.mutate({ email, otp: values.otp });
  }

  return (
    <Card>
      <CardHeader className="text-center">
        <CardTitle className="text-xl">Verify your email</CardTitle>
        <CardDescription>
          Enter the 6-digit code sent to {email}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={form.handleSubmit(submit)} noValidate>
          <FieldGroup>
            <Field data-invalid={!!otpError}>
              <FieldLabel htmlFor="otp">Verification code</FieldLabel>
              <Input
                {...form.register("otp", {
                  onChange: (event) => {
                    form.setValue(
                      "otp",
                      event.target.value.replace(/\D/g, "").slice(0, 6),
                      { shouldDirty: true, shouldValidate: true }
                    );
                  },
                })}
                id="otp"
                inputMode="numeric"
                autoComplete="one-time-code"
                spellCheck={false}
                pattern="[0-9]*"
                placeholder="000000"
                maxLength={6}
                aria-invalid={!!otpError}
                required
              />
              {otpError && <FieldError>{otpError}</FieldError>}
            </Field>

            <Field>
              <Button type="submit" disabled={verifyMutation.isPending}>
                {verifyMutation.isPending ? "Verifying…" : "Verify account"}
              </Button>
              {verifyMutation.isError && (
                <FieldError>{verifyMutation.error.message}</FieldError>
              )}
              <FieldDescription className="text-center">
                Didn&apos;t receive the code?{" "}
                <button
                  type="button"
                  className="underline underline-offset-4 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:no-underline disabled:opacity-60"
                  disabled={secondsRemaining > 0 || resendMutation.isPending}
                  onClick={() => {
                    resendMutation.reset();
                    resendMutation.mutate(email);
                  }}
                >
                  {resendMutation.isPending
                    ? "Sending…"
                    : secondsRemaining > 0
                      ? `Resend in ${secondsRemaining}s`
                      : "Resend code"}
                </button>
              </FieldDescription>
              {resendMutation.isError && (
                <FieldError>{resendMutation.error.message}</FieldError>
              )}
            </Field>
          </FieldGroup>
        </form>
      </CardContent>
      <CardContent className="pt-0 text-center">
        <FieldDescription>
          Wrong email? <Link href="/sign-up">Start again</Link>
        </FieldDescription>
      </CardContent>
    </Card>
  );
}

function useCooldown(retryAt: string | null) {
  const [now, setNow] = useState(() => Date.now());
  const retryTime = retryAt ? Date.parse(retryAt) : 0;

  useEffect(() => {
    const remaining = retryTime - Date.now();
    if (!Number.isFinite(remaining) || remaining <= 0) return;

    const timer = window.setTimeout(
      () => setNow(Date.now()),
      Math.min(remaining, 1_000)
    );
    return () => window.clearTimeout(timer);
  }, [now, retryTime]);

  if (!Number.isFinite(retryTime)) return 0;
  return Math.max(Math.ceil((retryTime - now) / 1_000), 0);
}

function verificationPath(email: string, retryAt: string) {
  const params = new URLSearchParams({ email, retryAt });
  return `/verify-sign-up?${params.toString()}`;
}

function getRetryAt(error: unknown) {
  if (!(error instanceof ApiError)) return undefined;
  const retryAt = error.meta?.retry_after_at ?? error.meta?.reset_at;
  return typeof retryAt === "string" ? retryAt : undefined;
}
