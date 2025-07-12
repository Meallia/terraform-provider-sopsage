{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { nixpkgs, flake-utils, ... }: flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; config.allowUnfree = true; };
    in {
    formatter = pkgs.alejandra;
    devShell = pkgs.mkShell {
      buildInputs = with pkgs; [
        go
        terraform
        gnumake
		golangci-lint
      ];
    };
  });
}
