package lexmodels

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	botAliasDeletedTimeout = 5 * time.Minute
	intentDeletedTimeout   = 5 * time.Minute
)

func waitBotVersionCreated(ctx context.Context, conn *lexmodelbuildingservice.LexModelBuildingService, name, version string, timeout time.Duration) (*lexmodelbuildingservice.GetBotOutput, error) { //nolint:unparam
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{lexmodelbuildingservice.StatusBuilding},
		Target: []string{
			lexmodelbuildingservice.StatusNotBuilt,
			lexmodelbuildingservice.StatusReady,
			lexmodelbuildingservice.StatusReadyBasicTesting,
		},
		Refresh: statusBotVersion(ctx, conn, name, version),
		Timeout: timeout,
	}

	outputRaw, err := stateChangeConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*lexmodelbuildingservice.GetBotOutput); ok {
		if status := aws.StringValue(output.Status); status == lexmodelbuildingservice.StatusFailed {
			tfresource.SetLastError(err, errors.New(aws.StringValue(output.FailureReason)))
		}

		return output, err
	}

	return nil, err
}

func waitBotDeleted(ctx context.Context, conn *lexmodelbuildingservice.LexModelBuildingService, name string, timeout time.Duration) (*lexmodelbuildingservice.GetBotOutput, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{
			lexmodelbuildingservice.StatusNotBuilt,
			lexmodelbuildingservice.StatusReady,
			lexmodelbuildingservice.StatusReadyBasicTesting,
		},
		Target:  []string{},
		Refresh: statusBotVersion(ctx, conn, name, BotVersionLatest),
		Timeout: timeout,
	}

	outputRaw, err := stateChangeConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*lexmodelbuildingservice.GetBotOutput); ok {
		if status := aws.StringValue(output.Status); status == lexmodelbuildingservice.StatusFailed {
			tfresource.SetLastError(err, errors.New(aws.StringValue(output.FailureReason)))
		}

		return output, err
	}

	return nil, err
}

func waitBotAliasDeleted(ctx context.Context, conn *lexmodelbuildingservice.LexModelBuildingService, botAliasName, botName string) (*lexmodelbuildingservice.GetBotAliasOutput, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{serviceStatusCreated},
		Target:  []string{}, // An empty slice indicates that the resource is gone
		Refresh: statusBotAlias(ctx, conn, botAliasName, botName),
		Timeout: botAliasDeletedTimeout,
	}
	outputRaw, err := stateChangeConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*lexmodelbuildingservice.GetBotAliasOutput); ok {
		return v, err
	}

	return nil, err
}

func waitIntentDeleted(ctx context.Context, conn *lexmodelbuildingservice.LexModelBuildingService, intentId string) (*lexmodelbuildingservice.GetIntentVersionsOutput, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{serviceStatusCreated},
		Target:  []string{}, // An empty slice indicates that the resource is gone
		Refresh: statusIntent(ctx, conn, intentId),
		Timeout: intentDeletedTimeout,
	}
	outputRaw, err := stateChangeConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*lexmodelbuildingservice.GetIntentVersionsOutput); ok {
		return v, err
	}

	return nil, err
}

func waitSlotTypeDeleted(ctx context.Context, conn *lexmodelbuildingservice.LexModelBuildingService, name string) (*lexmodelbuildingservice.GetSlotTypeOutput, error) {
	stateChangeConf := &resource.StateChangeConf{
		Pending: []string{serviceStatusCreated},
		Target:  []string{},
		Refresh: statusSlotType(ctx, conn, name, SlotTypeVersionLatest),
		Timeout: slotTypeDeleteTimeout,
	}
	outputRaw, err := stateChangeConf.WaitForStateContext(ctx)

	if v, ok := outputRaw.(*lexmodelbuildingservice.GetSlotTypeOutput); ok {
		return v, err
	}

	return nil, err
}
