# Go NASA Deep Space Network Now

This is unofficial application for live monitoring the [NASA Deep Space Network Now](https://www.nasa.gov/directorates/somd/space-communications-navigation-program/what-is-the-deep-space-network/).

The NASA Deep Space Network (DSN) is a global system of large antennas and communication facilities used to track, send commands, and receive data from spacecraft exploring deep space. Operated by NASA’s Jet Propulsion Laboratory (JPL), the DSN supports missions to the Moon, Mars, outer planets, and beyond, including the Voyager, Perseverance, and James Webb Space Telescope.

---

![Screenshot 1](_assets/_shot1.png)

## Features

- **Live monitoring with real-time, detailed data.** Data updates are almost instantaneous.
- **Display targets being tracked by the stations.** You can visualize the current target, distance, and round-trip light time value.
- **Display up and down signals for each station.** You can observe the signal strength, frequency, data rate, and modulation type.
- **Interactive.** You can select a station, antenna, current target, and both types of signals.
- **Informative.** You can access detailed information about the station, antenna, target, and signals.
- **Themes.** Cycle through 4 built-in color themes: dark, solarized, nord, and light.
- **Bookmarks.** Mark dishes of interest with a star indicator. Bookmarks persist across sessions.
- **Compact view.** Toggle between a detailed view and a compact table overview of all dishes.
- **Signal notifications.** Get notified in the status bar when signals are acquired or lost.
- **Connection indicator.** A colored status indicator shows whether the data feed is connected, degraded, or disconnected.

## Keyboard Shortcuts

| Key     | Action                          |
| ------- | ------------------------------- |
| `s`     | Cycle through stations          |
| `t`     | Cycle through targets           |
| `u`     | Cycle through up signals        |
| `d`     | Cycle through down signals      |
| `↑` `↓` | Navigate dishes list            |
| `b`     | Toggle bookmark on current dish |
| `c`     | Toggle compact view             |
| `T`     | Cycle through themes            |
| `j`     | Show JSON preview               |
| `i`     | Show antenna specifications     |
| `?`     | Show help                       |
| `Esc`   | Close modal                     |
| `q`     | Quit                            |

## Installation

0. Requires [Go](https://go.dev/) to be installed.
1. Install the application by running: `go install github.com/RustyDaemon/go-dsn-now@latest`
2. Run the application with `go-dsn-now`.

> **Note:** If you get `command not found`, make sure `$GOPATH/bin` is in your `PATH`. You can find your `GOPATH` by running `go env GOPATH`.

## License

Licensed under [Apache 2.0](LICENSE). Third-party dependencies are subject to their own licenses.

[NASA Deep Space Network Now](https://eyes.nasa.gov/dsn/dsn.html) data is sourced indirectly from NASA and used solely for educational and informational purposes. This project is not affiliated with NASA or any of its subsidiaries.
