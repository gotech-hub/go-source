# 🧪 Go Project - Test Coverage & Build Guide

This repository provides a Go project setup with test automation, coverage validation, mock generation, and a handy `Makefile` for common development tasks.

---

## 📁 Requirements

Before running the commands, make sure the following tools are installed:

- **Go** (1.22+ recommended)
- [`mockgen`](https://github.com/golang/mock) for interface mocking
- [`dlv`](https://github.com/go-delve/delve) for debugging (optional)
- Unix tools: `bash`, `awk`, `grep`, `bc`

Install `mockgen` if you haven't already:

```bash
go install github.com/golang/mock/mockgen@latest
```

Run check unit tests is recommended to be run in a containerized environment.
```bash
./coverage.sh
```