{
  inputs = {
    templ.url = "github:a-h/templ/v0.2.663";
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
  };

  outputs = inputs@{ templ, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      templPkg = templ.packages.${system}.templ;
      pkgs = inputs.nixpkgs.legacyPackages.${system};
      goverter = pkgs.buildGoModule rec {
        name = "goverter";
        version = "1.5.0";
        src = pkgs.fetchFromGitHub {
          owner = "jmattheis";
          repo = "goverter";
          rev = "v${version}";
          sha256 = "sha256-J0PS4ZxGtOa+0QOOSjfg0WeVYGyf757WuTnpQTWIV1w=";
        };
        nativeBuildInputs = [ pkgs.go ];
        buildPhase = ''
          runHook preBuild
          go build -o goverter ./cmd/goverter
          runHook postBuild
        '';
        installPhase = ''
          runHook preInstall
          mkdir -p $out/bin
          cp goverter $out/bin
          runHook postInstall
        '';
        vendorHash = "sha256-uQ1qKZLRwsgXKqSAERSqf+1cYKp6MTeVbfGs+qcdakE=";
      };
    in
    {
      devShell.${system} = pkgs.mkShell {
        name = "redmage-shell";
        buildInputs = with pkgs; [
          templPkg
          go
          modd
          nodePackages_latest.nodejs
          goose
          air
          upx
          buf
          buf-language-server
          protoc-gen-go
          protoc-gen-go-grpc
          protoc-gen-connect-go
          protoc-gen-validate
          goverter
        ];
      };
    };
}
