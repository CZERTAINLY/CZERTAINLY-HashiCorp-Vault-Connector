# export GO_POST_PROCESS_FILE="gofmt -w"

# openapi-generator-cli generate -g go-server -o generated --enable-post-process-file \
#  --additional-properties=packageName=authority,sourceFolder=authority,outputAsLibrary=true \
#  -i ./api/connector-api/doc-openapi-connector-authority-provider-v2.yaml

 openapi-generator-cli generate -g go-server -o generated --enable-post-process-file --additional-properties=packageName=discovery,sourceFolder=discovery,outputAsLibrary=true -i ./api/connector-api/doc-openapi-connector-discovery-provider-disc.yaml


# #  openapi-generator-cli generate -g go-server -o test --enable-post-process-file \
# #  --additional-properties=packageName=authority,sourceFolder=authority \
# #  -i ./api/connector-api/doc-openapi-connector-authority-provider-v2.yaml

# #  openapi-generator-cli generate -g go-server -o test --enable-post-process-file \
# #  --additional-properties=packageName=discovery,sourceFolder=discovery \
# #  -i ./api/connector-api/doc-openapi-connector-discovery-provider.yaml

# ~/swagger mixin api/connector-api/doc-openapi-connector-discovery-provider.yaml api/connector-api/doc-openapi-connector-authority-provider-v2.yaml -o api/connector-api/merged.yaml --format=yaml

# sed -i '1i openapi: 3.0.1'  api/connector-api/merged.yaml

openapi-generator-cli generate -g go-server -o generated --enable-post-process-file \
 --additional-properties=packageName=openapi,sourceFolder=test \
 -i ./api/connector-api/merged.yaml


 