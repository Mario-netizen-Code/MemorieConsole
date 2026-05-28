# MemorieConsole

Real-time console log viewer for JSONL session memory files. Tails `.txt` files in `.memories/` and displays entries as colored cards with ANSI formatting.

## Features

- **Full-screen takeover** — clears entire terminal buffer and hides cursor on start
- **Colored cards** — each entry type has a distinct color, icon, and background badge
- **Multi-file tailing** — watches multiple `.txt` files simultaneously
- **Auto-detect new files** — polls `.memories/` and loads new files on-the-fly
- **Graceful shutdown** — Ctrl+C restores cursor and closes file handles
- **Cross-platform** — Windows (native Win32 API) and Unix (ANII + COLUMNS env)
- **Windows file locking** — opens files with `FILE_SHARE_READ | FILE_SHARE_WRITE | FILE_SHARE_DELETE`

## Build

```sh
go build -o tailer.exe ./cmd/tailer/
```

## Usage

```sh
tailer.exe                    # reads .memories/*.txt
tailer.exe file1.txt file2.txt  # reads specific files
```

Press **Ctrl+C** to exit.

## Structure

```
cmd/tailer/main.go              entry point
internal/app/
  app.go                        main loop + signal handling
  display.go                    ANSI formatting, card rendering, entry types
  tailer.go                     file tailing, reading, glob
  platform_windows.go           Win32 console/screen/file API
  platform_other.go             Unix ANII + file fallbacks
```

## Entry Format

Each line in the `.txt` files is a JSON object:

```json
{"id":"mem_0001","timestamp":"2026-05-27T12:00:00+02:00","session_id":1,"scope":"project","type":"modification","summary":"Fixed issue","files":["main.go"],"tags":["bugfix"]}
```

## License

MIT
