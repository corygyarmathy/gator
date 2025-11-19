# gator

boot.dev project written in Go.

## Development Setup

### Nix

1. Install `direnv` with `nix-direnv`:

```nix
  programs.direnv = {
    enable = true;
    nix-direnv.enable = true;
  };
```

2. Enter the project directory:

```bash
  cd my-project
  direnv allow
```

A `nix develop` environment will automatically be created from the `flake.nix` and `flake.lock` files.

Manual option: run `nix develop` each time you navigate to the project directory in a shell or editor.
