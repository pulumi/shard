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

## Why?

Strategies like `go test ./... -list .` compile your code in order to discover test cases.
This can be especially slow in CI environments depending on the state of your build cache.

We can discover test cases significantly faster by essentially `grep`-ing for patterns like `Test..`, `Fuzz...`, etc.

## Caveats

* A test can potentially be executed more than once if another package shares a test with the same name.
  Renaming your tests to be globally unique is currently the best workaround if you want to guarantee a single execution per test function.
  You can discover test with name collisions by running `shard --total 1 --index 0`.
* Benchmarks aren't currently collected so running with `-bench` will not have any effect.
      


