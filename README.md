# YLX

**Your Local Exchange.**

YLX is a local classifieds platform where buyers and sellers discover each
other, communicate directly, and arrange payment and delivery themselves.

> YLX is in early development.

## Getting started

### Requirements

- Go 1.26.5+
- Node.js 18+
- pnpm 9

### Run the API

```sh
cd apps/api
make run
```

The API runs at `http://localhost:8080`. Check it with:

```sh
curl http://localhost:8080/healthz
```

## Workspace

This repository is a Turborepo monorepo managed with pnpm. Application code
lives in `apps/`, and shared packages will live in `packages/`.
