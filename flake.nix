{
  description = "Go development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
      {
        devShells.default = pkgs.mkShell {
          name = "go-dev-shell";

          packages = with pkgs; [
            # Go toolchain
            go

            # Go development tools
            # gopls # Language server
            golangci-lint # Linter aggregator
            gotools # Includes goimports, godoc, etc.
            delve # Debugger

            # Optional but useful
            # air # Live reload for Go apps
            # go-migrate # Database migrations
            # sqlc # SQL code generator

            # Project dependencies
            postgresql
          ];

          shellHook = ''
            echo "Go $(go version | awk '{print $3}') Nix development environment"

            # Set up project-specific environment variables
            # export DATABASE_URL="postgresql://localhost/mydb"
            # export GO_ENV="development"
          '';
        };
      }
    );
}
