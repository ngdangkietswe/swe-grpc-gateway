echo "Update dependency..."
GOPROXY=direct go get -u github.com/ngdangkietswe/swe-protobuf-shared
GOPROXY=direct go get -u github.com/ngdangkietswe/swe-go-common-shared
go mod tidy
go mod vendor
echo "Update dependency successful!"

echo "Update openapi..."
git clone https://github.com/ngdangkietswe/swe-protobuf-shared.git
cp -r swe-protobuf-shared/openapiv2/* swagger/
rm -rf swe-protobuf-shared
echo "Update openapi successful!"