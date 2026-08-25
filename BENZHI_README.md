# apriori-miner — Go 关联规则挖掘与频繁项集分析 HTTP 服务（含持久化）

Frequent itemset mining and association rule generation service using the
Apriori algorithm. Supports configurable minimum support/confidence thresholds,
rule ranking, sampling, and persistent itemset storage.

## Build / Run / Test

```bash
go build -o apriori-miner .
./apriori-miner serve -addr :8080
./apriori-miner -input example/transactions.txt -min-support 0.3
go test ./...
```

## Evaluation Image

Evaluation-specific files (do not overwrite project Dockerfile/README):

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md` (this file)

Build and verify in container:

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
