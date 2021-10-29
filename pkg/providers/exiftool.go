package providers

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"goexiftool/pkg"
	"goexiftool/pkg/dtos"
	"os/exec"
	"strings"
)

var _ pkg.Provider = &Exiftool{}

type Exiftool struct {
	cmd  *exec.Cmd
	exit chan struct{}
}

func (s *Exiftool) ConsumeFile(file string) (<-chan *dtos.Response, error) {
	responseChan := make(chan *dtos.Response)
	finalResponse := ""
	cmdName := "exiftool -listx"
	cmd := exec.Command("sh", "-c", cmdName)
	s.cmd = cmd

	stdout, _ := cmd.StdoutPipe()
	reader := bufio.NewReader(stdout)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-s.exit:
				if err := cmd.Process.Kill(); err != nil {
					fmt.Printf("failed to kill: %v", err)
				}
				return
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				finalResponse += line
				if strings.Contains(line, "</table>") {
					var taginfo string
					if !strings.Contains(line, "</taginfo>") {
						taginfo = "</taginfo>"
					}
					response, err := s.digest(finalResponse + taginfo)
					if err != nil {
						break
					}
					responseChan <- response
				}
			}
		}
	}()

	return responseChan, nil
}

//wrap exiftool in a function to make it easier to call
func (s *Exiftool) digest(data string) (*dtos.Response, error) {
	var exiftoolResponse dtos.ExiftoolXML
	err := xml.Unmarshal([]byte(data), &exiftoolResponse)
	if err != nil {
		return nil, err
	}

	response := dtos.Response{}
	for e := range exiftoolResponse.Tables {
		description := make(map[string]interface{})
		for _, e := range exiftoolResponse.Tables[e].Tag.Desc {
			description[e.Lang] = e.Value
		}
		tag := &dtos.Data{
			"writable":    exiftoolResponse.Tables[e].Tag.Writable,
			"tag":         fmt.Sprintf("%s:%s", exiftoolResponse.Tables[e].Name, exiftoolResponse.Tables[e].Tag.Name),
			"group":       exiftoolResponse.Tables[e].Name,
			"description": description,
		}
		response.Tags = append(response.Tags, tag)
	}

	return &response, nil
}

func (s *Exiftool) Close() {
	go func() {
		s.exit <- struct{}{}
	}()
}

func New() *Exiftool {
	return &Exiftool{
		exit: make(chan struct{}),
	}
}
