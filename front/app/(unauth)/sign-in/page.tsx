"use client";

import z from "zod";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card } from "@/components/ui/card";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
  FieldLegend,
  FieldSet,
} from "@/components/ui/field";
import Link from "next/link";
import clientFetch from "@/lib/client-side-fetching";
import { useMutation } from "@tanstack/react-query";
import { useRouter } from "next/navigation";

const schema = z.object({
  username: z
    .string()
    .trim()
    .min(1, { message: "Input a username" }),
  password: z.string().trim().min(1, { message: "Input a password" }),
});

export default function SignIn() {
  const router = useRouter();

  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      username: "",
      password: "",
    }
  });

  const { mutate: login } = useMutation({
    mutationFn: ({
      username,
      password,
    }: {
      username: string;
      password: string;
    }) =>
      clientFetch.post("/api/tokens/auth", {
        username,
        password,
      }),
    onSuccess: () => {
      router.push("/folders");
    },
    onError: () => {
      console.log("Login failed");
    },
  });

  const handleSubmit = ({
    username,
    password,
  }: {
    username: string;
    password: string;
  }) => {
    login({ username, password });
  }

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md p-8">
        <form id="signin-form" onSubmit={form.handleSubmit(handleSubmit)}>
          <FieldGroup>
            <FieldSet>
              <FieldLegend className="text-2xl font-bold">Sign In</FieldLegend>
              <FieldDescription>
                Welcome back! Sign in to access your notes
              </FieldDescription>

              <FieldGroup>
                <Controller
                  name="username"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="signin-username">Username</FieldLabel>
                      <Input
                        {...field}
                        id="signin-username"
                        placeholder="Enter your username"
                        aria-invalid={fieldState.invalid}
                        autoComplete="username"
                      />
                      {fieldState.invalid && (
                        <FieldError errors={[fieldState.error]} />
                      )}
                    </Field>
                  )}
                />

                <Controller
                  name="password"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="signin-password">Password</FieldLabel>
                      <Input
                        {...field}
                        id="signin-password"
                        type="password"
                        placeholder="Enter your password"
                        aria-invalid={fieldState.invalid}
                        autoComplete="current-password"
                      />
                      {fieldState.invalid && (
                        <FieldError errors={[fieldState.error]} />
                      )}
                    </Field>
                  )}
                />
              </FieldGroup>
            </FieldSet>

            <Field>
              <Button type="submit" className="w-full">
                Sign In
              </Button>
            </Field>

            <p className="text-center text-sm text-muted-foreground">
              Don&apos;t have an account?{" "}
              <Link href="/sign-up" className="text-primary hover:underline">
                Sign up
              </Link>
            </p>
          </FieldGroup>
        </form>
      </Card>
    </div>
  )
}
