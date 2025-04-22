import React from 'react'

const Home: React.FC = () => {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-14 items-center">
          
        </div>
      </header>
      <main className="flex-1">
        <section className="w-full py-12 md:py-24 lg:py-32 xl:py-48">
          <div className="container px-4 md:px-6">
            <div className="grid gap-6 lg:grid-cols-[1fr_400px] lg:gap-12 xl:grid-cols-[1fr_600px]">
              <div className="flex flex-col justify-center space-y-4">
                <div className="space-y-2">
                  <h1 className="text-3xl font-bold tracking-tighter sm:text-5xl xl:text-6xl/none">
                    Streamline Your Engineering 1:1s
                  </h1>
                  <p className="max-w-[600px] text-muted-foreground md:text-xl">
                    Track GitHub metrics, Jira tasks, and follow-up items all in one place. Make your 1:1s more
                    productive and data-driven.
                  </p>
                </div>
                <div className="flex flex-col gap-2 min-[400px]:flex-row">
                  <a href="/register">
                    <a className="w-full">
                      Get Started
                    </a>
                  </a>
                  <a href="/demo">
                    <a className="w-full">
                      View Demo
                    </a>
                  </a>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>
      <footer className="border-t py-6 md:py-0">
        <div className="container flex flex-col items-center justify-between gap-4 md:h-24 md:flex-row">
          <p className="text-center text-sm leading-loose text-muted-foreground md:text-left">
            © 2024 EM Dashboard. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  )
}

export default Home 