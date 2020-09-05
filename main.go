// Licensed to Michael Tougeron <github@e.tougeron.com> under
// one or more contributor license agreements. See the NOTICE
// file distributed with this work for additional information
// regarding copyright ownership.
// Michael Tougeron <github@e.tougeron.com> licenses this file
// to you under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Comcast/kuberhealthy/v2/pkg/checks/external/checkclient"
	"github.com/Comcast/kuberhealthy/v2/pkg/kubeClient"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var (
	buildVersion string = ""
	buildTime    string = ""

	kubeConfigFile string = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	namespace      string = os.Getenv("NAMESPACE")
	// In cluster API calls should be quick
	maxDuration           = 250
	debugEnv       string = os.Getenv("DEBUG")
	maxDurationEnv string = os.Getenv("MAX_DURATION_MILLISECONDS")
)

func init() {

	if maxDurationEnv != "" {
		maxDuration, _ = strconv.Atoi(os.Getenv("MAX_DURATION_MILLISECONDS"))
	}

	if namespace == "" {
		namespace = "kuberhealthy"
	}

	// Enable Debug just in case
	if len(debugEnv) != 0 {
		debug, err := strconv.ParseBool(debugEnv)
		if err != nil {
			log.Fatalln("Failed to parse DEBUG Environment variable:", err.Error())
		}
		checkclient.Debug = debug
		log.SetLevel(log.DebugLevel)
	}

	// APP Build information
	log.Debugln("Application Version:", buildVersion)
	log.Debugln("Application Build Time:", buildTime)
}

func main() {

	client, err := kubeClient.Create(kubeConfigFile)
	if err != nil {
		log.Fatalln("Unable to create kubernetes client")
		log.Debugln(err.Error())
		os.Exit(1)
	}

	// const apiCallLatencyThreshold = maxDuration * time.Second
	apiCallLatencyThreshold, _ := time.ParseDuration(strconv.Itoa(maxDuration) + "ms")

	log.Debugln("Beginning check of namespace: " + namespace)
	timeStart := time.Now()
	pods, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	apiResponseTime := time.Since(timeStart)
	if err != nil {
		log.Errorln("Error getting pods")
		log.Debugln(err.Error())
		err = checkclient.ReportFailure([]string{"Could not get pods: " + err.Error()})
		os.Exit(1)
	}
	if len(pods.Items) == 0 {
		log.Errorf("There were no pods found in namespace: %s", namespace)
		err = checkclient.ReportFailure([]string{"There were no pods found in namespace: " + namespace})
		os.Exit(1)
	}
	log.Debugf("Took %s time getting pods", apiResponseTime.String())
	if apiResponseTime > apiCallLatencyThreshold {
		err = checkclient.ReportFailure([]string{"Took too long getting pods: " + apiResponseTime.String()})
		if err != nil {
			log.Errorln("Error reporting failure to Kuberhealthy servers")
			log.Debugln(err.Error())
			os.Exit(1)
		}
		os.Exit(1)
	}

	// report success to Kuberhealthy servers if there were no failed pods in our list.
	err = checkclient.ReportSuccess()
	log.Debugln("Reporting Success, found pods.")
	if err != nil {
		log.Errorln("Error reporting success to Kuberhealthy servers")
		log.Debugln(err.Error())
		os.Exit(1)
	}
}
