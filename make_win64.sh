rm build_win64/*.exe
cp -rf assets/ build_win64/
cp -rf missions/ build_win64/

build_datetime=`date '+%Y_%m_%d__%H_%M_%S'`;
exename="TaffeRL_"$build_datetime".exe"

echo "Building "$exename"..."

env CGO_ENABLED="1" CC="/usr/bin/x86_64-w64-mingw32-gcc" GOOS="windows" CGO_LDFLAGS="-lmingw32 -lSDL2" CGO_CFLAGS="-D_REENTRANT" go build -o build_win64/$exename *.go

