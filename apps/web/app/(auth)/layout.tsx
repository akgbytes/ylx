import Link from "next/link";
import Image from "next/image";

export default function AuthLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <main
      id="main-content"
      className="grid min-h-svh place-items-center bg-muted p-6"
    >
      <div className="flex w-full max-w-sm flex-col gap-6">
        <Link href="/" className="self-center text-xl font-semibold">
          <Image
            src="/brand/favicon-64.png"
            alt="YLX"
            className="size-10"
            width="40"
            height="40"
          />
        </Link>
        {children}
      </div>
    </main>
  );
}
