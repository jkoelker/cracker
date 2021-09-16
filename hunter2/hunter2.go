//

package hunter2

import  (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/schollz/progressbar/v3"

	"github.com/jkoelker/cracker/hash"
)

func lines(done <-chan struct{}, path string) (<-chan string, <-chan error) {
	lines := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		defer close(lines)

		file, err := os.Open(path)
		if err != nil {
			errc <- err
			return
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			select {
			case lines <- scanner.Text():
			case <-done:
				break
			}
		}

		if err := scanner.Err(); err != nil {
			errc <- err
			return
		}
	}()

	return lines, errc
}

func digester(
	done <-chan struct{},
	passwords <-chan string,
	result chan <- string,
	progress *progressbar.ProgressBar,
	target string,
) {
	hasher := hash.NewBCrypt()

	for password := range passwords {
		select {
		case <-done:
			return
		default:
			progress.Add(1)
			if hasher.Check(password, target){
				result <- password
			}
		}
	}
}

func Search(target string, list string, workers int) error {
	fmt.Printf("Searching %s for hash %s\n", list, target)
	done := make(chan struct{})
	defer close(done)

	lines, linesErr := lines(done, list)
	result := make(chan string)

	var wg sync.WaitGroup
	wg.Add(workers)
	bar := progressbar.Default(-1)

	for i := 0; i < workers; i++ {
		go func() {
			digester(done, lines, result, bar, target)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(result) // HLc
	}()

	for r := range result {
		fmt.Printf("\n********\nFound password: %s\n********\n", r)
	}

	if err := <-linesErr; err != nil {
		return err
	}

	fmt.Println("")

	return nil
}
