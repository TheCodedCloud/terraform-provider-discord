{
  description = "terraform-provider-discord - development environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }: let
    supported = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
    forAllSystems = nixpkgs.lib.genAttrs supported;
    pkgsFor = system: import nixpkgs { inherit system; };
  in {
    devShells = forAllSystems (system: let
      pkgs = pkgsFor system;
    in {
      default = pkgs.mkShell {
        name = "terraform-provider-discord-dev";

        packages = [
          pkgs.pre-commit
          pkgs.opentofu 
          pkgs.tflint
        ];

        shellHook = ''
          echo "Tofu:       $(tofu version)"
          echo "Pre-commit: $(pre-commit --version)"
          echo ""
          if [ -d .git ]; then
            pre-commit install 2>/dev/null || true
            # Run pre-commit inside nix develop so Cursor's git client (which
            # doesn't load this dev shell) still has tofu/tflint on PATH.
            root="$(git rev-parse --show-toplevel)"
            printf '%s\n' '#!/bin/sh' 'cd "$(git rev-parse --show-toplevel)" && nix develop -c pre-commit run --hook-stage pre-commit' > "$root/.git/hooks/pre-commit"
            chmod +x "$root/.git/hooks/pre-commit"
          fi
        '';
      };
    });
  };
}
