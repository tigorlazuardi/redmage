{
  inputs = {
    templ.url = "github:a-h/templ/v0.2.542"; # 0.2.542
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
  };

  outputs = inputs@{ templ, nixpkgs, ... }:
    let
      system = "x86_64-linux";
      templPkg = templ.packages.${system}.templ;
      pkgs = inputs.nixpkgs.legacyPackages.${system};
    in
    {
      devShell.${system} = pkgs.mkShell rec {
        name = "redmage-shell";
        buildInputs = with pkgs; [
          templPkg
          go
          modd
          nodejs_21
        ];
      };
    };
}
