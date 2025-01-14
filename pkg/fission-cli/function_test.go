package fission_cli

import (
	"flag"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"

	fv1 "github.com/fission/fission/pkg/apis/fission.io/v1"
)

func TestGetInvokeStrategy(t *testing.T) {
	cases := []struct {
		testArgs               map[string]string
		existingInvokeStrategy *fv1.InvokeStrategy
		expectedResult         *fv1.InvokeStrategy
		expectError            bool
	}{
		{
			// case: use default executor poolmgr
			testArgs:               map[string]string{},
			existingInvokeStrategy: nil,
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType: fv1.ExecutorTypePoolmgr,
				},
			},
			expectError: false,
		},
		{
			// case: executor type set to poolmgr
			testArgs:               map[string]string{"executortype": fv1.ExecutorTypePoolmgr},
			existingInvokeStrategy: nil,
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType: fv1.ExecutorTypePoolmgr,
				},
			},
			expectError: false,
		},
		{
			// case: executor type set to newdeploy
			testArgs:               map[string]string{"executortype": fv1.ExecutorTypeNewdeploy},
			existingInvokeStrategy: nil,
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              DEFAULT_MIN_SCALE,
					MaxScale:              DEFAULT_MIN_SCALE,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectError: false,
		},
		{
			// case: executor type change from poolmgr to newdeploy
			testArgs: map[string]string{"executortype": fv1.ExecutorTypeNewdeploy},
			existingInvokeStrategy: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType: fv1.ExecutorTypePoolmgr,
				},
			},
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              DEFAULT_MIN_SCALE,
					MaxScale:              DEFAULT_MIN_SCALE,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectError: false,
		},
		{
			// case: executor type change from newdeploy to poolmgr
			testArgs: map[string]string{"executortype": fv1.ExecutorTypePoolmgr},
			existingInvokeStrategy: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              DEFAULT_MIN_SCALE,
					MaxScale:              DEFAULT_MIN_SCALE,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType: fv1.ExecutorTypePoolmgr,
				},
			},
			expectError: false,
		},
		{
			// case: minscale < maxscale
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
				"minscale":     "2",
				"maxscale":     "3",
			},
			existingInvokeStrategy: nil,
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              2,
					MaxScale:              3,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectError: false,
		},
		{
			// case: minscale > maxscale
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
				"minscale":     "5",
				"maxscale":     "3",
			},
			existingInvokeStrategy: nil,
			expectedResult:         nil,
			expectError:            true,
		},
		{
			// case: maxscale not specified
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
				"minscale":     "5",
			},
			existingInvokeStrategy: nil,
			expectedResult:         nil,
			expectError:            true,
		},
		{
			// case: minscale not specified
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
				"maxscale":     "3",
			},
			existingInvokeStrategy: nil,
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              DEFAULT_MIN_SCALE,
					MaxScale:              3,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectError: false,
		},
		{
			// case: maxscale set to 0
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
				"maxscale":     "0",
			},
			existingInvokeStrategy: nil,
			expectedResult:         nil,
			expectError:            true,
		},
		{
			// case: maxscale set to 9 when existing is 5
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
				"maxscale":     "9",
			},
			existingInvokeStrategy: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              2,
					MaxScale:              5,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              2,
					MaxScale:              9,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectError: false,
		},
		{
			// case: change nothing for existing strategy
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
			},
			existingInvokeStrategy: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              2,
					MaxScale:              5,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              2,
					MaxScale:              5,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectError: false,
		},
		{
			// case: set target cpu percentage
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
				"targetcpu":    "50",
			},
			existingInvokeStrategy: nil,
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              DEFAULT_MIN_SCALE,
					MaxScale:              DEFAULT_MIN_SCALE,
					TargetCPUPercent:      50,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectError: false,
		},
		{
			// case: change target cpu percentage
			testArgs: map[string]string{
				"executortype": fv1.ExecutorTypeNewdeploy,
				"targetcpu":    "20",
			},
			existingInvokeStrategy: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              2,
					MaxScale:              5,
					TargetCPUPercent:      88,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              2,
					MaxScale:              5,
					TargetCPUPercent:      20,
					SpecializationTimeout: DEFAULT_SPECIALIZATION_TIMEOUT,
				},
			},
			expectError: false,
		},
		{
			// case: change specializationtimeout
			testArgs: map[string]string{
				"executortype":          fv1.ExecutorTypeNewdeploy,
				"specializationtimeout": "200",
			},
			existingInvokeStrategy: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:     fv1.ExecutorTypeNewdeploy,
					MinScale:         2,
					MaxScale:         5,
					TargetCPUPercent: DEFAULT_TARGET_CPU_PERCENTAGE,
				},
			},
			expectedResult: &fv1.InvokeStrategy{
				StrategyType: fv1.StrategyTypeExecution,
				ExecutionStrategy: fv1.ExecutionStrategy{
					ExecutorType:          fv1.ExecutorTypeNewdeploy,
					MinScale:              2,
					MaxScale:              5,
					SpecializationTimeout: 200,
					TargetCPUPercent:      DEFAULT_TARGET_CPU_PERCENTAGE,
				},
			},
			expectError: false,
		},
		{
			// case: specializationtimeout should not work for poolmgr
			testArgs: map[string]string{
				"executortype":          fv1.ExecutorTypePoolmgr,
				"specializationtimeout": "10",
			},
			existingInvokeStrategy: nil,
			expectedResult:         nil,
			expectError:            true,
		},
		{
			// case: specializationtimeout should not be less than 120
			testArgs: map[string]string{
				"executortype":          fv1.ExecutorTypeNewdeploy,
				"specializationtimeout": "90",
			},
			existingInvokeStrategy: nil,
			expectedResult:         nil,
			expectError:            true,
		},
	}

	for i, c := range cases {
		fmt.Printf("=== Test Case %v ===\n", i)

		app := NewCliApp()
		set := flag.NewFlagSet("test-cmd", 0)
		ctx := cli.NewContext(app, set, nil)

		for k, v := range c.testArgs {
			set.String(k, v, "")
			ctx.Set(k, v)
		}

		strategy, err := getInvokeStrategy(ctx, c.existingInvokeStrategy)
		if c.expectError {
			assert.NotNil(t, err)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			assert.Nil(t, err)
			assert.NoError(t, strategy.Validate(), fmt.Sprintf("Failed at test case %v", i))
			assert.Equal(t, *c.expectedResult, *strategy)
		}
	}
}
