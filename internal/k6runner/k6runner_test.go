package k6runner

import (
	"context"
	"fmt"
	"github.com/loadimpact/k6/api"
	"github.com/loadimpact/k6/cmd"
	"github.com/loadimpact/k6/core"
	"github.com/loadimpact/k6/core/local"
	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)


const testJS = `
import { check } from "k6";
import http from 'k6/http'
import { randomIntBetween, randomItem,uuidv4 } from "https://jslib.k6.io/k6-utils/1.0.0/index.js";

let body = {
    "docId": 892530,
    "docTitle": "《北京银保监局行政处罚信息公开表》-（京银保监罚决字〔2020〕13号）",
    "docSubtitle": "《北京银保监局行政处罚信息公开表》-（京银保监罚决字〔2020〕13号）",
    "publishDate": "2020-03-02 18:16:50.0"
};

export let options = {
    tags: {
        case_name: 'bm'
    }
};


export default function () {
    const options = {
        headers: {
            'Content-Type': 'application/json'
        }
    };


    let docId = randomIntBetween(1,100);
    let res = http.post('http://127.0.0.1:8181/detail/' + docId, JSON.stringify(body), options)
};
`

func TestRunner(t *testing.T) {

	filesystems := loader.CreateFilesystems()
	//pwd, err := os.Getwd()
	//pwd, err := os.Getwd()
	//if err != nil {
	//	t.Error(err)
	//}
	//src, err := loader.ReadSource("/home/master/works/k6/bm-test.js", pwd, filesystems, os.Stdin)


	src := loader.SourceData{URL: &url.URL{Path: "/-", Scheme: "file"}, Data: []byte(testJS)}
	//if err != nil {
	//	t.Error(err)
	//}

	//runtimeOptions, err := getRuntimeOptions(cmd.Flags())
	//if err != nil {
	//	return err
	//}

	rtOpts := lib.RuntimeOptions{
		IncludeSystemEnvVars: null.BoolFrom(false),
		CompatibilityMode:    null.StringFrom("extended"),
		Env:                  make(map[string]string),
	}

	r, err := js.New(&src, filesystems, rtOpts)
	if err != nil {
		t.Error(err)
	}

	ex := local.New(r)

	opts := lib.Options{
		VUs:      null.IntFrom(1),
		VUsMax:   null.IntFrom(1),
		Duration: types.NullDurationFrom(time.Second * 3),
		//Iterations:            getNullInt64(flags, "iterations"),
		//Paused:                getNullBool(flags, "paused"),
		//MaxRedirects:          getNullInt64(flags, "max-redirects"),
		//Batch:                 getNullInt64(flags, "batch"),
		//BatchPerHost:          getNullInt64(flags, "batch-per-host"),
		//RPS:                   getNullInt64(flags, "rps"),
		//UserAgent:             getNullString(flags, "user-agent"),
		//HTTPDebug:             getNullString(flags, "http-debug"),
		//InsecureSkipTLSVerify: getNullBool(flags, "insecure-skip-tls-verify"),
		//NoConnectionReuse:     getNullBool(flags, "no-connection-reuse"),
		//NoVUConnectionReuse:   getNullBool(flags, "no-vu-connection-reuse"),
		//MinIterationDuration:  getNullDuration(flags, "min-iteration-duration"),
		//Throw:                 getNullBool(flags, "throw"),
		//DiscardResponseBodies: getNullBool(flags, "discard-response-bodies"),
		// Default values for options without CLI flags:
		// TODO: find a saner and more dev-friendly and error-proof way to handle options
		SetupTimeout:    types.NullDuration{Duration: types.Duration(10 * time.Second), Valid: false},
		TeardownTimeout: types.NullDuration{Duration: types.Duration(10 * time.Second), Valid: false},

		MetricSamplesBufferSize: null.NewInt(1000, false),
	}

	conf := cmd.Config{
		Options: opts,
		Out:     []string{},
		//Linger:        getNullBool(flags, "linger"),
		//NoUsageReport: getNullBool(flags, "no-usage-report"),
		//NoThresholds:  getNullBool(flags, "no-thresholds"),
		//NoSummary:     getNullBool(flags, "no-summary"),
		//SummaryExport: getNullString(flags, "summary-export"),
	}

	engine, err := core.NewEngine(ex, conf.Options)

	if err != nil {
		t.Error(err)
	}

	go func() {
		if err := api.ListenAndServe(":9191", engine); err != nil {
			logrus.WithError(err).Warn("Error from API server")
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	errC := make(chan error)
	go func() { errC <- engine.Run(ctx) }()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigC)

	//updateFreq := 50 * time.Millisecond
	//ticker := time.NewTicker(updateFreq)
	mainLoop:
	for {

		select {
		//case <-ticker.C:
		//	fmt.Printf("running %s \n", time.Now())
		case err:=<-errC:
			cancel()
			if err == nil {
				fmt.Printf("Done %s \n", time.Now())

				logrus.Debug("Engine terminated cleanly")
				break mainLoop
			}
		case sig := <-sigC:
			logrus.WithField("sig", sig).Debug("Exiting in response to signal")
			cancel()
		}

	}
}
