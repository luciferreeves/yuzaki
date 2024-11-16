[![Yuzaki.gif](https://i.postimg.cc/VkxMN0jP/Yuzaki.gif)](https://postimg.cc/PCMN3x3S)

# Yuzaki — Discord Bot

> _Note_: This bot is still in development. Only build instructions are provided for now.

## Build Instructions

### Prerequisites

- [Go](https://go.dev) (v1.23 or higher)

After forking the repository, clone it to your local machine. Commands are run via the [`Makefile`](Makefile).

Setup your environment by running the following command:

```bash
make setup
```

Run the bot in development mode:

```bash
make dev
```

Build the bot:

```bash
make build
```

Run the bot (this will build the bot if it hasn't been built yet):

```bash
make run
```
