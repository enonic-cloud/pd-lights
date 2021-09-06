package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

type Light int
type State int

const (
	Green Light = iota
	Yellow
	Red
)

const (
	On State = iota
	Off
)

func (l Light) Set(ctx context.Context, s State) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s/ctrl.cgi?F%d=%d", viper.GetString("ip"), l, s), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %s", err)
	}
	return err
}

func SetLights(ctx context.Context, red State, yellow State, green State) error {
	if err := Red.Set(ctx, red); err != nil {
		return fmt.Errorf("failed to set red light: %s", err)
	}
	if err := Yellow.Set(ctx, yellow); err != nil {
		return fmt.Errorf("failed to set yellow light: %s", err)
	}
	if err := Green.Set(ctx, green); err != nil {
		return fmt.Errorf("failed to set green light: %s", err)
	}
	return nil
}
