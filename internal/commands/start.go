package commands

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:          "start",
	SilenceUsage: true,
	Short:        "Start API srv",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
		//if err != nil {
		//	return fmt.Errorf("failed to parse config: %w", err)
		//}
		//
		//instance, err := notifier.NewNotifier(cfg)
		//if err != nil {
		//	return fmt.Errorf("failed to create application: %w", err)
		//}
		//
		//errChan := make(chan error, 1)
		//okChan := make(chan struct{}, 1)
		//go func() {
		//	if err := instance.Start(); err != nil {
		//		errChan <- fmt.Errorf("failed to start application: %w", err)
		//	} else {
		//		okChan <- struct{}{}
		//	}
		//}()
		//
		//sigChan := make(chan os.Signal, 1)
		//signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		//
		//defer logger.Info("notifier stopped")
		//
		//select {
		//case <-sigChan:
		//	logger.Info("shutting down gracefully by signal")
		//
		//	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownDeadline)
		//	defer cancel()
		//
		//	instance.Shutdown(ctx)
		//	return nil
		//case <-okChan:
		//	return nil
		//case err := <-errChan:
		//	return err
		//}
	},
}
