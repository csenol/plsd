package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	. "github.com/csenol/plsd/pkg/api"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var index string
var scriptFile string
var watch bool
var contextFile string

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Executes Painless Script and returns output",
	Long: `Executes a given script against a running Elasticsearch Node.
`,
	Run: func(cmd *cobra.Command, args []string) {
		watcher, _ := fsnotify.NewWatcher()

		var tcs *TestCaseSetup
		var err error
		if contextFile == "" {
			if index == "" {
				os.Stderr.WriteString("index param is needed when context file is not specified \n")
				os.Exit(-1)
			}
			tcs = &TestCaseSetup{Context: "score", Index: index, Params: make(map[string]interface{}), Document: make(map[string]interface{})}
		} else {
			tcs, err = LoadTestCaseSetup(contextFile)
			if err != nil {
				fmt.Print(err)
				os.Exit(-1)
				return
			}

		}

		var script = ""
		if scriptFile == "" {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				script = script + "\n" + scanner.Text()
			}

		} else {
			b, err := ioutil.ReadFile(scriptFile)
			if err != nil {
				os.Stderr.WriteString(err.Error())
			}
			script = string(b)
		}

		resp, err := ExecuteQuery(script, *tcs, esEndpoint)

		if err != nil {
			os.Stderr.WriteString(err.Error())
		} else {
			fmt.Print(resp.Result)
		}

		if watch {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			watcher.Add(scriptFile)
			watcher.Add(contextFile)
			for {
				select {
				case event, ok := <-watcher.Events:
					fmt.Println(fmt.Sprintf("%s changed", event.Name))
					if !ok {
						return
					}
					tcs, err := LoadTestCaseSetup(contextFile)
					if err != nil {
						os.Stderr.WriteString(err.Error())
					}
					b, err := ioutil.ReadFile(scriptFile)
					if err != nil {
						os.Stderr.WriteString(err.Error())
					}
					script := string(b)
					resp, err := ExecuteQuery(script, *tcs, esEndpoint)

					if err != nil {
						os.Stderr.WriteString(err.Error())
					} else {
						fmt.Println(resp.Result)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						os.Stderr.WriteString(err.Error())
						return
					}

				case <-sigs:
					os.Exit(-1)
				}
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVar(&contextFile, "context-file", "", "Json File with Document and Params Field. If this is left script will run with an empty document and empty params. {\"index\": \"Index Name\", \"params\":{}, \"document\": {}, \"context\": \"score|filter\"}")
	execCmd.Flags().StringVar(&index, "index", "", "Index for running scripts in case file param is missed")
	execCmd.Flags().StringVar(&scriptFile, "script-file", "", "Painless script file to run. If not given, stdin will be used")
	execCmd.Flags().BoolVar(&watch, "watch", false, "Set flag if you want to watch script and paramater file and rerun plse on every change")
}
