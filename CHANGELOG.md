# Changelog

## 1.0.0

First stable release of Go DSN Now.

- Search stations, dishes, and spacecraft with live results.
- Combine active-only and bookmarked-only filters in both views.
- Keep compact selection consistent through sorting, refreshes, bookmarks, and opening details.
- Improve rendering performance by skipping unused detail rendering.
- Show startup errors and retry failed requests without restarting.
- Prevent overlapping refresh requests and keep long errors within the status bar.
- Add regression tests for navigation, filtering, startup recovery, and terminal layout.
