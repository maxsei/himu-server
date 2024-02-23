{ stdenv, lib, fetchzip, autoPatchelfHook }:

assert builtins.currentSystem == "x86_64-linux";

stdenv.mkDerivation rec {
  name = "duckdb";
  version = "0.10.0";
  src = fetchzip {
    url = "https://github.com/duckdb/duckdb/releases/download/v${version}/duckdb_cli-linux-amd64.zip";
    sha256 = "sha256-A2U7GFr3qRQibvBOOhFvMtrakB+EkXttRqFzpVeaViU=";
  };

  nativeBuildInputs = [ autoPatchelfHook ];

  buildInputs = [ stdenv.cc.cc.lib ];

  buildPhase = "";

  installPhase = ''
    runHook preInstall
    install -m755 -D ${name} $out/bin/${name}
    runHook postInstall
  '';
}
