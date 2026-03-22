type RoutePageProps = {
  title: string
  description: string
}

export function RoutePage({ title, description }: RoutePageProps) {
  return (
    <main className="mx-auto flex min-h-svh w-full max-w-4xl flex-col gap-3 p-6">
      <h1 className="text-2xl font-semibold">{title}</h1>
      <p className="text-sm text-muted-foreground">{description}</p>
    </main>
  )
}
