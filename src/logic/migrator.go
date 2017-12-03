// Copyright 2017 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author:javasgl	songganglin@gmail.com

package logic

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/polaris1119/config"
	"github.com/polaris1119/logger"
)

type MigratorLogic struct{}

var (
	wg           sync.WaitGroup
	liquibaseLib string
)

var DefaultMigrator = MigratorLogic{}

func (MigratorLogic) Migrator(changeVersion string) {

	liquibaseLib = config.ConfigFile.MustValue("migrator", "liquibase_lib_dir")

	if !filepath.IsAbs(liquibaseLib) {
		liquibaseLib = config.ROOT + "/" + liquibaseLib
	}

	changeLogDir := config.ConfigFile.MustValue("migrator", "change_log_dir")
	if !filepath.IsAbs(changeLogDir) {
		changeLogDir = config.ROOT + "/" + changeLogDir
	}
	versionDir := changeLogDir + "/" + changeVersion
	changeLogVersion, err := os.Stat(versionDir)

	logger.Infoln("migrator:changelog dir is:", versionDir)

	if err == nil && changeLogVersion.IsDir() {

		logger.Infoln("migrator:exec changelog version:", changeVersion, ", files in ", versionDir)

		changeLogs, err := ioutil.ReadDir(versionDir)
		if err != nil {
			logger.Errorln("migrator:read changelog files error:", err)
			os.Exit(1)
		}

		for _, changeLog := range changeLogs {
			if strings.HasSuffix(changeLog.Name(), ".xml") {

				database := strings.TrimSuffix(changeLog.Name(), ".xml")

				wg.Add(1)
				go execDatabaseChange(database, versionDir+"/"+changeLog.Name())

			}
		}
		wg.Wait()

	} else {
		logger.Errorln("migrator:read changelog version dir error:", err)
		os.Exit(1)
	}
}
func execDatabaseChange(database, changeLog string) {
	defer wg.Done()

	args := []string{
		"--driver=com.mysql.jdbc.Driver",
		"--username=" + config.ConfigFile.MustValue("mysql", "user"),
		"--password=" + config.ConfigFile.MustValue("mysql", "password"),
		"--url=" + fmt.Sprintf("jdbc:mysql://%s:%s/%s?characterEncoding=utf8", config.ConfigFile.MustValue("mysql", "host"), config.ConfigFile.MustValue("mysql", "port"), database),
		"--changeLogFile=" + changeLog,
		"--classpath=" + liquibaseLib + "/lib/mysql-connector-java-5.1.25-bin.jar",
		"--logLevel=info",
		"update",
	}
	cmd := exec.Command(liquibaseLib+"/liquibase", args...)

	fmt.Println(strings.Join(cmd.Args, " "))

	output, err := cmd.CombinedOutput()

	fmt.Println(string(output))

	if err != nil {
		fmt.Println("Error:", err)
	}
}
