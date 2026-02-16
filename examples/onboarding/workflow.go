package onboarding

import (
	"durable-engine/engine"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

func OnboardingWorkflow(ctx *engine.Context) error {

	// Step 1: Create Record (Sequential)
	userID, err := engine.Step(ctx, "create-record", func() (string, error) {

		fmt.Println("Creating employee record...")

		time.Sleep(2 * time.Second)

		return "EMP-123", nil
	})

	if err != nil {
		return err
	}

	// Step 2 & 3: Parallel Steps
	var g errgroup.Group

	g.Go(func() error {

		_, err := engine.Step(ctx, "provision-laptop", func() (string, error) {

			fmt.Println("Provisioning laptop...")

			time.Sleep(3 * time.Second)

			return "LAPTOP-ASSIGNED", nil
		})

		return err
	})

	g.Go(func() error {

		_, err := engine.Step(ctx, "provision-access", func() (string, error) {

			fmt.Println("Provisioning system access...")

			time.Sleep(3 * time.Second)

			return "ACCESS-GRANTED", nil
		})

		return err
	})

	err = g.Wait()
	if err != nil {
		return err
	}

	// Step 4: Send Welcome Email
	_, err = engine.Step(ctx, "send-email", func() (string, error) {

		fmt.Println("Sending welcome email to:", userID)

		time.Sleep(2 * time.Second)

		return "EMAIL-SENT", nil
	})

	return err
}
