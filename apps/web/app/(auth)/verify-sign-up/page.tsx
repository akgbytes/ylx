import { Suspense } from "react";

import { VerifySignUpForm } from "@/components/auth/verify-sign-up-form";

export default function VerifySignUpPage() {
  return (
    <Suspense>
      <VerifySignUpForm />
    </Suspense>
  );
}
