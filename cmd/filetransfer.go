// Copyright 2023 Northern.tech AS
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	    http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.
package cmd

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mendersoftware/mender-cli/client/deviceconnect"
	"github.com/mendersoftware/mender-cli/log"
)

const (
	deviceDelimiter = ":"
)

var fileTransferCmd = &cobra.Command{
	Use:   "cp device_id:file_path file_path",
	Short: "Transfer files from/to a device",
	Long:  "A CLI interface for copying files from/to devices in your setup",
	Args:  cobra.MinimumNArgs(2),
	Run: func(c *cobra.Command, args []string) {
		cmd, err := NewFileTransfer(c, args)
		CheckErr(err)
		CheckErr(cmd.Run())
	},
}

type FileTransferCmd struct {
	server      string
	skipVerify  bool
	source      string
	destination string
	token       string
}

func NewFileTransfer(cmd *cobra.Command, args []string) (*FileTransferCmd, error) {
	server := viper.GetString(argRootServer)
	if server == "" {
		return nil, errors.New("Empty server value. This should never happen")
	}

	skipVerify, err := cmd.Flags().GetBool(argRootSkipVerify)
	if err != nil {
		return nil, err
	}

	token, err := getAuthToken(cmd)
	if err != nil {
		return nil, err
	}

	return &FileTransferCmd{
		server:      server,
		skipVerify:  skipVerify,
		token:       token,
		source:      args[0],
		destination: args[1],
	}, nil
}

func (c *FileTransferCmd) Run() error {
	if strings.Contains(c.destination, ":") {
		return c.upload()
	}
	return c.download()
}

func (c *FileTransferCmd) checkDevice(deviceID string) error {
	// check if the device is connected
	client := deviceconnect.NewClient(c.server, c.token, c.skipVerify)
	device, err := client.GetDevice(deviceID)
	if err != nil {
		return errors.Wrap(err, "unable to get the device")
	} else if device.Status != deviceconnect.CONNECTED {
		return errors.New("the device is not connected")
	}

	return nil
}

func (c *FileTransferCmd) upload() error {
	d, err := deviceSpecification(c.destination)
	if err == nil {
		err = c.checkDevice(d.DeviceID)
	}
	if err != nil {
		return err
	}
	client := deviceconnect.NewFileTransferClient(c.server, c.token, c.skipVerify)
	if err = client.Upload(c.source, d); err != nil {
		return err
	}
	log.Infof("Successfully uploaded the file %q to device %q at location %q\n",
		c.source, d.DeviceID, d.DevicePath)
	return nil
}

func deviceSpecification(s string) (*deviceconnect.DeviceSpec, error) {
	d := strings.Split(s, deviceDelimiter)
	if len(d) > 2 {
		return nil, errors.New("The device specification contains multiple ':' delimeters")
	}
	if len(d) != 2 {
		return nil, errors.New("The device specification is missing the ':' separator")
	}
	return &deviceconnect.DeviceSpec{DeviceID: d[0], DevicePath: d[1]}, nil
}

func (c *FileTransferCmd) download() error {
	d, err := deviceSpecification(c.source)
	if err == nil {
		err = c.checkDevice(d.DeviceID)
	}
	if err != nil {
		return err
	}
	client := deviceconnect.NewFileTransferClient(c.server, c.token, c.skipVerify)
	if err = client.Download(d, c.destination); err != nil {
		return err
	}
	log.Infof("Successfully downloaded the file: %q from device %q to %q\n",
		d.DevicePath, d.DeviceID, c.source)
	return nil
}
