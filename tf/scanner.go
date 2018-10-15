package tf

type MatchFunc func(s string) bool

// Scanner watches output from io.Writers for known lines
type Scanner struct {
	Matcher  MatchFunc
	Response chan string
	Error    error
}

func NewScanner(m MatchFunc) *Scanner {
	return &Scanner{
		Matcher:  m,
		Response: make(chan string, 1000),
	}
}

// // Watch iterates all lines passed through the buffer, sending matches down the channel
// func (s *Scanner) Watch(wg *sync.WaitGroup, scanner *bufio.Scanner) {
// 	defer wg.Done()
// 	for scanner.Scan() {
// 		text := scanner.Text()
// 		if s.Matcher(text) {
// 			s.Response <- text
// 			wg.Add(1)
// 		}
// 	}
// }

// // HookOutput taps into the stream's output
// func (s *Scanner) HookOutput(stream io.ReadCloser, wg *sync.WaitGroup) {
// 	linestream := bufio.NewScanner(stream)
// 	go s.Watch(wg, linestream)
// 	linestream.Err()
// 	return
// }
