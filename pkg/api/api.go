package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
)

const (
	delta      = 0.00001
	InfoColor  = "\033[1;32m%s\033[0m"
	ErrorColor = "\033[1;31m%s\033[0m"
)

type TestCaseSetup struct {
	Index          string                 `json:"index"`
	Params         map[string]interface{} `json:"params"`
	Document       map[string]interface{} `json:"document"`
	ExpectedResult interface{}            `json:"expected_result"`
	Description    string                 `json:"description"`
	Context        string                 `json:"context"`
}

type TestCase struct {
	TestCasePath  string
	Script        string
	TestCaseSetup []TestCaseSetup
}

func LoadTestCaseSetups(filepath string) ([]TestCaseSetup, error) {
	file, err := os.Open(filepath)
	if err != nil {
		os.Stderr.WriteString("error while opening test case setup file")
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(file)
	var testCaseSetups = make([]TestCaseSetup, 0)
	defer file.Close()

	err = json.Unmarshal(byteValue, &testCaseSetups)
	if err != nil {
		os.Stderr.WriteString("error while unmarshaling test case setup file" + filepath)
	}
	return testCaseSetups, err

}

func LoadTestCaseSetup(filepath string) (*TestCaseSetup, error) {
	file, err := os.Open(filepath)
	if err != nil {
		os.Stderr.WriteString("error while opening test case setup file")
		return nil, err
	}
	byteValue, _ := ioutil.ReadAll(file)
	defer file.Close()
	var testCaseSetup TestCaseSetup

	err = json.Unmarshal(byteValue, &testCaseSetup)
	if err != nil {
		os.Stderr.WriteString("error while unmarshaling test case setup file" + filepath)
	}
	return &testCaseSetup, err

}

type ESReq struct {
	Script       `json:"script"`
	ContextSetup `json:"context_setup"`
	Context      string `json:"context"`
}

type Script struct {
	Source string                 `json:"source"`
	Params map[string]interface{} `json:"params"`
}

type ContextSetup struct {
	Index    string                 `json:"index"`
	Document map[string]interface{} `json:"document"`
}

type Response struct {
	Result interface{} `json:"result"`
}

func ExecuteQuery(script string, setup TestCaseSetup, elasticsearhEndpoint string) (*Response, error) {
	context := ""
	if setup.Context == "" {
		context = "score"
	} else {
		context = setup.Context
	}
	req := ESReq{
		Script: Script{
			Source: script,
			Params: setup.Params,
		},
		ContextSetup: ContextSetup{
			Index:    setup.Index,
			Document: setup.Document,
		},
		Context: context,
	}

	jsonValue, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(elasticsearhEndpoint, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		json.Indent(buf, body, "", "  ")
		return nil, fmt.Errorf(buf.String())
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

///_scripts/painless/_execute
func RunTest(script string, setup TestCaseSetup, elasticsearhEndpoint string) error {

	response, err := ExecuteQuery(script, setup, elasticsearhEndpoint)
	if err != nil {
		return err
	}
	switch v := response.Result.(type) {
	case float64:
		a, ok := setup.ExpectedResult.(float64)
		if !ok {
			return fmt.Errorf("Error type %T %t", setup.ExpectedResult, setup.ExpectedResult)
		}
		if math.Abs(v-a) > delta {
			return fmt.Errorf("Expeted %f Got %f", setup.ExpectedResult, v)
		}

	case string:
		a, ok := setup.ExpectedResult.(string)
		if !ok {
			return fmt.Errorf("Error type %T %t", setup.ExpectedResult, setup.ExpectedResult)
		}

		if v != a {
			return fmt.Errorf("Error type %T %t", setup.ExpectedResult, setup.ExpectedResult)
		}

	case int:
		a, ok := setup.ExpectedResult.(int)
		if !ok {
			return fmt.Errorf("Error type %T %t", setup.ExpectedResult, setup.ExpectedResult)
		}

		if v != a {
			return fmt.Errorf("Expeted %d Got %d", a, v)
		}

	case bool:
		a, ok := setup.ExpectedResult.(bool)
		if !ok {
			return fmt.Errorf("Error type %T %t", setup.ExpectedResult, setup.ExpectedResult)
		}

		if v != a {
			return fmt.Errorf("Expeted %t Got %t", a, v)
		}

	default:
		return fmt.Errorf("I don't know about type %T!\n", v)
	}

	return nil

}

func RunTestCase(testCase TestCase, elasticsearhEndpoint string, parallel bool) int32 {

	if parallel {
		var wg sync.WaitGroup
		var failed int32
		for _, setup := range testCase.TestCaseSetup {
			wg.Add(1)
			s := setup
			go func() {
				defer wg.Done()
				testResult := RunTest(testCase.Script, s, elasticsearhEndpoint)
				WriteTestResult(testCase, s, testResult)
				if testResult != nil {
					atomic.AddInt32(&failed, 1)
				}
			}()
		}
		wg.Wait()
		failedFinal := atomic.LoadInt32(&failed)
		return failedFinal
	} else {
		var failed int32 = 0
		for _, setup := range testCase.TestCaseSetup {

			testResult := RunTest(testCase.Script, setup, elasticsearhEndpoint)
			WriteTestResult(testCase, setup, testResult)
			if testResult != nil {
				failed += 1
			}
		}
		return failed
	}

}

func WriteTestResult(testCase TestCase, setup TestCaseSetup, err error) {
	testMessage := setup.Description
	if err != nil {
		text := fmt.Sprintf("%s \n TestCase Failed at %s", testMessage, testCase.TestCasePath)
		errorMessage := fmt.Sprintf(" %s", err)
		fmt.Printf(ErrorColor, text)
		fmt.Println("")
		fmt.Printf(ErrorColor, errorMessage)
		fmt.Println("")
	} else {
		text := fmt.Sprintf("Test Passed: %s", testMessage)
		fmt.Printf(InfoColor, text)
		fmt.Println("")
	}

}
