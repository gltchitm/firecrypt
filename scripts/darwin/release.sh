#!/usr/bin/env sh

set -e

cd $(dirname "$0")/../..

if [ -d release/darwin ]; then
    rm -rf release/darwin
fi

if ! [ -f icon/darwin/icon.icns ]; then
    cd icon/darwin

    if [ -d Firecrypt.iconset ]; then
        rm -rf Firecrypt.iconset
    fi

    mkdir Firecrypt.iconset
    cd Firecrypt.iconset

    sips -z 16 16 ../icon.png --out icon_16x16.png
    sips -z 32 32 ../icon.png --out icon_16x162x.png
    sips -z 32 32 ../icon.png --out icon_32x32.png
    sips -z 64 64 ../icon.png --out icon_32x322x.png
    sips -z 64 64 ../icon.png --out icon_64x64.png
    sips -z 128 128 ../icon.png --out icon_64x642x.png
    sips -z 128 128 ../icon.png --out icon_128x128.png
    sips -z 256 256 ../icon.png --out icon_128x1282x.png
    sips -z 256 256 ../icon.png --out icon_256x256.png
    sips -z 512 512 ../icon.png --out icon_256x2562x.png
    sips -z 512 512 ../icon.png --out icon_512x512.png
    sips -z 1024 1024 ../icon.png --out icon_512x5122x.png
    sips -z 1024 1024 ../icon.png --out icon_1024.png

    cd ..

    iconutil -c icns Firecrypt.iconset -o icon.icns

    rm -rf Firecrypt.iconset

    cd ../..
fi

mkdir -p release/darwin/Firecrypt.app/Contents
mkdir release/darwin/Firecrypt.app/Contents/MacOS
mkdir release/darwin/Firecrypt.app/Contents/Resources
mkdir release/darwin/Firecrypt.app/Contents/Resources/icon

set +e
version=$(git describe --tags --abbrev=0 2>/dev/null)
set -e

if [[ $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    version=${version:1}
else
    version=0.0.0
fi

cat << EOF > release/darwin/Firecrypt.app/Contents/Info.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>
    <string>Firecrypt</string>
    <key>CFBundleDisplayName</key>
    <string>Firecrypt</string>
    <key>CFBundleExecutable</key>
    <string>firecrypt</string>
    <key>CFBundleIconFile</key>
    <string>icon/icon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.github.gltchitm.Firecrypt</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleVersion</key>
    <string>$version</string>
    <key>CFBundleShortVersionString</key>
    <string>$version</string>
</dict>
</plist>
EOF

cp NOTICE release/darwin/Firecrypt.app/Contents/Resources/NOTICE
cp icon/darwin/icon.icns release/darwin/Firecrypt.app/Contents/Resources/icon/icon.icns
cp icon/NOTICE release/darwin/Firecrypt.app/Contents/Resources/icon/NOTICE
cp -r ui/* release/darwin/Firecrypt.app/Contents/Resources

go clean --cache
CGO_CFLAGS="-DFIRECRYPT_VERSION=@\"$version\"" go build -tags release -o release/darwin/Firecrypt.app/Contents/MacOS/firecrypt .

cd release/darwin

zip -r Firecrypt.zip Firecrypt.app

cd ../..
