# Shard 🔪

A quick hack to shard large Go test suites.

Example usage with GitHub Actions:

```
jobs:
  tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        shard: [0, 1, 2, 3]
  steps:
    - name: Shard tests
      run: go run github.com/blampe/shard --total ${{ strategy.job-total }} --index ${{ strategy.job-index }} > tests
    - name: Run tests
      run: go test $(cat tests)
```
      


