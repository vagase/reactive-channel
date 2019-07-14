# reactive-channel
Reactive programming utilities for go, implmenting with go channels. Inspired by [ReactiveX](http://reactivex.io).

# install
go get -u github.com/vagase/reactive-channel

# example
### Distinct
```go
inChan := From([]interface{}{1, 2, 2, 3, 3, 5})
outChan := Distinct(inChan)
```
*outChan will emit: 1, 2, 3, 5*
### Filter
```go
inChan := From([]interface{}{1, 2, 3, 4})
outChan := Filter(inChan, func(i interface{}) bool {
  num := i.(int)
  return num%2 == 0
})
```
*outChan will emit: 2 4*
### CombineLatest
```go
// 1, 2, 3, 4
inChan1 := IntervalRange(1, 4, time.Millisecond*30, 0)
// a, b
inChan2 := Map(IntervalRange(0, 2, time.Millisecond*50, 0), func(i interface{}) interface{} {
  return string(97 + i.(int))
})
outChan := Map(CombineLatest(inChan1, inChan2), func(i interface{}) interface{} {
  values := i.([]interface{})
  return strconv.Itoa(values[0].(int)) + values[1].(string)
})
```
*outChan will emit: "1a", "2a", "3a", "3b", "4b"*
### Debounce
```go
inChan := IntervalRange(0, 4, 40*time.Millisecond, 0)
outChan := Debounce(inChan, time.Millisecond*100)
```
*outChan will emit: 0, 3*
# usage 
reactive-channel implements most useful ReactiveX operators, checkout [full api list here](http://reactivex.io/documentation/operators.html#categorized).
* [From](http://reactivex.io/documentation/operators/from.html)
* [Interval](http://reactivex.io/documentation/operators/interval.html)
* [IntervalRange](http://reactivex.io/documentation/operators/range.html)
* [Range](http://reactivex.io/documentation/operators/range.html)
* [Repeat](http://reactivex.io/documentation/operators/repeat.html)
* [Just](http://reactivex.io/documentation/operators/just.html)
* [Timer](http://reactivex.io/documentation/operators/timer.html)
* [Map](http://reactivex.io/documentation/operators/map.html)
* [Reduce](http://reactivex.io/documentation/operators/reduce.html)
* [Filter](http://reactivex.io/documentation/operators/filter.html)
* [Merge](http://reactivex.io/documentation/operators/merge.html)
* [Subscribe](http://reactivex.io/documentation/operators/subscribe.html)
* [Buffer](http://reactivex.io/documentation/operators/buffer.html)
* [FlatMap](http://reactivex.io/documentation/operators/flatmap.html)
* [GroupBy](http://reactivex.io/documentation/operators/groupby.html)
* [Debounce](http://reactivex.io/documentation/operators/debounce.html)
* [Distinct](http://reactivex.io/documentation/operators/distinct.html)
* [ElementAt](http://reactivex.io/documentation/operators/elementat.html)
* [First](http://reactivex.io/documentation/operators/first.html)
* [Last](http://reactivex.io/documentation/operators/last.html)
* [IgnoreElements](http://reactivex.io/documentation/operators/ignoreelements.html)
* [Sample](http://reactivex.io/documentation/operators/sample.html)
* [Skip](http://reactivex.io/documentation/operators/skip.html)
* [SkipLast](http://reactivex.io/documentation/operators/skiplast.html)
* [Take](http://reactivex.io/documentation/operators/take.html)
* [TakeLast](http://reactivex.io/documentation/operators/takelast.html)
* [CombineLatest](http://reactivex.io/documentation/operators/combinelatest.html)
* [StartWith](http://reactivex.io/documentation/operators/startwith.html)
* [Switch](http://reactivex.io/documentation/operators/switch.html)
* [Zip](http://reactivex.io/documentation/operators/zip.html)
* [Delay](http://reactivex.io/documentation/operators/delay.html)
* [TimeInterval](http://reactivex.io/documentation/operators/timeinterval.html)
* [Timeout](http://reactivex.io/documentation/operators/timeout.html)
* [Timestamp](http://reactivex.io/documentation/operators/timestamp.html)
* [All](http://reactivex.io/documentation/operators/all.html)
* [Contains](http://reactivex.io/documentation/operators/contains.html)
* [Amb](http://reactivex.io/documentation/operators/amb.html)
* [DefaultIfEmpty](http://reactivex.io/documentation/operators/defaultifempty.html)
* [SequenceEqual](http://reactivex.io/documentation/operators/sequenceequal.html)
* [SkipUntil](http://reactivex.io/documentation/operators/skipuntil.html)
* [SkipWhile](http://reactivex.io/documentation/operators/skipwhile.html)
* [TakeUntil](http://reactivex.io/documentation/operators/takeuntil.html)
* [TakeWhile](http://reactivex.io/documentation/operators/takewhile.html)
* [Sum](http://reactivex.io/documentation/operators/sum.html)
* [Average](http://reactivex.io/documentation/operators/average.html)
* [Count](http://reactivex.io/documentation/operators/count.html)
* [Min](http://reactivex.io/documentation/operators/min.html)
* [Max](http://reactivex.io/documentation/operators/max.html)
* [Concat](http://reactivex.io/documentation/operators/concat.html)
* [IsEmpty](http://reactivex.io/documentation/operators/contains.html)
* [WindowWithCount](http://reactivex.io/documentation/operators/window.html)
# license
MIT Licence © [vagase](https://github.com/vagase)
