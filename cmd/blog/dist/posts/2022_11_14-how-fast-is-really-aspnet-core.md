<!--
    Tags: aspnet-core dotnet-core csharp
-->

# How fast is ASP.NET Core?

In recent years the .NET Team has been heavily advertising ASP.NET Core as one of the fastest web frameworks on the market. The source of those claims has always been the [TechEmpower Framework Benchmarks](https://www.techempower.com/benchmarks).

Take this slide from [BUILD 2021](https://www.youtube.com/watch?v=2Ky28Et3gy0), which [Scott Hunter](https://twitter.com/coolcsh) - Director of Program Management .NET, presented last year:

![Dubious .NET 5 performance claims](https://cdn.dusted.codes/images/blog-posts/2022-11-14/dotnet-5-performance-claims.png)

According to him .NET is **more than 10 times faster** than Node.js.

Scott also claims that .NET is faster than Java, Go and **even C++**, which is a huge boast if this is true!

Only recently [Sébastien Ros](https://twitter.com/sebastienros), from the ASP.NET Core team, wrote this on Reddit:

![Reddit comment by member of the ASP.NET Core team](https://cdn.dusted.codes/images/blog-posts/2022-11-14/reddit-comment.png)

In particular this sentence was super interesting to read:

> Finally, even with the fastest Go web framework, .NET is still faster when using a high level stack (middleware, minimal APIs, …).

That is a bold claim and equally super impressive if true, so I was naturally curious to find out more about ASP.NET Core's performance and the TechEmpower Framework Benchmarks.

## TechEmpower Benchmarks

[TechEmpower](https://www.techempower.com) is a software agency located in Los Angeles, California who run an independent framework benchmark on their servers. They publish all the [results on their website](https://www.techempower.com/benchmarks/) as well as the [framework code on GitHub](https://github.com/TechEmpower/FrameworkBenchmarks).

The first thing that stood out to me was that the last official round ([Round 21](https://www.techempower.com/benchmarks/#section=data-r21)) was captured on 19th July 2022. The round before that ([Round 20](https://www.techempower.com/benchmarks/#section=data-r20)) ran in February 2021, which means there was a gap of more than a year between those two official rounds. I am not sure why they have only so few official rounds but I also discovered that they have a continuos benchmark run which can be viewed on their [Results Dashboard](https://tfb-status.techempower.com). However, since the last official round was not that long ago and the difference between the results from Round 21 and the [last completed run from the continuous benchmarks](https://www.techempower.com/benchmarks/#section=test&runid=da435f15-1b5b-4347-acbe-a68ced6efb39) is not that big I decided to stick with Round 21 for my further analysis.

TechEmpower divides their tests into the following categories:

- JSON serializers
- Single query
- Multiple queries
- Cached queries
- Fortunes
- Data updates
- Plaintext

The [Fortunes](https://www.techempower.com/benchmarks/#section=data-r21&test=fortune) benchmark is the gold standard of all benchmarks. It is the only one which tries to resemble a "real world scenario" which involves some reading from a database, sorting data by text, XSS prevention and it includes some server-side HTML template rendering too.

All the other test categories focus on an isolated aspect of a framework which makes it interesting for reading but useless when ranking web frameworks by general performance.

So let's take a closer look at the Fortunes benchmark from Round 21:

[![TechEmpower Benchmark Results Top 20 from Round 21](https://cdn.dusted.codes/images/blog-posts/2022-11-14/techempower-benchmarks-round-21.png)](https://cdn.dusted.codes/images/blog-posts/2022-11-14/techempower-benchmarks-round-21.png)

To my astonishment ASP.NET Core ranks 9th in place amongst the top 10 fastest frameworks! Two further flavours of the ASP.NET Core benchmark also rank 13th and 14th out of the 439 completed benchmark runs. That is very impressive indeed!

### What are the different ASP.NET Core benchmarks?

Why does ASP.NET Core appear more than once in the benchmark results with varying performance metrics?

It turns out that there are in fact 15 different ASP.NET Core benchmarks which can be broadly subdivided into these four categories:

- ASP.NET Core stripped naked
- ASP.NET Core with middleware
- ASP.NET Core MVC
- ASP.NET Core on Mono

![ASP.NET Core Benchmark Frameworks](https://cdn.dusted.codes/images/blog-posts/2022-11-14/aspnet-core-benchmark-frameworks.png)

However, those are self-chosen names (by the .NET Team) and in order to get a real picture of what is being tested one has to look at the actual code itself. Luckily all the [code is publicly available on GitHub](https://github.com/TechEmpower/FrameworkBenchmarks).

I'm not interested in checking out 15 different implementations of various ASP.NET Core benchmarks so I decided to focus on the top performing ones by further narrowing down the 15 benchmarks into the best 7 out of the bunch:

![ASP.NET Core Benchmark Frameworks without MySQL and without Mono tests](https://cdn.dusted.codes/images/blog-posts/2022-11-14/aspnet-core-benchmark-frameworks-without-mysql-and-without-mono.png)

I removed the Mono benchmarks and all the tests which used MySQL as the underlying database, because those tests performed significantly worse in comparison to the .NET Core with Postgres equivalents (which has the `pg` suffix in the labels).

Slowly the picture becomes clearer. The above screenshot also includes the framework "classification" which can be seen on the right hand side of the image. The top benchmark (which is the impressive one that ranks 9th overall) is classified as **"Platform"**. The next three benchmarks are classified as **"Micro"** and the last three benchmarks are classified as **"Full"**. There seems to be a very significant performance drop as one moves from the "Platform" tests down to the "Full" tests.

Similar to the naming of the framework benchmarks, the classification is not standardised or audited by TechEmpower employees either. Anyone can submit code with an arbitrary name and classification and get very little or no scrutiny at all by the repository maintainers. At least that was my impression (I once submitted an F# benchmark test).

Only the code itself can be used as a reliable source of truth to draw conclusions from those tests.

Luckily the code for all ASP.NET Core (on .NET Core) benchmarks can be found inside the [/frameworks/CSharp/aspnetcore](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore) folder of the GitHub repository.

On 19th July 2022 (when Round 21 took place) the ASP.NET Core benchmark was divided into two projects:

- [/Benchmarks](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/Benchmarks)
- [/PlatformBenchmarks](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/PlatformBenchmarks)

Both of these web applications are **very different** so it is important to understand which one is used by which benchmark. This can be done by inspecting the [`config.toml`](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/config.toml) file and the associated `Dockerfile` for the respective test case.

For example, the best ranking ASP.NET Core benchmark (`aspcore-ado-pg`) has the following configuration:

##### config.toml

```
[ado-pg]
urls.db = "/db"
urls.query = "/queries/"
urls.fortune = "/fortunes"
urls.cached_query = "/cached-worlds/"
approach = "Realistic"
classification = "Platform"
database = "Postgres"
database_os = "Linux"
os = "Linux"
orm = "Raw"
platform = ".NET"
webserver = "Kestrel"
versus = "aspcore-ado-pg"
```

##### aspcore-ado-pg.dockerfile

```
FROM mcr.microsoft.com/dotnet/sdk:6.0.100 AS build
WORKDIR /app
COPY PlatformBenchmarks .
RUN dotnet publish -c Release -o out /p:DatabaseProvider=Npgsql

FROM mcr.microsoft.com/dotnet/aspnet:6.0.0 AS runtime
ENV ASPNETCORE_URLS http://+:8080

# Full PGO
ENV DOTNET_TieredPGO 1
ENV DOTNET_TC_QuickJitForLoops 1
ENV DOTNET_ReadyToRun 0

WORKDIR /app
COPY --from=build /app/out ./
COPY PlatformBenchmarks/appsettings.postgresql.json ./appsettings.json

EXPOSE 8080

ENTRYPOINT ["dotnet", "PlatformBenchmarks.dll"]
```

The Dockerfile tells us that this test uses the `/PlatformBenchmakrs` code:

```
COPY PlatformBenchmarks .
```

From the `config.toml` file we can derive that the Fortune test invokes the `/fortunes` endpoint during the benchmark run.

Also the .NET Team specified this particular benchmark to be classified as a realistic approach in the `config.toml` file:

```
approach = "Realistic"
```

## The "ASP.NET Core Platform" Benchmark

Cool, so what's inside this highly performant realistic ASP.NET Core application?

![ASP.NET Core PlatformBenchmarks code repository](https://cdn.dusted.codes/images/blog-posts/2022-11-14/aspnet-core-source-code.png)

On first glance I didn't recognise a lot of what I'd normally consider a typical ASP.NET Core application (I've been developing professionally on ASP.NET and later ASP.NET Core since 2010).

The only thing that looked slightly familiar was the use of Kestrel (the .NET web server) inside [Program.cs](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/PlatformBenchmarks/Program.cs#L62-L73):

![Kestrel setup](https://cdn.dusted.codes/images/blog-posts/2022-11-14/kestrel-setup.png)

To my surprise this was also the **only thing** which I could recognise as an "ASP.NET Core" thing. The web application itself is not even initialised via one of the many ASP.NET Core idioms. Instead it creates a custom `BenchmarkApplication` as the listener on the configured endpoint.

An untrained eye might be thinking that `builder.UseHttpApplication<T>()` is a method that comes with Kestrel, but that is not the case either. The extension method as well as the `HttpApplication` class [which is in use here](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/PlatformBenchmarks/HttpApplication.cs) are not things which you'd find in the actual ASP.NET Core framework. It is yet another custom class specifically written for this benchmark:

![Fake HttpApplication](https://cdn.dusted.codes/images/blog-posts/2022-11-14/fake-aspnet-core-http-application.png)

Not even the interface `IHttpApplication` comes from ASP.NET Core. This is also a [custom made type](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/PlatformBenchmarks/IHttpConnection.cs) which was specifically designed for the benchmark tests.

Looking further into the `BenchmarkApplication.cs` I was shocked by the sheer [amount of finely tuned low level C# code](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/PlatformBenchmarks/BenchmarkApplication.cs) that was tailor made for this (extremely simple) application.

Everything inside the `/PlatformBenchmarks` folder is custom code which you won't find anywhere in an official ASP.NET Core package.

A good example is the `AsciiString` class which is used to statically initialise huge chunks of the expected HTTP responses in advance:

![AsciiString Usage](https://cdn.dusted.codes/images/blog-posts/2022-11-14/ascii-strings.png)

Even though it is called `AsciiString` it's only a string in name:

![AsciiString Implementation](https://cdn.dusted.codes/images/blog-posts/2022-11-14/ascii-string-implementation.png)

In reality the `AsciiString` class is just a fancy (highly optimised) wrapper around a byte array which converts a string into bytes during initialisation. In the case of the Fortunes test the entire HTTP header (which the application is expected to return during a test run) is created upfront during application startup and then kept in memory for the entirety of the benchmark:

![Hardcoded HTTP Headers](https://cdn.dusted.codes/images/blog-posts/2022-11-14/http-header-trick.png)

This is supposed to be a very simple application, something which a framework could probably squeeze into a single file of code, but the `/PlatformBenchmarks` project has **many dozens of expertly crafted classes** with all sorts of **trickery** applied to produce a desired outcome.

The extent to which the .NET Team went is extraordinary.

ASP.NET Core has many ways of implementing routing. They have Actions and Controllers, Endpoint Routing, Minimal APIs, or if someone wanted to operate on the **lowest level of ASP.NET Core** (= Platform), then they could work directly with the `Request` and `Response` objects from the `HttpContext`.

Neither of these options can be found `/PlatformBenchmarks`:

![Highly optimised routing](https://cdn.dusted.codes/images/blog-posts/2022-11-14/optimised-routing.png)

In fact, you won't find a `HttpContext` anywhere at all. It's almost like the .NET Team tried to avoid using ASP.NET Core at all cost, which is strange to say the least.

Sieving through the project reveals even more bizarre code which the .NET Team applied to "tweak" the benchmark score.

For instance take a look at the HTML templating implementation of the ASP.NET Core solution:

![ASP.NET Core Fortunes output writer](https://cdn.dusted.codes/images/blog-posts/2022-11-14/aspnet-core-fortunes-test.png)

There is no HTML template at all. The whole point of the Fortunes benchmark is - amongst others - to test different web frameworks for how fast they can output templated HTML. In ASP.NET Core we have two templating engines, [Razor Views](https://learn.microsoft.com/en-us/aspnet/core/mvc/views/razor?view=aspnetcore-7.0) and [Razor Pages](https://learn.microsoft.com/en-us/aspnet/core/razor-pages/?view=aspnetcore-7.0&tabs=visual-studio), of which none is being used here.

Instead there are more hardcoded statically initialised byte arrays:

![HTML Template Rendering Cheat](https://cdn.dusted.codes/images/blog-posts/2022-11-14/html-template-cheat.png)

Of course the question remains if these sort of tricks are allowed? The lines might be a bit blurry but I am certain that this implementation pushes the boundaries of what one might consider a real templating engine.

Web frameworks don't have to participate in every category of the TechEmpower Benchmark tests. In fact it is encouraged to only enter the categories which apply to a particular framework. For example, if a low level ASP.NET Core implementation (a real one which uses ASP.NET Core with `HttpContext` and so on) doesn't have template rendering included then it shouldn't enter the competition for Fortunes. If a higher level framework such as ASP.NET Core MVC has HTML template rendering available then it can enter the Fortunes benchmark. Entering the Fortunes competition with random C# code that doesn't resemble a real web framework at all makes very little sense and really just tarnishes the credibility of the entire TechEmpower Framework Benchmark test.

Perhaps I am being a little bit overly critical here, but this line of code really got me thinking:

![Date Header Cheat](https://cdn.dusted.codes/images/blog-posts/2022-11-14/writing-date-header.png)

Setting the `Date` HTTP header with a date time value is such a small task that you don't even need a framework to do this job. It should be no more than a single line of code:

```
response.WriteHeader("Date", DateTime.UtcNow.ToString())
```

However, the ASP.NET Core benchmark has a "*slightly more optimised*" [solution](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/PlatformBenchmarks/DateHeader.cs) to this task:

![Date Header Code](https://cdn.dusted.codes/images/blog-posts/2022-11-14/date-header-implementation.png)

Setting a date time value has been so highly optimised that I can't even fit the entire code into a single screen. **The creativity of finding ways to save computation cycles and therefore score higher in the benchmarks is truly astonishing.** The `DateHeader` class is a static class (which means it only gets initialised once as a singleton and is then kept in memory) with a static `DateTimeOffset` value (of course already stored as a byte array). Additionally a `System.Threading.Timer` object is also statically initialised with a **one second** interval. This [Timer](https://learn.microsoft.com/en-us/dotnet/api/system.threading.timer?view=net-7.0) will run on a separate thread and set a new date time value once every second:

```
private static readonly Timer s_timer = new Timer((s) => {
    SetDateValues(DateTimeOffset.UtcNow);
}, null, 1000, 1000);
```

You wonder how this is an optimisation? Well, the TechEmpower Benchmark will hit a web server many hundreds of thousand times **per second** to really test the limits of each framework. The `DateHeader` class will return the exact same timestamp for all those thousand requests and henceforth save itself from computing a new timestamp many thousand times. Then after one second the `Timer` (which runs on a separate thread) will sync a new timestamp exactly once and cache the timestamp for the next 300+ thousand requests. I'm impressed by the ingenuity. In all fairness the HTTP `Date` header doesn't accept timestamps more granular than a second and the [TechEmpower guidelines](https://github.com/TechEmpower/FrameworkBenchmarks/wiki/Project-Information-Framework-Tests-Overview) mention this to be an accepted optimisation.

The only question I have is if this benchmark is testing ASP.NET Core why does it need to replicate something which ASP.NET Core already has [out of the box](https://github.com/dotnet/aspnetcore/blob/v5.0.17/src/Servers/Kestrel/Core/src/Internal/Http/DateHeaderValueManager.cs)?

Now I ask myself, are all the ASP.NET Core benchmarks "tweaked" like this?

What about other frameworks?

I needed to further investigate this!

## ASP.NET Core Micro Benchmarks

After dissecting the "Platform" benchmark it was time to look at the "Micro" frameworks:

![ASP.NET Core Benchmark Frameworks without MySQL and without Mono tests](https://cdn.dusted.codes/images/blog-posts/2022-11-14/aspnet-core-benchmark-frameworks-without-mysql-and-without-mono.png)

Looking at the respective [Dockerfile](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/aspcore-mw-ado-pg.dockerfile) it turns out that the "Micro" benchmarks use the code from the `/Benchmarks` folder, which looks like an actual ASP.NET Core application:

![ASP.NET Core Benchmarks folder](https://cdn.dusted.codes/images/blog-posts/2022-11-14/aspnet-core-benchmarks-folder.png)

This benchmark immediately has a different vibe than the one before. I'm very pleased to see that it's actually using elements which come from ASP.NET Core itself. The Fortunes tests are initialised via conventional middleware like this:

![Fortunes Raw Middleware](https://cdn.dusted.codes/images/blog-posts/2022-11-14/fortunes-raw-middleware.png)

The `aspcore-mw-ado-pg` benchmark is what most .NET developers would probably call a low level "Platform" ASP.NET Core implementation. There is no higher level routing, no content negotiation, no other cross-cutting middlewares, no EntityFramework and still no actual HTML template rendering either, but at least it's ASP.NET Core.

The [middleware](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/Benchmarks/Middleware/FortunesRawMiddleware.cs) operates directly on the `HttpContext` to do basic routing:

![Middleware Routing](https://cdn.dusted.codes/images/blog-posts/2022-11-14/middleware-routing.png)

This is okay and [inline with the TechEmpower guidelines](https://github.com/TechEmpower/FrameworkBenchmarks/wiki/Project-Information-Framework-Tests-Overview), because operating directly on the `HttpContext` is canonical for the framework (as opposed to the benchmark before):

> In some cases, it is considered normal and sufficiently production-grade to use hand-crafted minimalist routing using control structures such as if/else branching. This is acceptable where it is considered canonical for the framework.

Although the middleware benchmark doesn't apply the `AsciiString` trickery any more, it still resorts to a "fake" templating engine:

![StringBuilder Templates](https://cdn.dusted.codes/images/blog-posts/2022-11-14/stringbuilder-templates.png)

Overall it is a much more realistic (albeit not perfect) benchmark!

## ASP.NET Core Full Benchmarks

Finally it was time to check out the "MVC" benchmarks. It also derives its code from the `/Benchmarks` folder but instead of operating on the raw `HttpContext` it actually initialises the [least required MVC middleware](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/Benchmarks/Startup.cs#L104-L114) with the Razor View Engine:

![MVC Core Middleware](https://cdn.dusted.codes/images/blog-posts/2022-11-14/mvc-middleware.png)

The Controller Action is also kept very realistic and finally uses the actual ASP.NET Core templating engine:

```
[HttpGet("raw")]
public async Task<IActionResult> Raw()
{
    var db = HttpContext.RequestServices.GetRequiredService<RawDb>();
    return View("Fortunes", await db.LoadFortunesRows());
}
```

The Razor view matches what one would expect from this simple benchmark:

![Razor View Template](https://cdn.dusted.codes/images/blog-posts/2022-11-14/mvc-view-template.png)

This is the most realistic ASP.NET Core application which actually meets the spirit of the Fortunes benchmark.

However, the results of this benchmark are **very different** to what Microsoft actively advertised to the .NET Community. The performance difference between a "fake" templating engine where a HTML response is being created **in memory** via a [cached StringBuilder](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/CSharp/aspnetcore/Benchmarks/Data/StringBuilderCache.cs) versus an actual templating engine which has to incur additional (expensive) I/O operations to read, parse and apply HTML templates from disk is enormous.

The latter only manages to serve **184k requests/sec** and **only ranks 109<sup>th</sup>** overall in the TechEmpower Framework Benchmarks for the Fortunes test. That is a staggering difference and something to be kept in mind when comparing ASP.NET Core to frameworks written in Java, Go or C++.

## Other frameworks

Now that I've established a clearer picture of what the various ASP.NET Core benchmarks are, it was time to look at other frameworks too.

### Java

The fastest Java benchmark which also uses Postgres as the underlying database is [Jooby](https://jooby.io).

Their [benchmark implementation](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Java/jooby) is astonishingly simple. The entire fortune implementation is basically this block of code:

![Jooby Fortunes](https://cdn.dusted.codes/images/blog-posts/2022-11-14/jooby-fortunes.png)

It uses a higher level router (`get("/fortunes", ctx -> {}`) as well as conventional database access methods and a real templating engine too:

![Jooby Template](https://cdn.dusted.codes/images/blog-posts/2022-11-14/jooby-template.png)

This is pretty much the Java equivalent to the ASP.NET Core MVC (aka Full) benchmark.

The interesting part is that this completely unoptimised fully fledged Java MVC framework **ranks overall 12<sup>th</sup>** in the Fortunes benchmark with an incredible **404k requests/sec**. It is essentially **more than twice as fast** as the ASP.NET Core equivalent, still beats the "Micro" implementation of the ASP.NET Core benchmark (which skips all the expensive I/O operations by using a fake templating engine) and even manages to compete with the infamous `/PlatformBenchmarks` application which in all honesty due to its differences is not even worth a comparison.

No disrespect to ASP.NET Core (because 184k requests/sec is still an amazing result) but it doesn't come anywhere near this Java framework when it comes to performance. Credit where credit is due.

### Go

What about Go?

Sébastien Ros (developer working on ASP.NET Core performance at Microsoft) specifically called out Go and claimed that ASP.NET Core is still faster than Go in a like-for-like comparison. I was personally very interested in this claim as I have migrated several .NET Core projects to Go and seen dramatic performance increases as a result of it.

At the time of writing this post the fasted Fortune benchmark is [atreugo](https://github.com/savsgio/atreugo) for Go.

Similar to Java, the actual [Go implementation](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Go/atreugo) is kept extremely simple.

Routing is done via the framework provided idioms:

![atreugo routing](https://cdn.dusted.codes/images/blog-posts/2022-11-14/atreugo-routing.png)

No shortcuts or trickery to be found here. The entire application for the Fortunes benchmark is basically [less than 20 lines of code](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Go/atreugo/src/views/views.go#L97-L123).

Templating is done the proper way too:

![atreugo template](https://cdn.dusted.codes/images/blog-posts/2022-11-14/atreugo-template.png)

So where does this leave us overall? Well, just like with the fasted Java framework, the Go benchmark also compares to ASP.NET Core's "Full" implementation best. Anything other would simply not be fair. You cannot compare a benchmark which spits out in-memory crafted HTML (which is not even part of ASP.NET Core) versus one that actually uses a real templating engine that goes through expensive cycles of reading files from I/O, parsing them at runtime and having to execute their logic for every request (loops, variables, etc. in the template).

Nevertheless, the expensive Go implementation **ranks 22<sup>nd</sup>** overall in the TechEmpower Fortunes Benchmark with an equally impressive **381k requests/sec**. Not quite as fast as the Java one but **still more than 2x faster than the equivalent test in ASP.NET Core**.

### C++

Hopefully this shouldn't be a big surprise, but currently C++ with the [drogon](https://github.com/drogonframework/drogon) framework **leads the Fortunes** benchmarks with a breathtaking **616k requests/sec** which beats every other framework by a long stretch (except Rust where the gap is not that big)! What makes this achievement even more astonishing is that it manages to do this with a [fully fledged MVC implementation](https://github.com/drogonframework/drogon). There is absolutely no shortcuts or trickery at play.

It even uses the [CSP templating engine](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/C%2B%2B/drogon/drogon_benchmark/views) which looks like this:

![drogon template](https://cdn.dusted.codes/images/blog-posts/2022-11-14/drogon-template.png)

I love .NET but there is no mental gymnastics that one could convincingly apply in which .NET comes on top of C++. Any benchmark that suggests otherwise knows it's not being honest with itself.

### Rust, Node.js, Kotlin and PHP

Since the .NET Team started to campaign ASP.NET Core as a much faster web framework than many others I thought it would only be fair to further probe those claims.

#### Rust

[Rust](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Rust/xitca-web) delivers **588k requests/sec** and comes **2<sup>nd</sup>** in the overall Fortunes benchmark. It's the only other language platform which gives C++ a run for its money. The [xitca-web](https://github.com/HFQR/xitca-web) framework accomplishes this unbelievable result with another proper [MVC-like implementation](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Rust/xitca-web/src/main.rs#L130-L136) and an [actual templating engine](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Rust/xitca-web/templates/fortune.stpl).

#### Kotlin

Another great result is achieved by a [Kotlin web framework](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Kotlin/vertx-web-kotlin-coroutines) with a very honest [Fortunes implementation](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Kotlin/vertx-web-kotlin-coroutines/src/main/kotlin/io/vertx/benchmark/App.kt#L147-L172) which uses the [Rocker engine](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/Kotlin/vertx-web-kotlin-coroutines/src/main/resources/templates/Fortunes.rocker.html) for its HTML templating. It pegs at **350k requests/sec** and comes **29<sup>th</sup>** overall which is still **80 places ahead of the equivalent ASP.NET Core** implementation.

#### Node.js

One claim which turned out to be (partially) true is that **ASP.NET Core is faster than Node.js**. Although only **3x and not 10x faster** as it was claimed, ASP.NET Core still beats [Polkadot](https://github.com/lukeed/polkadot) convincingly, which is the highest ranking Node.js framework which had a comparable implementation to the "Micro" benchmark in ASP.NET Core. With only **125k requests/sec** Node.js trails behind .NET.

#### PHP

Now this might actually take people by surprise, but if you haven't been paying attention then you might have missed all the work that has gone into PHP over the many years. Not least because Facebook invested a lot of effort into making PHP a better platform. It is now capable of serving an incredible **309k requests/sec** with it's [MVC-like implementation](https://github.com/TechEmpower/FrameworkBenchmarks/blob/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/PHP/mixphp/views/fortunes.php) delivered by [mixphp](https://github.com/TechEmpower/FrameworkBenchmarks/tree/62aaac842e6bf51540bb838bb9ffaaad0d7c9e73/frameworks/PHP/mixphp). That is still significantly faster than ASP.NET Core's MVC framework and certainly deserves a mention too!

#### Just(js)

If you are a **JavaScript** developer don't feel too bad about the Node.js benchmarks, because [Just(js)](https://github.com/just-js/just) will knock you off your socks with a spectacular **538k requests/sec**. This is no joke, [Just(js)](https://just.billywhizz.io) comes **5<sup>th</sup> in the Fortunes benchmark** and is the only framework which competes in the realms of C++ and Rust. It is a remarkable achievement which is [not something that happened by mistake](https://just.billywhizz.io/blog/on-javascript-performance-01/). It is far ahead of every other ASP.NET Core benchmark and had to be mentioned as part of this post!

## Is ASP.NET Core actually fast?

**Yes**, it certainly is!

Especially if you think back to what Classic ASP.NET was during the .NET Framework times then it becomes very clear that ASP.NET Core is world's apart from its darker past.

Make no mistake, **ASP.NET Core is very fast** and certainly doesn't need to shy away from a healthy competition. However, it is **evidently not faster than Java, Go or C++**. Perhaps it will get there one day but at the moment this is not the case. I am certain that we haven't seen the ceiling for ASP.NET Core just yet and I look forward to what the .NET Team will deliver next. ASP.NET Core is a great platform and even though it's not the fastest (yet), it is still a joy!

I wish Scott Hunter and the rest of the ASP.NET Core Team didn't feel the need to market ASP.NET Core based on soft lies and bad-faith claims to make ASP.NET Core stand out amongst its peers. I'm sure there is more to be proud of!

#### Sidenotes

One final interesting thing which came up during my research is that TechEmpower [switched their cloud hosting environment from AWS to Azure](https://www.techempower.com/benchmarks/#section=environment) around the time when Microsoft got interested in the tests. TechEmpower also receives its physical hardware for all their on-premise tests by Microsoft today.