package cmd

import (
	"bufio"
	"fmt"
	"os"

	. "github.com/csenol/plsd/pkg/api"
	"github.com/spf13/cobra"
)

var parallel bool
var testFile string

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs test files with the script",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if testFile == "" {
			fmt.Println("test-file is not set")
			os.Exit(-1)
			return
		}

		var script = ""
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			script = script + "\n" + scanner.Text()
		}

		tcs, err := LoadTestCaseSetups(testFile)
		if err != nil {
			fmt.Print(err)
			os.Exit(-1)
			return
		}
		testCase := TestCase{testFile, script, tcs}

		failed := RunTestCase(testCase, esEndpoint, parallel)

		if failed != 0 {
			os.Exit(int(failed))
		}
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringVar(&testFile, "test-file", "", "Json File Path with test setup and expected value [{\"description\": \"Test Description\", \"index\": \"Index Name\", \"params\":{}, \"document\": {}, \"expected_result\": \"Number, String or Boolean\", \"context\": \"score|filter\"} ]")
	testCmd.Flags().BoolVar(&parallel, "parallel", false, "Set the flag to run tests on parallel")

}
