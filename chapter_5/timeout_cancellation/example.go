package main

func main() {

}

func preemptableFunc(done <-chan any, valueStream <-chan any) {
	resultStream := make(chan any)

	longCalculation := func(value any) any {
		// do some long calculation
		return nil
	}
	betterLongCalculation := func(done <-chan any, value any) any {
		// do some long calculation
		return nil
	}
	// almost preemptable function
	reallyLongCalculation := func(done <-chan any, value any) any {
		// long calculation not preemptable
		intermediateResult := longCalculation(value)
		select {
		case <-done:
			return nil
		default:
		}
		return longCalculation(intermediateResult)
	}
	betterReallyLongCalculation := func(done <-chan any, value any) any {
		// long calculation not preemptable
		intermediateResult := betterLongCalculation(done, value)
		return betterLongCalculation(done, intermediateResult)
	}
	go func() {
		defer close(resultStream)

		var value any
		select {
		case <-done:
			return
		case value = <-valueStream:
		}
		// almost preemptable
		result := reallyLongCalculation(done, value)

		// preemptable
		result = betterReallyLongCalculation(done, value)
		select {
		case <-done:
			return
		case resultStream <- result:
		}
	}()
}

func notPreemptableFunc(done <-chan any, valueStream <-chan any) {
	resultStream := make(chan any)

	// non preemptable function
	reallyLongCalculation := func(value any) any { return nil }

	go func() {
		defer close(resultStream)

		var value any
		select {
		case <-done:
			return
		case value = <-valueStream:
		}
		result := reallyLongCalculation(value)

		select {
		case <-done:
			return
		case resultStream <- result:
		}
	}()
}
