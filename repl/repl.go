package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/NishanthSpShetty/monkey/lexer"
	"github.com/NishanthSpShetty/monkey/parser"
	"github.com/NishanthSpShetty/monkey/runtime/evaluator"
	"github.com/NishanthSpShetty/monkey/runtime/evaluator/runtime"
)

const (
	PROMPT      = ">>"
	MONKEY_FACE = `           __,__
  .--.  .-"    "-. .   --.
 / .. \/ .-.  .-.   \/ .. \
| |  '| /    Y    \  |'  | |
| \   \ \  0 | 0  / /   /  |
 \ '- ,\.-"""""""-./, -'  /
  ''-' /_   ^ ^   _\ '-''
      |  \._   _./  |
      \   \ '~' /   /
       '._ '-=-' _.'
          '-----'
`
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	r := runtime.New()
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			// no more tokens to read
			return
		}

		line := scanner.Text()

		if strings.HasPrefix(line, ":") {
			// runtime/repl commands
			run(r, out, line)
			continue
		}

		l := lexer.New(line)

		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Erors()) != 0 {
			printParseError(out, p.Erors())
			continue
		}

		eval := evaluator.Eval(r, program)
		if eval != nil {

			io.WriteString(out, eval.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseError(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func run(r *runtime.Runtime, out io.Writer, line string) {
	cmd := strings.Split(line, ":")[1]
	switch cmd {
	case "env":
		r.PrintVars()
	default:

		io.WriteString(out, MONKEY_FACE)
		io.WriteString(out, "Woops! We ran into some monkey business here!\n")
		io.WriteString(out, fmt.Sprintf("invalid command: %s\n", cmd))
	}
}
