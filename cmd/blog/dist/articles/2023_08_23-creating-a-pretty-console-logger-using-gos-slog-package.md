<!--
    Tags: golang slog logging
-->

# Creating a pretty console logger using Go's slog package

I had the privilege of attending [GopherCon UK](https://www.gophercon.co.uk) last week, and among the many captivating talks, one that stood out to me was "Structured Logging for the Standard Library" presented by [Jonathan Amsterdam](https://twitter.com/JonathanAmster2).

The presentation provided an insightful dive into [Go's `log/slog` package](https://pkg.go.dev/log/slog). This talk couldn't have come at a better time given that I've just started on a new Go project, where I was eager to use Go's structured logging approach. The `slog` package in Go distinctly draws its inspiration from [Uber's `zap`](https://github.com/uber-go/zap), making it a seamless transition for those who are well-acquainted with the latter. If you're already at ease with `zap`, you'll find yourself quickly at home with `slog` as well.

Currently, the `slog` library offers two built-in logging handlers: the `TextHandler` and the `JSONHandler`. The `TextHandler` formats logs as a series of `key=value` pairs, while the `JSONHandler` produces logs in JSON format. These handlers are greatly optimised for production scenarios, but can be somewhat verbose when troubleshooting applications during the development phase.

Recognizing this, I realized the necessity for a more visually friendly console logger tailored for local development purposes. Despite stumbling upon a code snippet within a blog post titled [A Comprehensive Guide to Logging in Go with Slog](https://betterstack.com/community/guides/logging/logging-in-go/), the implementation was both incomplete and flawed, rendering it unsuitable for my use case.

So my next question was: How difficult can it be to create one myself? Luckily, with a bit of clever hackery not that difficult at all!

## Objectives

First let's address some of the shortcomings in the implementation outlined in the blog post referenced above. Unfortunately the custom handler fails to print preformatted attributes coming from the `WithAttrs` method. This means that attributes set on a parent logger are not propagated to child loggers at all. Additionally, the handler struggles to manage groups established on a parent logger too. Alongside these issues, there was no support for appending an `error` attribute, as well as addressing other edge cases which the original `JSONHandler` intuitively dealt with.

This didn't come as a huge surprise, given the difficulty of crafting a custom `slog.Handler` as highlighted by the Go team themselves. In fact, writing a custom `slog.Handler` is not to be taken lightly, and the Go team anticipates that only a select number of package authors will find themselves in the need of undertaking this task. To facilitate this, the Go team has thoughtfully provided a [comprehensive guide to writing slog handlers](https://github.com/golang/example/tree/master/slog-handler-guide) to assist with this process.

Either way, I have no desire of writing a complete `slog.Handler` myself for something which is only ever going to be relevant during development time. With this in mind I set myself the following requirements:

#### Requirements:

- Logs must be visually pleasing
- Implementation must be complete
- Only use packages from the standard library
- Keep it super simple (as I'm lazy)

#### Non requirements:

- Doesn't have to be very fast
- Doesn't have to be very memory efficient

It's good to keep in mind that this "pretty" handler is tailored for development purposes and doesn't require blazing speed or intense memory efficiency. Since I won't be generating millions of logs on my local machine, this greatly simplifies the upcoming solution.

#### Final output:

The final logs should look something like this:

![Example 1](https://cdn.dusted.codes/images/blog-posts/2023-08-23/prettylog-example-1.png)

Here's another example by enabling debug logs and adding source information to them:

![Example 2](https://cdn.dusted.codes/images/blog-posts/2023-08-23/prettylog-example-2.png)

Evidently, these logs are designed for human readability through colouring and spacing. The default log attributes (time, level, message) are presented in a single line, while extra structured attributes are attached as a JSON object.

If you like this log style then keep on reading.

## Creating a pretty console logger

For the purpose of this blog post I am calling this package `prettylog` but you can copy paste this logger into your own codebase and call it whatever you want.

Let's start with the function that will add color to the console output:

```go
package prettylog

import (
	"fmt"
	"strconv"
)

const (
	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}
```

That's all that's required for straightforward coloured output, avoiding the need for an external dependency on the [color](https://github.com/fatih/color) package. I'm not even going to use all the colours listed above but I've included them regardless so you can adjust your logs to your own liking.

Moving forward, we'll create a struct called `Handler` (later used as `prettylog.Handler`):

```go
type Handler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}
```

The handler has three dependencies:

- A "nested" `slog.Handler` which we wrap to effectively fulfil most of our handler's logic.
- A `*bytes.Buffer` with the purpose to capture the output from the "nested" handler.
- A mutex to guarantee thread safe access to our `*bytes.Buffer`.

All three dependencies will make more sense once we implement the `Handle` function.

The `slog.Handler` interface requires four methods to be implemented:

- `Enabled`
- `WithAttrs`
- `WithGroup`
- `Handle`

The `Enabled` method denotes whether a given handler handles a `slog.Record` of a particular `slog.Level`. The `WithAttrs` and `WithGroup` methods create child loggers with predefined Attrs.

For all three methods we can use the implementation of our nested handler:

```go
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}
```

The `Handle` method is where things get interesting.

Writing a log line is actually remarkably easy if one completely ignores groups and attributes to begin with:

```go
const (
	timeFormat = "[15:04:05.000]"
)

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	fmt.Println(
		colorize(lightGray, r.Time.Format(timeFormat)),
		level,
		colorize(white, r.Message),
	)

	return nil
}
```

What about printing all the attributes that are added to the `slog.Record` or a parent logger? This is where the bytes buffer and the nested handler come into play.

##### The concept is simple:

We'll invoke the `Handle` function of the nested handler, but have it write to the `*bytes.Buffer` instead of the final `io.Writer`. We'll exclude the default log attributes such as time, level, and message from the nested handler to prevent repetition. Then, we'll append the remaining output as an indented JSON string to our log line. Since loggers need to function correctly when a single `slog.Logger` is shared among multiple goroutines, we also need to synchronize read and write access to the `*bytes.Buffer` using the mutex.

Let's encapsulate this behaviour in a function called `computeAttrs`:

```go
func (h *Handler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	return attrs, nil
}
```

The `computeAttrs` works as following:

1. It initially locks the mutex to ensure synchronized access across all goroutines utilizing the same logger or a child logger that shares the same `*bytes.Buffer`.
2. It defers the process of resetting the buffer (necessary to prevent outdated Attrs from previous `Log` calls) and releasing the mutex once the task is complete.
3. The `Handle` function of the inner `slog.Handler` is then invoked. This is where we compute a JSON object within the `*bytes.Buffer`, leveraging the capabilities of a `slog.JSONHandler`.
4. Lastly, the JSON buffer is transformed into a `map[string]any` after which the resulting object is returned to the caller.

Now, let's revisit our own `Handle` function and integrate the following code:

```go
attrs, err := h.computeAttrs(ctx, r)
if err != nil {
    return err
}

bytes, err := json.MarshalIndent(attrs, "", "  ")
if err != nil {
    return fmt.Errorf("error when marshaling attrs: %w", err)
}
```

Through invoking `computeAttrs`, we can obtain the `attrs` map, which we subsequently convert into a neatly formatted (indented) JSON string using marshaling. Admittedly, this code isn't the most efficient approach (writing a JSON string into a buffer, deserializing it into an object, and then re-serializing it as a string), but it's worth mentioning that I couldn't identify a more effective method to obtain an indented JSON string from the `slog.JSONHandler`. Fortunately, as highlighted earlier, this handler isn't designed to achieve peak speed performance in any case.

Finally, we attach the formatted JSON string in a dark gray hue to our "pretty" log entry:

```go
fmt.Println(
    colorize(lightGray, r.Time.Format(timeFormat)),
    level,
    colorize(white, r.Message),
    colorize(darkGray, string(bytes)),
)
```

The final `Handle` method looks as following:

```go
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {

	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(darkGray, level)
	case slog.LevelInfo:
		level = colorize(cyan, level)
	case slog.LevelWarn:
		level = colorize(lightYellow, level)
	case slog.LevelError:
		level = colorize(lightRed, level)
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	fmt.Println(
		colorize(lightGray, r.Time.Format(timeFormat)),
		level,
		colorize(white, r.Message),
		colorize(darkGray, string(bytes)),
	)

	return nil
}
```

Only one last task remains. Currently, the nested `slog.Handler` writes the time, log level, and log message in addition to other custom attributes. However, since our handler is responsible for displaying these three default attributes, we need to configure the inner `slog.Handler` to bypass the `slog.TimeKey`, `slog.LevelKey` and `slog.MessageKey` attributes.

The most straightforward approach is to provide a function to the `ReplaceAttr` property of the `slog.HandlerOptions`. However, we wish to preserve the ability for an application to specify its individual `ReplaceAttr` function and `slog.HandlerOptions`. Therefore we must apply a final touch of trickery to "merge" a custom `ReplaceAttr` function with our own requirements:

```go
func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}
```

A helpful analogy for understanding the `suppressDefaults` function is to compare it to a middleware in an HTTP server. It takes in a `next` function that matches the same function signature as the `ReplaceAttr` property. It then performs filtering on `slog.TimeKey`, `slog.LevelKey`, and `slog.MessageKey` before continuing with `next` (if it's not nil).

With this in place, we're ready to create a constructor for our `prettylog.Handler` and assemble everything together:

```go
func NewHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &Handler{
		b: b,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		m: &sync.Mutex{},
	}
}
```

The entire code can be found on [GitHub](https://github.com/dusted-go/logging/tree/main).

## Final result

Below are a few examples of how those pretty logs will look like.

Example of a logger with no `*slog.HandlerOptions`:

![Example 3](https://cdn.dusted.codes/images/blog-posts/2023-08-23/prettylog-example-3.png)

Creating a child logger with an additional group of attributes attached to it:

![Example 4](https://cdn.dusted.codes/images/blog-posts/2023-08-23/prettylog-example-4.png)

Making sure custom `ReplaceAttr` functions are supported:

![Example 5](https://cdn.dusted.codes/images/blog-posts/2023-08-23/prettylog-example-5.png)

Hopefully this blog post proved to be useful. It certainly provided me a valuable exercise for delving into the new `log/slog` package and gaining a better understanding of Go's latest structured logging capabilities.