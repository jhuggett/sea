package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

const OptInDebug slog.Level = -5

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

type Handler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex

	Level    slog.Level
	UseColor bool

	Allowlist []string
	BlockList []string

	WriteToFile string // path to file
}

func (h *Handler) colorize(colorCode int, v string) string {
	if !h.UseColor {
		return v
	}
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

const (
	timeFormat = "15:04:05.000"
)

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Custom slog handler panicked:", err)
			fmt.Println(r)
		}
	}()

	level := r.Level.String()
	message := r.Message

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	source := attrs["source"].(map[string]any)
	delete(attrs, "source")

	if len(h.Allowlist) > 0 {
		if _, ok := source["file"]; ok {
			for _, allow := range h.Allowlist {
				if strings.Contains(source["file"].(string), allow) {
					goto Allowed
				}
			}
		}
	}

	if len(h.BlockList) > 0 {
		if _, ok := source["file"]; ok {
			for _, block := range h.BlockList {
				if strings.Contains(source["file"].(string), block) {
					return nil
				}
			}
		}
	}

	if len(h.Allowlist) == 0 {
		goto Allowed
	}

	return nil

Allowed:

	pkg := source["function"].(string)
	bits := strings.Split(pkg, "/")
	pkg = bits[len(bits)-1]
	bits = strings.Split(pkg, ".")
	pkg = bits[0]

	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	if h.WriteToFile != "" {
		file, err := os.OpenFile(h.WriteToFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error when opening file for writing: %w", err)
		}
		_, err = file.Write([]byte(fmt.Sprintf("%s %s %s %s (%s:%v %s) %s\n", level, pkg, r.Time.Format(timeFormat), message, source["file"], source["line"], source["function"], string(bytes))))
		if err != nil {
			return fmt.Errorf("error when writing to file: %w", err)
		}
	}

	switch r.Level {
	case slog.LevelDebug:
		level = h.colorize(cyan, "🟦 DBG")
		message = h.colorize(lightGray, message)
	case slog.LevelInfo:
		level = h.colorize(green, "🟩 INF")
	case slog.LevelWarn:
		level = h.colorize(yellow, "🟨 WRN")
		message = h.colorize(yellow, message)
	case slog.LevelError:
		level = h.colorize(red, "🟥 ERR")
		message = h.colorize(red, message)
	}

	if len(attrs) == 0 {
		bytes = []byte{}
	}

	fmt.Println(
		level,
		h.colorize(magenta, pkg),
		h.colorize(darkGray, r.Time.Format(timeFormat)),
		message,
		h.colorize(darkGray, fmt.Sprintf("(%s:%v %s)", source["file"], source["line"], source["function"])),
		h.colorize(darkGray, string(bytes)),
	)

	return nil
}

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

type HandlerOptions struct {
	slog.HandlerOptions

	UseColor bool

	Allowlist []string
	BlockList []string

	WriteToFile string // path to file
}

func NewHandler(opts *HandlerOptions) *Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &Handler{
		b: b,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		m:        &sync.Mutex{},
		Level:    opts.Level.Level(),
		UseColor: opts.UseColor,

		Allowlist: opts.Allowlist,
		BlockList: opts.BlockList,

		WriteToFile: opts.WriteToFile,
	}
}

func RandID() string {
	return fmt.Sprintf("%d", rand.Int63())
}

func Package(name string) *slog.Logger {
	return slog.With("package", name)
}

func UnderTest() {
	slog.SetDefault(
		slog.New(NewHandler(&HandlerOptions{
			HandlerOptions: slog.HandlerOptions{
				AddSource: true,
				Level:     OptInDebug,
			},
			UseColor: true,
		})),
	)
}
