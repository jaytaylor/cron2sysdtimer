package crontab

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func TestParse(t *testing.T) {
	filename := filepath.Join("..", "testdata", "crontab")
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to open testdata. filename: %s, error: %s", filename, err)
	}

	expected := []*Schedule{
		&Schedule{
			Spec:    "0,5,10,15,20,25,30,35,40,45,50,55 * * * *",
			Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task01.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task01 RAILS_ENV=production'",
		},
		&Schedule{
			Spec:    "15 * * * *",
			Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task02.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task02 RAILS_ENV=production'",
		},
		&Schedule{
			Spec:    "30 * * * *",
			Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task04.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task04 RAILS_ENV=production'",
		},
	}

	actual, err := Parse(string(body))
	if err != nil {
		t.Fatalf("Unexpected error parsing body string %q: %s", body, err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Schedules do not match expected values.\n  expected: %+v\n  actual:      %+v", expected, actual)
	}
}

func TestConvertToSystemdCalendar(t *testing.T) {
	testcases := []struct {
		schedule *Schedule
		expected string
	}{
		{
			schedule: &Schedule{
				Spec:    "*/5 * * * *",
				Command: "",
			},
			expected: "*:0,5,10,15,20,25,30,35,40,45,50,55",
		},
		{
			schedule: &Schedule{
				Spec:    "0,5,10,15,20,25,30,35,40,45,50,55 10-12 * * *",
				Command: "",
			},
			expected: "10,11,12:0,5,10,15,20,25,30,35,40,45,50,55", // TODO: 10-12:0,5,...
		},
		{
			schedule: &Schedule{
				Spec:    "0-5 * 1 * *",
				Command: "",
			},
			expected: "*-1 *:0,1,2,3,4,5", // TODO: *:0-5
		},
		{
			schedule: &Schedule{
				Spec:    "23 2,1 * 12 1,6",
				Command: "",
			},
			expected: "Mon,Sat 12-* 1,2:23",
		},
		{
			schedule: &Schedule{
				Spec:    "0,20,40 8-17 * * 1-5",
				Command: "",
			},
			expected: "Mon,Tue,Wed,Thu,Fri 8,9,10,11,12,13,14,15,16,17:0,20,40",
		},
		{
			schedule: &Schedule{
				Spec:    "0 17 * * *",
				Command: "",
			},
			expected: "17:0",
		},
		{
			schedule: &Schedule{
				Spec:    "* * * * 0",
				Command: "",
			},
			expected: "Sun *:*",
		},
		{
			schedule: &Schedule{
				Spec:    "5 * * * *",
				Command: "",
			},
			expected: "*:5",
		},
	}

	for _, tc := range testcases {
		actual, err := tc.schedule.ConvertToSystemdCalendar()
		if err != nil {
			t.Errorf("Error should not be raised. error: %s", err)
		}

		if actual != tc.expected {
			t.Errorf("Calendar does not match. expected: %q, actual: %q", tc.expected, actual)
		}
	}
}

func TestNameByRegexp(t *testing.T) {
	testcases := []struct {
		schedule   *Schedule
		nameRegexp *regexp.Regexp
		expected   string
	}{
		{
			schedule: &Schedule{
				Spec:    "0,5,10,15,20,25,30,35,40,45,50,55 * * * *",
				Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task01.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task01 RAILS_ENV=production'",
			},
			nameRegexp: regexp.MustCompile(`--name ([a-zA-Z0-9.]+)`),
			expected:   "scheduler.task01",
		},
		{
			schedule: &Schedule{
				Spec:    "15 * * * *",
				Command: "/bin/echo hello",
			},
			nameRegexp: regexp.MustCompile(`--name ([a-zA-Z0-9.]+)`),
			expected:   "",
		},
		{
			schedule: &Schedule{
				Spec:    "30 * * * *",
				Command: "/bin/docker run --name hello ubuntu:16.04 echo hello",
			},
			nameRegexp: regexp.MustCompile(`--name ([a-zA-Z0-9.]+)`),
			expected:   "hello",
		},
		{
			schedule: &Schedule{
				Spec:    "30 * * * *",
				Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task01_--.._.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task01 RAILS_ENV=production'",
			},
			nameRegexp: regexp.MustCompile(`--name ([a-zA-Z0-9.]+)`),
			expected:   "scheduler.task01",
		},
		{
			schedule: &Schedule{
				Spec:    "30 * * * *",
				Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task01.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task01 RAILS_ENV=production'",
			},
			nameRegexp: regexp.MustCompile(``),
			expected:   "",
		},
		{
			schedule: &Schedule{
				Spec:    "30 * * * *",
				Command: "/bin/bash -l -c 'docker run --rm=true --name scheduler.task01.`date +\\%Y\\%m\\%d\\%H\\%M` --memory=5g 123456789012.dkr.ecr.ap-northeast-1.amazonaws.com/app:latest bundle exec rake task01 RAILS_ENV=production'",
			},
			nameRegexp: nil,
			expected:   "",
		},
	}

	for _, tc := range testcases {
		if actual := tc.schedule.NameByRegexp(tc.nameRegexp); actual != tc.expected {
			t.Errorf("Name does not match. expected: %q, actual: %q", tc.expected, actual)
		}
	}
}

func TestSHA256Sum(t *testing.T) {
	schedule := &Schedule{
		Spec:    "15 * * * *",
		Command: "echo 'hello'",
	}
	expected := "4ab7fd35a3996a8b58483a640a52976d5c974372c12e5f7a973be86d96a0096e"

	if actual := schedule.SHA256Sum(); actual != expected {
		t.Errorf("Checksum does not match. expected: %q, actual: %q", expected, actual)
	}
}
