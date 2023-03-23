# Written answers

## Task 1

##### Q:
What happens if you remove the `go-command` from the `Seek` call in the `main` function?

##### A:
Anna will always message Bob, Cody will always message Dave and Eva will not send a message to anyone.
This happens because only the main goroutine runs which means that the names will be used in the order of the `people` array

##### Q:
What happens if you switch the declaration `wg := new(sync.WaitGroup)` to `var wg sync.WaitGroup` and the parameter `wg *sync.WaitGroup` to `wg sync.WaitGroup`?

##### A:
`wg`would then be a WaitGroup object instead of a pointer to it. This means that `wg.Done()`only affects
the local copy of it which makes results in a deadlock. Declaring `wg` as an object works if the `Seek`function looks like this: `Seek(name, match, &wg)` and the parameter remains `wg *sync.WaitGroup`

##### Q:
What happens if you remove the buffer on the channel match?

##### A:
Making the channel unbuffered casues a deadlock. The last name sent to the channel will be blocked
as there isn't any goroutine that can receive it. The buffer allows the channel to receive 1 name without blocking it.

##### Q:
What happens if you remove the default-case from the case-statement in the `main` function?

##### A:
If the array of names instead was even in length, it would be a deadlock. This happens because the select statement in main will wait for data from the channel match (this data will never come since everyone has received a message).

