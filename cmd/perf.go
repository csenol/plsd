package cmd

import (
	"bufio"
	"fmt"
	. "github.com/csenol/plsd/pkg/api"
	"github.com/jamiealquiza/tachymeter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"time"
)

var queryFile string
var contextFile1 string
var index1 string
var from int
var size int
var timeout string
var terminateAfter int
var repeat int
var debug bool

// perfCmd represents the perf command
var perfCmd = &cobra.Command{
	Use:   "perf",
	Short: "This command can be used to help understanding impacts of sorting with painless script. This feature is experimental. Output/methodology might change",
	Long: `This commands takes a query, a script to sort and context file same as in exec command for parameters. 
It runs the query number-of-times aggregates the response times and prints latency distribution of all collectors. This feature is experimental. Output/methodology might change
`,
	Run: func(cmd *cobra.Command, args []string) {
		var script = ""
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			script = script + "\n" + scanner.Text()
		}
		resp, _ := RunPerf(queryFile, script, contextFile1, esEndpoint, index1, from, size, timeout, terminateAfter, repeat, debug)
		t := tachymeter.New(&tachymeter.Config{Size: repeat * 100})
		for _, timing := range resp {
			s := strconv.FormatInt(timing, 10) + "ns"
			d, _ := time.ParseDuration(s)
			t.AddTime(d)
		}

		fmt.Println(t.Calc().String())

	},
}

func init() {
	rootCmd.AddCommand(perfCmd)
	perfCmd.Flags().StringVar(&contextFile1, "context-file", "", "Json File with Document and Params Field. If this is left script will run with an empty document and empty params. {\"index\": \"Index Name\", \"params\":{}, \"document\": {}, \"context\": \"score|filter\"}")
	perfCmd.Flags().StringVar(&queryFile, "query-file", "", "ES Query to run performance with script")
	perfCmd.Flags().StringVar(&index1, "index", "", "Index for running scripts in case file param is missed")
	perfCmd.Flags().IntVar(&from, "from", 0, "Starting offset of ES Query")
	perfCmd.Flags().IntVar(&size, "size", 100, "Number of documents to retrieve from ES")
	perfCmd.Flags().StringVar(&timeout, "timeout", "300ms", "Timeout to complete ES Query")
	perfCmd.Flags().IntVar(&terminateAfter, "terminate-after", 100000, "ES Terminate After parameter")
	perfCmd.Flags().IntVar(&repeat, "repeat", 10, "Number of Repetition for Performance Test")
	perfCmd.Flags().BoolVar(&debug, "debug", false, "If set, logs first ES Request and Response")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// perfCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// perfCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
