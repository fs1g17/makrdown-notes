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

const schema = z.object({
  username: z
    .string()
    .trim()
    .min(1, { message: "Input a username" })
    .max(50, { message: "username can't be longer than 50 characters" }),
  email: z.email(),
  password: z.string().trim().min(8, { message: "Password must be at least 8 characters" }),
  confirmPassword: z.string().trim().min(8, { message: "Password must be at least 8 characters" }),
})
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords must match",
    path: ["confirmPassword"]
  });

export default function SignUp() {
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
    }
  });

  const { mutate: register } = useMutation({
    mutationFn: ({
      username,
      email,
      password,
    }: {
      username: string;
      email: string;
      password: string;
    }) =>
      clientFetch.post("/api/user/register", {
        username,
        email,
        password,
      }),
    onSuccess: () => {
      console.log("Registration successful");
    },
    onError: () => {
      console.log("Registration failed");
    },
  });

  const handleSubmit = ({
    username,
    email,
    password,
  }: {
    username: string;
    email: string;
    password: string;
  }) => {
    register({ username, email, password });
  }

  return (
    <div className="flex min-h-screen items-center justify-center p-4">
      <Card className="w-full max-w-md p-8">
        <form id="signup-form" onSubmit={form.handleSubmit(handleSubmit)}>
          <FieldGroup>
            <FieldSet>
              <FieldLegend className="text-2xl font-bold">Create Account</FieldLegend>
              <FieldDescription>
                Sign up to start organizing your markdown notes
              </FieldDescription>

              <FieldGroup>
                <Controller
                  name="username"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="signup-username">Username</FieldLabel>
                      <Input
                        {...field}
                        id="signup-username"
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
                  name="email"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="signup-email">Email</FieldLabel>
                      <Input
                        {...field}
                        id="signup-email"
                        type="email"
                        placeholder="you@example.com"
                        aria-invalid={fieldState.invalid}
                        autoComplete="email"
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
                      <FieldLabel htmlFor="signup-password">Password</FieldLabel>
                      <Input
                        {...field}
                        id="signup-password"
                        type="password"
                        placeholder="Create a password"
                        aria-invalid={fieldState.invalid}
                        autoComplete="new-password"
                      />
                      <FieldDescription>
                        Must be at least 8 characters
                      </FieldDescription>
                      {fieldState.invalid && (
                        <FieldError errors={[fieldState.error]} />
                      )}
                    </Field>
                  )}
                />

                <Controller
                  name="confirmPassword"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Field data-invalid={fieldState.invalid}>
                      <FieldLabel htmlFor="signup-confirm-password">Confirm Password</FieldLabel>
                      <Input
                        {...field}
                        id="signup-confirm-password"
                        type="password"
                        placeholder="Confirm your password"
                        aria-invalid={fieldState.invalid}
                        autoComplete="new-password"
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
                Sign Up
              </Button>
            </Field>

            <p className="text-center text-sm text-muted-foreground">
              Already have an account?{" "}
              <Link href="/sign-in" className="text-primary hover:underline">
                Sign in
              </Link>
            </p>
          </FieldGroup>
        </form>
      </Card>
    </div>
  )
}