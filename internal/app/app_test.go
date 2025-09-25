package app

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/require"
)

func TestAppCreateAndRun(t *testing.T) {
	t.Run("Returns error if it can not wire dependencies", func(t *testing.T) {
		original := wireDependencies
		defer func() {
			wireDependencies = original
		}()
		
		wireDependencies = func() (Dependencies, error) {
			return Dependencies{}, errors.New("test error")
		}
		
		err := CreateAndRun("path")
		require.Error(t, err)		
	})
	
	t.Run("Creates the application, if wire dependencies don't return error", func(t *testing.T) {
		original := wireDependencies
		defer func() {
			wireDependencies = original
		}()
		
		wireDependencies = func() (Dependencies, error) {
			return Dependencies{}, nil
		}
		
		_, err := newApp(wireDependencies)
		require.NoError(t, err)	
	})
}