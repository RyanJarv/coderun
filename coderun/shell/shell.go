package shell

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/chzyer/readline"
)

func dir(d string) []string {
	list := []string{}
	files, err := ioutil.ReadDir(d)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			list = append(list, "./"+path.Join(d, f.Name())+"/")
			list = append(list, "./"+path.Join(d, f.Name())+"/..")
		} else {
			list = append(list, "./"+path.Join(d, f.Name()))
		}
	}
	return list
}
func listDirs(search string) func(string) []string {
	list := dir(search)

	return func(line string) []string {
		if f, err := os.Stat(path.Dir(line)); err == nil && f.IsDir() {
			return dir(path.Dir(line))
		} else {
			return list
		}
	}
}

func listFiles(search string) func(string) []string {
	return func(list string) []string { return complete("-f", search) }
}

func listCmds(search string) func(string) []string {
	return func(list string) []string { return complete("-c", search) }
}

func complete(action, search string) []string {
	c := exec.Command("bash", "-l", "-c", fmt.Sprintf("compgen %s '%s'", action, search))
	out, err := c.CombinedOutput()
	if err != nil {
		log.Printf(string(out))
		log.Fatal(err)
	}
	return strings.Split(string(out), "\n")
}

func NewShell() *Shell {
	p := &Shell{
		completer: readline.NewPrefixCompleter(
			readline.PcItem("mode",
				readline.PcItem("vi"),
				readline.PcItem("emacs"),
			),
			readline.PcItem("setprompt"),
			readline.PcItem("bye"),
			readline.PcItem("help"),
			readline.PcItemDynamic(listDirs("./")),
			//readline.PcItemDynamic(listCmds("")),
		),
	}
	return p
}

type Shell struct {
	completer *readline.PrefixCompleter
	instance  *readline.Instance
}

func (p *Shell) Instance() *readline.Instance { return p.instance }

func (p *Shell) AddCompleters(add *readline.PrefixCompleter) {
	c := p.completer.GetChildren()
	p.completer.SetChildren(append(c, add))
}

func (p *Shell) usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, p.completer.Tree("    "))
}

func (p *Shell) filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func (p *Shell) Teardown() {
	p.instance.Close()
}

func (p *Shell) Start(callback func(line string)) {
	var err error
	p.instance, err = readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    p.completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: p.filterInput,
	})
	if err != nil {
		panic(err)
	}

	log.SetOutput(p.instance.Stderr())
	for {
		line, err := p.instance.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "mode "):
			switch line[5:] {
			case "vi":
				p.instance.SetVimMode(true)
			case "emacs":
				p.instance.SetVimMode(false)
			default:
				println("invalid mode:", line[5:])
			}
		case line == "mode":
			if p.instance.IsVimMode() {
				println("current mode: vim")
			} else {
				println("current mode: emacs")
			}
		case line == "help":
			p.usage(p.instance.Stderr())
		case strings.HasPrefix(line, "setprompt"):
			if len(line) <= 10 {
				log.Println("setprompt <prompt>")
				break
			}
			p.instance.SetPrompt(line[10:])
		case line == "bye":
			goto exit
		case line == "":
		default:
			callback(line)
		}
	}
exit:
}
